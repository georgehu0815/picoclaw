package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/sipeed/picoclaw/pkg/auth"
)

// AzureConfig holds Azure OpenAI configuration with managed identity support
// Similar to azure-openai-models.ts configuration
type AzureConfig struct {
	Endpoint             string // Azure OpenAI endpoint URL
	Deployment           string // Azure OpenAI deployment name
	APIVersion           string // Azure OpenAI API version
	Scope                string // Azure OpenAI scope for authentication
	ManagedIdentityID    string // Client ID for user-assigned managed identity (optional)
	UseManagedIdentity   bool   // Enable managed identity authentication
	Verbose              bool   // Enable debug logging
}

type CodexProvider struct {
	client      *openai.Client
	accountID   string
	tokenSource func() (string, string, error)
	azureConfig *AzureConfig // Azure-specific configuration
}

const defaultCodexInstructions = "You are Codex, a coding assistant."

func NewCodexProvider(token, accountID string) *CodexProvider {
	opts := []option.RequestOption{
		option.WithBaseURL("https://chatgpt.com/backend-api/codex"),
		option.WithAPIKey(token),
	}
	if accountID != "" {
		opts = append(opts, option.WithHeader("Chatgpt-Account-Id", accountID))
	}
	client := openai.NewClient(opts...)
	return &CodexProvider{
		client:    &client,
		accountID: accountID,
	}
}

func NewCodexProviderWithTokenSource(token, accountID string, tokenSource func() (string, string, error)) *CodexProvider {
	p := NewCodexProvider(token, accountID)
	p.tokenSource = tokenSource
	return p
}

// NewCodexProviderWithAzure creates a provider configured for Azure OpenAI
// Similar to azure-openai-models.ts configuration
func NewCodexProviderWithAzure(azureConfig *AzureConfig, initialToken string) (*CodexProvider, error) {
	if azureConfig == nil {
		return nil, fmt.Errorf("Azure configuration is required")
	}

	// Build Azure OpenAI endpoint URL
	// Base URL without query parameters (added per-request)
	baseURL := fmt.Sprintf("%s/openai/deployments/%s",
		strings.TrimRight(azureConfig.Endpoint, "/"),
		azureConfig.Deployment,
	)

	opts := []option.RequestOption{
		option.WithBaseURL(baseURL),
		// api-version will be added per-request in chatAzure()
	}

	if initialToken != "" {
		opts = append(opts, option.WithAPIKey(initialToken))
	}

	client := openai.NewClient(opts...)

	// Create token source with Azure managed identity support
	tokenSource := createDynamicCodexTokenSource(azureConfig)

	return &CodexProvider{
		client:      &client,
		tokenSource: tokenSource,
		azureConfig: azureConfig,
	}, nil
}

// NewCodexProviderAuto creates a provider with automatic configuration detection
// Checks for Azure configuration first, then falls back to standard OpenAI
func NewCodexProviderAuto() (*CodexProvider, error) {
	// Try to load Azure configuration from environment
	azureConfig, err := LoadAzureConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load Azure config: %w", err)
	}

	// If Azure is configured, use Azure provider
	if azureConfig != nil {
		if azureConfig.Verbose {
			fmt.Println("[CodexProvider] Using Azure OpenAI configuration - codex_provider.go:109")
		}
		return NewCodexProviderWithAzure(azureConfig, "")
	}

	// Otherwise, use standard OpenAI with dynamic token source
	tokenSource := createCodexTokenSource()
	token, accountID, err := tokenSource()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve OpenAI token: %w", err)
	}

	provider := NewCodexProvider(token, accountID)
	provider.tokenSource = tokenSource
	return provider, nil
}

// NewCodexProviderWithDynamicAuth creates a provider with enhanced authentication
// Supports both Azure Managed Identity and OpenAI OAuth
func NewCodexProviderWithDynamicAuth(azureConfig *AzureConfig) (*CodexProvider, error) {
	if azureConfig != nil {
		return NewCodexProviderWithAzure(azureConfig, "")
	}
	return NewCodexProviderAuto()
}

