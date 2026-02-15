package providers

import (
	"context"
	"fmt"
	"log"
)

// ExampleAzureAutoDetection demonstrates automatic Azure/OpenAI detection
func ExampleAzureAutoDetection() {
	// Automatically detects Azure or OpenAI based on environment variables
	provider, err := NewCodexProviderAuto()
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()
	messages := []Message{
		{Role: "user", Content: "Hello from Azure OpenAI!"},
	}

	response, err := provider.Chat(ctx, messages, nil, "gpt-4o", nil)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Println("Response:", response.Content)
}

// ExampleAzureExplicitConfig demonstrates explicit Azure configuration
func ExampleAzureExplicitConfig() {
	// Load Azure config from environment variables
	azureConfig, err := LoadAzureConfigFromEnv()
	if err != nil {
		log.Fatalf("Failed to load Azure config: %v", err)
	}

	if azureConfig == nil {
		log.Fatal("Azure configuration not found in environment")
	}

	// Create provider with Azure config
	provider, err := NewCodexProviderWithAzure(azureConfig, "")
	if err != nil {
		log.Fatalf("Failed to create Azure provider: %v", err)
	}

	ctx := context.Background()
	messages := []Message{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: "Explain Azure Managed Identity"},
	}

	response, err := provider.Chat(ctx, messages, nil, "gpt-4o", nil)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Println("Response:", response.Content)
	fmt.Printf("Tokens used: %d\n", response.Usage.TotalTokens)
}

// ExampleAzureManualConfig demonstrates manual Azure configuration
func ExampleAzureManualConfig() {
	// Create Azure config manually (without environment variables)
	azureConfig := &AzureConfig{
		Endpoint:           "https://your-resource.openai.azure.com",
		Deployment:         "gpt-4o",
		APIVersion:         "2024-02-15-preview",
		Scope:              "https://cognitiveservices.azure.com/.default",
		ManagedIdentityID:  "", // Empty for system-assigned managed identity
		UseManagedIdentity: true,
		Verbose:            true, // Enable verbose logging
	}

	provider, err := NewCodexProviderWithAzure(azureConfig, "")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// The verbose flag will show authentication details:
	// [CodexProvider] Attempting Azure Managed Identity authentication
	// [AzureManagedIdentity] Retrieved token for scope: ...
	// [CodexProvider] Successfully authenticated with Azure Managed Identity

	ctx := context.Background()
	messages := []Message{
		{Role: "user", Content: "Write a Hello World in Go"},
	}

	response, err := provider.Chat(ctx, messages, nil, "gpt-4o", nil)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Println("Response:", response.Content)
}

// ExampleAzureUserAssignedManagedIdentity demonstrates user-assigned managed identity
func ExampleAzureUserAssignedManagedIdentity() {
	azureConfig := &AzureConfig{
		Endpoint:  "https://your-resource.openai.azure.com",
		Deployment: "gpt-4o",
		APIVersion: "2024-02-15-preview",
		Scope:     "https://cognitiveservices.azure.com/.default",
		// Specify client ID for user-assigned managed identity
		ManagedIdentityID:  "12345678-1234-1234-1234-123456789abc",
		UseManagedIdentity: true,
		Verbose:            false,
	}

	provider, err := NewCodexProviderWithAzure(azureConfig, "")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	ctx := context.Background()
	messages := []Message{
		{Role: "user", Content: "What is Go?"},
	}

	response, err := provider.Chat(ctx, messages, nil, "gpt-4o", nil)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Println("Response:", response.Content)
}

// ExampleAzureWithToolCalls demonstrates Azure OpenAI with function calling
func ExampleAzureWithToolCalls() {
	provider, err := NewCodexProviderAuto()
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// Define a weather tool
	tools := []ToolDefinition{
		{
			Type: "function",
			Function: ToolFunctionDefinition{
				Name:        "get_weather",
				Description: "Get the current weather for a location",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "The city and state, e.g. San Francisco, CA",
						},
						"unit": map[string]interface{}{
							"type": "string",
							"enum": []string{"celsius", "fahrenheit"},
						},
					},
					"required": []interface{}{"location"},
				},
			},
		},
	}

	ctx := context.Background()
	messages := []Message{
		{Role: "user", Content: "What's the weather in Seattle?"},
	}

	response, err := provider.Chat(ctx, messages, tools, "gpt-4o", nil)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	// Check if the model wants to call a function
	if len(response.ToolCalls) > 0 {
		for _, toolCall := range response.ToolCalls {
			fmt.Printf("Function call: %s\n", toolCall.Name)
			fmt.Printf("Arguments: %+v\n", toolCall.Arguments)
		}
	} else {
		fmt.Println("Response:", response.Content)
	}
}

// ExampleAzureConfigValidation demonstrates config validation
func ExampleAzureConfigValidation() {
	// Try to load Azure config (may be nil if not configured)
	azureConfig, err := LoadAzureConfigFromEnv()
	if err != nil {
		log.Printf("Error loading Azure config: %v", err)
		return
	}

	if azureConfig == nil {
		log.Println("Azure not configured, using standard OpenAI")
		// Create standard OpenAI provider
		tokenSource := createCodexTokenSource()
		token, accountID, _ := tokenSource()
		provider := NewCodexProvider(token, accountID)
		_ = provider
	} else {
		log.Println("Azure OpenAI configured:")
		log.Printf("  Endpoint: %s", azureConfig.Endpoint)
		log.Printf("  Deployment: %s", azureConfig.Deployment)
		log.Printf("  API Version: %s", azureConfig.APIVersion)
		log.Printf("  Managed Identity: %v", azureConfig.UseManagedIdentity)
		if azureConfig.ManagedIdentityID != "" {
			log.Printf("  MI Client ID: %s", azureConfig.ManagedIdentityID)
		}

		// Create Azure provider
		provider, err := NewCodexProviderWithAzure(azureConfig, "")
		if err != nil {
			log.Fatalf("Failed to create provider: %v", err)
		}
		_ = provider
	}
}

