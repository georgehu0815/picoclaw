package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/sipeed/picoclaw/pkg/auth"
)

// TokenManagerConfig holds configuration for token retrieval
type TokenManagerConfig struct {
	Verbose bool
	Account string
}

// ClaudeCredentials represents authentication credentials from various sources
type ClaudeCredentials struct {
	APIKey         string
	MCPOAuthTokens map[string]interface{}
	SessionToken   string
}

type ClaudeProvider struct {
	client      *anthropic.Client
	tokenSource func() (string, error)
	config      TokenManagerConfig
}

func NewClaudeProvider(token string) *ClaudeProvider {
	client := anthropic.NewClient(
		option.WithAPIKey(token),
		option.WithBaseURL("https://api.anthropic.com"),
	)
	return &ClaudeProvider{
		client: &client,
		config: TokenManagerConfig{Verbose: false},
	}
}

func NewClaudeProviderWithTokenSource(token string, tokenSource func() (string, error)) *ClaudeProvider {
	p := NewClaudeProvider(token)
	p.tokenSource = tokenSource
	return p
}

// NewClaudeProviderWithDynamicToken creates a provider with dynamic token management
// This implements the token-manager.ts functionality with automatic fallback:
// 1. Environment variable ANTHROPIC_API_KEY
// 2. macOS keychain (if on macOS)
// 3. Auth package credentials
func NewClaudeProviderWithDynamicToken(config TokenManagerConfig) (*ClaudeProvider, error) {
	// Create dynamic token source
	tokenSource := createDynamicTokenSource(config)

	// Get initial token
	token, err := tokenSource()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve API token: %w", err)
	}

	client := anthropic.NewClient(
		option.WithAPIKey(token),
		option.WithBaseURL("https://api.anthropic.com"),
	)

	return &ClaudeProvider{
		client:      &client,
		tokenSource: tokenSource,
		config:      config,
	}, nil
}

// NewClaudeProviderAuto creates a provider with automatic token detection
// Convenience function that uses default config
func NewClaudeProviderAuto() (*ClaudeProvider, error) {
	return NewClaudeProviderWithDynamicToken(TokenManagerConfig{
		Verbose: false,
	})
}

func (p *ClaudeProvider) Chat(ctx context.Context, messages []Message, tools []ToolDefinition, model string, options map[string]interface{}) (*LLMResponse, error) {
	var opts []option.RequestOption
	if p.tokenSource != nil {
		tok, err := p.tokenSource()
		if err != nil {
			return nil, fmt.Errorf("refreshing token: %w", err)
		}
		opts = append(opts, option.WithAPIKey(tok))
	}

	params, err := buildClaudeParams(messages, tools, model, options)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Messages.New(ctx, params, opts...)
	if err != nil {
		return nil, fmt.Errorf("claude API call: %w", err)
	}

	return parseClaudeResponse(resp), nil
}

func (p *ClaudeProvider) GetDefaultModel() string {
	return "claude-sonnet-4-5-20250929"
}