func (p *CodexProvider) Chat(ctx context.Context, messages []Message, tools []ToolDefinition, model string, options map[string]interface{}) (*LLMResponse, error) {
	var opts []option.RequestOption
	if p.tokenSource != nil {
		tok, accID, err := p.tokenSource()
		if err != nil {
			return nil, fmt.Errorf("refreshing token: %w", err)
		}
		opts = append(opts, option.WithAPIKey(tok))
		if accID != "" {
			opts = append(opts, option.WithHeader("Chatgpt-Account-Id", accID))
		}
	}

	// Azure OpenAI uses Chat Completions API, not Responses API
	if p.azureConfig != nil {
		if p.azureConfig.Verbose {
			fmt.Println("[CodexProvider] Using Azure OpenAI Chat Completions API - codex_provider.go:151")
		}
		return p.chatAzure(ctx, messages, tools, model, options, opts)
	}

	// Standard OpenAI uses Responses API
	params := buildCodexParams(messages, tools, model, options)

	resp, err := p.client.Responses.New(ctx, params, opts...)
	if err != nil {
		return nil, fmt.Errorf("codex API call: %w", err)
	}

	return parseCodexResponse(resp), nil
}

// chatAzure handles Azure OpenAI Chat Completions API
func (p *CodexProvider) chatAzure(ctx context.Context, messages []Message, tools []ToolDefinition, model string, options map[string]interface{}, opts []option.RequestOption) (*LLMResponse, error) {
	// Build chat completion parameters for Azure
	chatMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			chatMessages = append(chatMessages, openai.SystemMessage(msg.Content))
		case "user":
			chatMessages = append(chatMessages, openai.UserMessage(msg.Content))
		case "assistant":
			chatMessages = append(chatMessages, openai.AssistantMessage(msg.Content))
		case "tool":
			chatMessages = append(chatMessages, openai.ToolMessage(msg.ToolCallID, msg.Content))
		}
	}

	params := openai.ChatCompletionNewParams{
		Messages: chatMessages,
		Model:    model,
	}

	// Azure OpenAI uses max_completion_tokens instead of max_tokens
	if maxTokens, ok := options["max_tokens"].(int); ok {
		params.MaxCompletionTokens = openai.Int(int64(maxTokens))
	}

	// Skip temperature for Azure OpenAI gpt-5.2-chat (only supports default value of 1)
	// if temp, ok := options["temperature"].(float64); ok {
	// 	params.Temperature = openai.Float(temp)
	// }

	// Add api-version query parameter (required by Azure OpenAI)
	opts = append(opts, option.WithQuery("api-version", p.azureConfig.APIVersion))

	// Call Azure OpenAI Chat Completions API
	resp, err := p.client.Chat.Completions.New(ctx, params, opts...)
	if err != nil {
		return nil, fmt.Errorf("Azure OpenAI API call: %w", err)
	}

	// Parse Azure response
	return parseChatCompletionResponse(resp), nil
}

// parseChatCompletionResponse converts Azure OpenAI chat completion response to LLMResponse
func parseChatCompletionResponse(resp *openai.ChatCompletion) *LLMResponse {
	if len(resp.Choices) == 0 {
		return &LLMResponse{
			Content:      "",
			FinishReason: "error",
		}
	}

	choice := resp.Choices[0]
	message := choice.Message

	var toolCalls []ToolCall
	if len(message.ToolCalls) > 0 {
		toolCalls = make([]ToolCall, 0, len(message.ToolCalls))
		for _, tc := range message.ToolCalls {
			var args map[string]interface{}
			if tc.Function.Arguments != "" {
				json.Unmarshal([]byte(tc.Function.Arguments), &args)
			}
			toolCalls = append(toolCalls, ToolCall{
				ID:        tc.ID,
				Name:      tc.Function.Name,
				Arguments: args,
			})
		}
	}

	var usage *UsageInfo
	if resp.Usage.TotalTokens > 0 {
		usage = &UsageInfo{
			PromptTokens:     int(resp.Usage.PromptTokens),
			CompletionTokens: int(resp.Usage.CompletionTokens),
			TotalTokens:      int(resp.Usage.TotalTokens),
		}
	}

	return &LLMResponse{
		Content:      message.Content,
		ToolCalls:    toolCalls,
		FinishReason: string(choice.FinishReason),
		Usage:        usage,
	}
}