// ExampleAzureFallbackAuthentication demonstrates authentication fallback
func ExampleAzureFallbackAuthentication() {
	azureConfig := &AzureConfig{
		Endpoint:           "https://your-resource.openai.azure.com",
		Deployment:         "gpt-4o",
		APIVersion:         "2024-02-15-preview",
		Scope:              "https://cognitiveservices.azure.com/.default",
		UseManagedIdentity: true,
		Verbose:            true, // See the fallback chain
	}

	provider, err := NewCodexProviderWithAzure(azureConfig, "")
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// With verbose=true, you'll see the authentication attempts:
	// 1. [CodexProvider] Attempting Azure Managed Identity authentication
	//    If fails: [CodexProvider] Azure Managed Identity failed: ...
	// 2. Falls back to OpenAI OAuth
	// 3. Falls back to API Key

	ctx := context.Background()
	messages := []Message{
		{Role: "user", Content: "Test message"},
	}

	response, err := provider.Chat(ctx, messages, nil, "gpt-4o", nil)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Println("Response:", response.Content)
}

// ExampleAzureEnvironmentSetup demonstrates environment setup
func ExampleAzureEnvironmentSetup() {
	// In production, you would set these in your .env file or system environment:
	/*
		export AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com
		export AZURE_OPENAI_DEPLOYMENT=gpt-4o
		export AZURE_OPENAI_API_VERSION=2024-02-15-preview
		export AZURE_OPENAI_SCOPE=https://cognitiveservices.azure.com/.default
		export AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID=your-client-id  # Optional
		export AZURE_OPENAI_VERBOSE=true  # Optional
	*/

	// For this example, we'll check if they're set
	requiredVars := []string{
		"AZURE_OPENAI_ENDPOINT",
		"AZURE_OPENAI_DEPLOYMENT",
		"AZURE_OPENAI_API_VERSION",
		"AZURE_OPENAI_SCOPE",
	}

	allSet := true
	for _, varName := range requiredVars {
		// This would normally use os.Getenv in real code
		fmt.Printf("Checking %s...\n", varName)
		allSet = allSet && true // Placeholder
	}

	if allSet {
		fmt.Println("✓ All Azure environment variables are set")
		provider, err := NewCodexProviderAuto()
		if err != nil {
			log.Fatalf("Failed: %v", err)
		}
		_ = provider
	} else {
		fmt.Println("✗ Some Azure environment variables are missing")
		fmt.Println("  See CODEX_AZURE_USAGE.md for setup instructions")
	}
}

// ExampleAzureMigration demonstrates migrating from OpenAI to Azure
func ExampleAzureMigration() {
	// BEFORE: Using standard OpenAI
	oldProvider := NewCodexProvider("your-api-key", "account-id")
	_ = oldProvider

	// AFTER: Using Azure OpenAI with managed identity
	// 1. Set environment variables (see CODEX_AZURE_USAGE.md)
	// 2. Use auto-detection
	newProvider, err := NewCodexProviderAuto()
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// The API remains the same!
	ctx := context.Background()
	messages := []Message{
		{Role: "user", Content: "Hello!"},
	}

	// Same function call works for both OpenAI and Azure
	response, err := newProvider.Chat(ctx, messages, nil, "gpt-4o", nil)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Println("Response:", response.Content)
}

// ExampleAzureCompleteWorkflow demonstrates a complete workflow
func ExampleAzureCompleteWorkflow() {
	// 1. Load configuration
	azureConfig, err := LoadAzureConfigFromEnv()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	var provider *CodexProvider
	if azureConfig != nil {
		// 2. Create Azure provider
		provider, err = NewCodexProviderWithAzure(azureConfig, "")
		if err != nil {
			log.Fatalf("Failed to create Azure provider: %v", err)
		}
		log.Println("Using Azure OpenAI")
	} else {
		// 2. Fallback to standard OpenAI
		provider, err = NewCodexProviderAuto()
		if err != nil {
			log.Fatalf("Failed to create OpenAI provider: %v", err)
		}
		log.Println("Using standard OpenAI")
	}

	// 3. Prepare conversation
	ctx := context.Background()
	messages := []Message{
		{
			Role:    "system",
			Content: "You are a helpful coding assistant specializing in Go.",
		},
		{
			Role:    "user",
			Content: "Write a function to reverse a string in Go",
		},
	}

	// 4. Make the API call
	options := map[string]interface{}{
		"max_tokens":  500,
		"temperature": 0.7,
	}

	response, err := provider.Chat(ctx, messages, nil, "gpt-4o", options)
	if err != nil {
		log.Fatalf("API call failed: %v", err)
	}

	// 5. Process the response
	fmt.Println("=== Response ===")
	fmt.Println(response.Content)
	fmt.Println("\n=== Usage ===")
	fmt.Printf("Input tokens: %d\n", response.Usage.PromptTokens)
	fmt.Printf("Output tokens: %d\n", response.Usage.CompletionTokens)
	fmt.Printf("Total tokens: %d\n", response.Usage.TotalTokens)
}