func buildClaudeParams(messages []Message, tools []ToolDefinition, model string, options map[string]interface{}) (anthropic.MessageNewParams, error) {
	var system []anthropic.TextBlockParam
	var anthropicMessages []anthropic.MessageParam

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			system = append(system, anthropic.TextBlockParam{Text: msg.Content})
		case "user":
			if msg.ToolCallID != "" {
				anthropicMessages = append(anthropicMessages,
					anthropic.NewUserMessage(anthropic.NewToolResultBlock(msg.ToolCallID, msg.Content, false)),
				)
			} else {
				anthropicMessages = append(anthropicMessages,
					anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content)),
				)
			}
		case "assistant":
			if len(msg.ToolCalls) > 0 {
				var blocks []anthropic.ContentBlockParamUnion
				if msg.Content != "" {
					blocks = append(blocks, anthropic.NewTextBlock(msg.Content))
				}
				for _, tc := range msg.ToolCalls {
					blocks = append(blocks, anthropic.NewToolUseBlock(tc.ID, tc.Arguments, tc.Name))
				}
				anthropicMessages = append(anthropicMessages, anthropic.NewAssistantMessage(blocks...))
			} else {
				anthropicMessages = append(anthropicMessages,
					anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content)),
				)
			}
		case "tool":
			anthropicMessages = append(anthropicMessages,
				anthropic.NewUserMessage(anthropic.NewToolResultBlock(msg.ToolCallID, msg.Content, false)),
			)
		}
	}

	maxTokens := int64(4096)
	if mt, ok := options["max_tokens"].(int); ok {
		maxTokens = int64(mt)
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		Messages:  anthropicMessages,
		MaxTokens: maxTokens,
	}

	if len(system) > 0 {
		params.System = system
	}

	if temp, ok := options["temperature"].(float64); ok {
		params.Temperature = anthropic.Float(temp)
	}

	if len(tools) > 0 {
		params.Tools = translateToolsForClaude(tools)
	}

	return params, nil
}

func translateToolsForClaude(tools []ToolDefinition) []anthropic.ToolUnionParam {
	result := make([]anthropic.ToolUnionParam, 0, len(tools))
	for _, t := range tools {
		tool := anthropic.ToolParam{
			Name: t.Function.Name,
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: t.Function.Parameters["properties"],
			},
		}
		if desc := t.Function.Description; desc != "" {
			tool.Description = anthropic.String(desc)
		}
		if req, ok := t.Function.Parameters["required"].([]interface{}); ok {
			required := make([]string, 0, len(req))
			for _, r := range req {
				if s, ok := r.(string); ok {
					required = append(required, s)
				}
			}
			tool.InputSchema.Required = required
		}
		result = append(result, anthropic.ToolUnionParam{OfTool: &tool})
	}
	return result
}

func parseClaudeResponse(resp *anthropic.Message) *LLMResponse {
	var content string
	var toolCalls []ToolCall

	for _, block := range resp.Content {
		switch block.Type {
		case "text":
			tb := block.AsText()
			content += tb.Text
		case "tool_use":
			tu := block.AsToolUse()
			var args map[string]interface{}
			if err := json.Unmarshal(tu.Input, &args); err != nil {
				args = map[string]interface{}{"raw": string(tu.Input)}
			}
			toolCalls = append(toolCalls, ToolCall{
				ID:        tu.ID,
				Name:      tu.Name,
				Arguments: args,
			})
		}
	}

	finishReason := "stop"
	switch resp.StopReason {
	case anthropic.StopReasonToolUse:
		finishReason = "tool_calls"
	case anthropic.StopReasonMaxTokens:
		finishReason = "length"
	case anthropic.StopReasonEndTurn:
		finishReason = "stop"
	}

	return &LLMResponse{
		Content:      content,
		ToolCalls:    toolCalls,
		FinishReason: finishReason,
		Usage: &UsageInfo{
			PromptTokens:     int(resp.Usage.InputTokens),
			CompletionTokens: int(resp.Usage.OutputTokens),
			TotalTokens:      int(resp.Usage.InputTokens + resp.Usage.OutputTokens),
		},
	}
}

func createClaudeTokenSource() func() (string, error) {
	return createDynamicTokenSource(TokenManagerConfig{Verbose: false})
}

// createDynamicTokenSource creates a token source with multiple fallback mechanisms
// Similar to token-manager.ts getAnthropicApiKey()
func createDynamicTokenSource(config TokenManagerConfig) func() (string, error) {
	return func() (string, error) {
		// 1. Try environment variable first (highest priority)
		if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
			if config.Verbose {
				fmt.Println("[TokenManager] Using ANTHROPIC_API_KEY from environment")
			}
			return apiKey, nil
		}

		// 2. Try macOS keychain (if on macOS)
		if runtime.GOOS == "darwin" {
			credentials := getClaudeCredentialsFromKeychain(config)
			if credentials.APIKey != "" {
				if config.Verbose {
					fmt.Println("[TokenManager] Retrieved API key from keychain")
				}
				return credentials.APIKey, nil
			}
		}

		// 3. Fallback to auth package (existing mechanism)
		cred, err := auth.GetCredential("anthropic")
		if err != nil {
			return "", fmt.Errorf("loading auth credentials: %w", err)
		}
		if cred == nil {
			return "", fmt.Errorf("no credentials for anthropic. Run: picoclaw auth login --provider anthropic")
		}

		if config.Verbose {
			fmt.Println("[TokenManager] Using credential from auth package")
		}
		return cred.AccessToken, nil
	}
}