func (p *CodexProvider) GetDefaultModel() string {
	return "gpt-4o"
}

func buildCodexParams(messages []Message, tools []ToolDefinition, model string, options map[string]interface{}) responses.ResponseNewParams {
	var inputItems responses.ResponseInputParam
	var instructions string

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			instructions = msg.Content
		case "user":
			if msg.ToolCallID != "" {
				inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
					OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
						CallID: msg.ToolCallID,
						Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{OfString: openai.Opt(msg.Content)},
					},
				})
			} else {
				inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
					OfMessage: &responses.EasyInputMessageParam{
						Role:    responses.EasyInputMessageRoleUser,
						Content: responses.EasyInputMessageContentUnionParam{OfString: openai.Opt(msg.Content)},
					},
				})
			}
		case "assistant":
			if len(msg.ToolCalls) > 0 {
				if msg.Content != "" {
					inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
						OfMessage: &responses.EasyInputMessageParam{
							Role:    responses.EasyInputMessageRoleAssistant,
							Content: responses.EasyInputMessageContentUnionParam{OfString: openai.Opt(msg.Content)},
						},
					})
				}
				for _, tc := range msg.ToolCalls {
					argsJSON, _ := json.Marshal(tc.Arguments)
					inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
						OfFunctionCall: &responses.ResponseFunctionToolCallParam{
							CallID:    tc.ID,
							Name:      tc.Name,
							Arguments: string(argsJSON),
						},
					})
				}
			} else {
				inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
					OfMessage: &responses.EasyInputMessageParam{
						Role:    responses.EasyInputMessageRoleAssistant,
						Content: responses.EasyInputMessageContentUnionParam{OfString: openai.Opt(msg.Content)},
					},
				})
			}
		case "tool":
			inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
				OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
					CallID: msg.ToolCallID,
					Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{OfString: openai.Opt(msg.Content)},
				},
			})
		}
	}

	params := responses.ResponseNewParams{
		Model: model,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: inputItems,
		},
		Store: openai.Opt(false),
	}

	if instructions != "" {
		params.Instructions = openai.Opt(instructions)
	} else {
		// ChatGPT Codex backend requires instructions to be present.
		params.Instructions = openai.Opt(defaultCodexInstructions)
	}

	if maxTokens, ok := options["max_tokens"].(int); ok {
		params.MaxOutputTokens = openai.Opt(int64(maxTokens))
	}

	if temp, ok := options["temperature"].(float64); ok {
		params.Temperature = openai.Opt(temp)
	}

	if len(tools) > 0 {
		params.Tools = translateToolsForCodex(tools)
	}

	return params
}

func translateToolsForCodex(tools []ToolDefinition) []responses.ToolUnionParam {
	result := make([]responses.ToolUnionParam, 0, len(tools))
	for _, t := range tools {
		ft := responses.FunctionToolParam{
			Name:       t.Function.Name,
			Parameters: t.Function.Parameters,
			Strict:     openai.Opt(false),
		}
		if t.Function.Description != "" {
			ft.Description = openai.Opt(t.Function.Description)
		}
		result = append(result, responses.ToolUnionParam{OfFunction: &ft})
	}
	return result
}