// getClaudeCredentialsFromKeychain retrieves credentials from macOS keychain
// Similar to token-manager.ts getClaudeCredentials()
func getClaudeCredentialsFromKeychain(config TokenManagerConfig) ClaudeCredentials {
	credentials := ClaudeCredentials{
		MCPOAuthTokens: make(map[string]interface{}),
	}

	// Only attempt keychain access on macOS
	if runtime.GOOS != "darwin" {
		if config.Verbose {
			fmt.Println("[TokenManager] Not on macOS, skipping keychain access")
		}
		return credentials
	}

	// Try multiple keychain services in order of preference
	// This supports both Agency Claude and Claude Code
	keychainServices := []string{
		"Anthropic",   // Direct Anthropic API key
		"Agency",      // Agency Claude API key
		"Claude Code", // Claude Code API key
	}

	for _, service := range keychainServices {
		if apiKey := getKeychainPassword(service, config.Account); apiKey != "" {
			if strings.HasPrefix(apiKey, "sk-ant-") {
				if config.Verbose {
					fmt.Printf("[TokenManager] Found Anthropic API key in '%s' keychain service\n", service)
				}
				credentials.APIKey = apiKey
				return credentials
			} else if config.Verbose {
				fmt.Printf("[TokenManager] Found credential in '%s' but not a valid API key format (starts with: %s)\n", service, apiKey[:min(10, len(apiKey))])
			}
		}
	}

	// Try "Claude Code-credentials" (contains MCP OAuth tokens)
	if credsJSON := getKeychainPassword("Claude Code-credentials", config.Account); credsJSON != "" {
		var credsData map[string]interface{}
		if err := json.Unmarshal([]byte(credsJSON), &credsData); err == nil {
			// Extract MCP OAuth tokens
			if mcpOAuth, ok := credsData["mcpOAuth"].(map[string]interface{}); ok {
				credentials.MCPOAuthTokens = mcpOAuth
			}

			// Try to extract API key from credentials
			if apiKey := extractAPIKeyFromMCPCredentials(credsData); apiKey != "" {
				credentials.APIKey = apiKey
			}
		}
	}

	// Try "Claude Safe Storage" (encryption keys/session tokens)
	if safeStorage := getKeychainPassword("Claude Safe Storage", ""); safeStorage != "" {
		if credentials.APIKey == "" {
			credentials.SessionToken = safeStorage
		}
	}

	return credentials
}

// getKeychainPassword retrieves a password from macOS keychain
// Similar to token-manager.ts getKeychainPassword()
func getKeychainPassword(service, account string) string {
	// Build command
	args := []string{"find-generic-password", "-s", service, "-w"}
	if account != "" {
		args = append(args, "-a", account)
	}

	cmd := exec.Command("security", args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

// extractAPIKeyFromMCPCredentials extracts API key from MCP credentials structure
// Similar to token-manager.ts extractApiKeyFromMcpCredentials()
func extractAPIKeyFromMCPCredentials(data map[string]interface{}) string {
	// Check for direct API key
	if apiKey, ok := data["apiKey"].(string); ok {
		return apiKey
	}

	// Check for anthropic API key in various locations
	if anthropic, ok := data["anthropic"].(map[string]interface{}); ok {
		if apiKey, ok := anthropic["apiKey"].(string); ok {
			return apiKey
		}
	}

	return ""
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