func parseCodexResponse(resp *responses.Response) *LLMResponse {
	var content strings.Builder
	var toolCalls []ToolCall

	for _, item := range resp.Output {
		switch item.Type {
		case "message":
			for _, c := range item.Content {
				if c.Type == "output_text" {
					content.WriteString(c.Text)
				}
			}
		case "function_call":
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(item.Arguments), &args); err != nil {
				args = map[string]interface{}{"raw": item.Arguments}
			}
			toolCalls = append(toolCalls, ToolCall{
				ID:        item.CallID,
				Name:      item.Name,
				Arguments: args,
			})
		}
	}

	finishReason := "stop"
	if len(toolCalls) > 0 {
		finishReason = "tool_calls"
	}
	if resp.Status == "incomplete" {
		finishReason = "length"
	}

	var usage *UsageInfo
	if resp.Usage.TotalTokens > 0 {
		usage = &UsageInfo{
			PromptTokens:     int(resp.Usage.InputTokens),
			CompletionTokens: int(resp.Usage.OutputTokens),
			TotalTokens:      int(resp.Usage.TotalTokens),
		}
	}

	return &LLMResponse{
		Content:      content.String(),
		ToolCalls:    toolCalls,
		FinishReason: finishReason,
		Usage:        usage,
	}
}

func createCodexTokenSource() func() (string, string, error) {
	return func() (string, string, error) {
		cred, err := auth.GetCredential("openai")
		if err != nil {
			return "", "", fmt.Errorf("loading auth credentials: %w", err)
		}
		if cred == nil {
			return "", "", fmt.Errorf("no credentials for openai. Run: picoclaw auth login --provider openai")
		}

		if cred.AuthMethod == "oauth" && cred.NeedsRefresh() && cred.RefreshToken != "" {
			oauthCfg := auth.OpenAIOAuthConfig()
			refreshed, err := auth.RefreshAccessToken(cred, oauthCfg)
			if err != nil {
				return "", "", fmt.Errorf("refreshing token: %w", err)
			}
			if err := auth.SetCredential("openai", refreshed); err != nil {
				return "", "", fmt.Errorf("saving refreshed token: %w", err)
			}
			return refreshed.AccessToken, refreshed.AccountID, nil
		}

		return cred.AccessToken, cred.AccountID, nil
	}
}

// LoadAzureConfigFromEnv loads Azure OpenAI configuration from environment variables
// Similar to azure-openai-models.ts getRequiredEnv() pattern
func LoadAzureConfigFromEnv() (*AzureConfig, error) {
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	deployment := os.Getenv("AZURE_OPENAI_DEPLOYMENT")
	apiVersion := os.Getenv("AZURE_OPENAI_API_VERSION")
	scope := os.Getenv("AZURE_OPENAI_SCOPE")
	managedIdentityID := os.Getenv("AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID")

	// Check if Azure config is present
	if endpoint == "" && deployment == "" && apiVersion == "" {
		return nil, nil // Not using Azure
	}

	// If any Azure config is present, all required fields must be set
	missing := []string{}
	if endpoint == "" {
		missing = append(missing, "AZURE_OPENAI_ENDPOINT")
	}
	if deployment == "" {
		missing = append(missing, "AZURE_OPENAI_DEPLOYMENT")
	}
	if apiVersion == "" {
		missing = append(missing, "AZURE_OPENAI_API_VERSION")
	}
	if scope == "" {
		missing = append(missing, "AZURE_OPENAI_SCOPE")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required Azure OpenAI environment variables: %v\nPlease set them in your .env file. See .env.example for reference", missing)
	}

	return &AzureConfig{
		Endpoint:           endpoint,
		Deployment:         deployment,
		APIVersion:         apiVersion,
		Scope:              scope,
		ManagedIdentityID:  managedIdentityID,
		UseManagedIdentity: true, // Always use Azure auth when Azure config is present
		Verbose:            os.Getenv("AZURE_OPENAI_VERBOSE") == "true",
	}, nil
}

// createAzureManagedIdentityTokenSource creates a token source using Azure Managed Identity
// This requires the Azure Identity SDK to be installed
func createAzureManagedIdentityTokenSource(config *AzureConfig) func() (string, string, error) {
	return func() (string, string, error) {
		if config == nil {
			return "", "", fmt.Errorf("Azure configuration is nil")
		}

		// NOTE: This is a placeholder implementation
		// To fully implement Azure Managed Identity, you need to:
		// 1. Add Azure Identity SDK: go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
		// 2. Add Azure Core SDK: go get github.com/Azure/azure-sdk-for-go/sdk/azcore
		// 3. Implement token retrieval using DefaultAzureCredential or ManagedIdentityCredential

		// Azure authentication using DefaultAzureCredential or ManagedIdentityCredential
		var cred azcore.TokenCredential
		var err error

		if config.ManagedIdentityID != "" {
			// User-assigned managed identity (for Azure deployment)
			if config.Verbose {
				fmt.Printf("[AzureAuth] Using userassigned managed identity: %s\n - codex_provider.go:540", config.ManagedIdentityID)
			}
			options := &azidentity.ManagedIdentityCredentialOptions{
				ID: azidentity.ClientID(config.ManagedIdentityID),
			}
			cred, err = azidentity.NewManagedIdentityCredential(options)
		} else {
			// DefaultAzureCredential supports multiple auth methods:
			// - Managed Identity (when running in Azure)
			// - Azure CLI (local testing with 'az login')
			// - Environment variables
			// - Interactive browser (if needed)
			if config.Verbose {
				fmt.Println("[AzureAuth] Using DefaultAzureCredential (supports local Azure CLI auth) - codex_provider.go:553")
			}
			cred, err = azidentity.NewDefaultAzureCredential(nil)
		}

		if err != nil {
			return "", "", fmt.Errorf("failed to create Azure credential: %w", err)
		}

		// Get access token for the specified scope
		token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
			Scopes: []string{config.Scope},
		})
		if err != nil {
			return "", "", fmt.Errorf("failed to get Azure access token: %w", err)
		}

		if config.Verbose {
			fmt.Printf("[AzureAuth] Retrieved token for scope: %s\n - codex_provider.go:571", config.Scope)
		}

		return token.Token, "", nil
	}
}

// createDynamicCodexTokenSource creates a token source with multiple authentication methods
// Priority: 1) Azure Managed Identity, 2) OAuth, 3) API Key
func createDynamicCodexTokenSource(azureConfig *AzureConfig) func() (string, string, error) {
	return func() (string, string, error) {
		// 1. Try Azure Managed Identity first (if configured)
		if azureConfig != nil && azureConfig.UseManagedIdentity {
			if azureConfig.Verbose {
				fmt.Println("[CodexProvider] Attempting Azure Managed Identity authentication - codex_provider.go:585")
			}
			tokenSource := createAzureManagedIdentityTokenSource(azureConfig)
			token, accountID, err := tokenSource()
			if err == nil && token != "" {
				if azureConfig.Verbose {
					fmt.Println("[CodexProvider] Successfully authenticated with Azure Managed Identity - codex_provider.go:591")
				}
				return token, accountID, nil
			}
			if azureConfig.Verbose {
				fmt.Printf("[CodexProvider] Azure Managed Identity failed: %v\n - codex_provider.go:596", err)
			}
		}

		// 2. Fallback to standard OpenAI authentication
		cred, err := auth.GetCredential("openai")
		if err != nil {
			return "", "", fmt.Errorf("loading auth credentials: %w", err)
		}
		if cred == nil {
			return "", "", fmt.Errorf("no credentials for openai. Run: picoclaw auth login --provider openai")
		}

		// 3. Try OAuth token refresh if needed
		if cred.AuthMethod == "oauth" && cred.NeedsRefresh() && cred.RefreshToken != "" {
			oauthCfg := auth.OpenAIOAuthConfig()
			refreshed, err := auth.RefreshAccessToken(cred, oauthCfg)
			if err != nil {
				return "", "", fmt.Errorf("refreshing token: %w", err)
			}
			if err := auth.SetCredential("openai", refreshed); err != nil {
				return "", "", fmt.Errorf("saving refreshed token: %w", err)
			}
			return refreshed.AccessToken, refreshed.AccountID, nil
		}

		// 4. Use existing token
		return cred.AccessToken, cred.AccountID, nil
	}
}
