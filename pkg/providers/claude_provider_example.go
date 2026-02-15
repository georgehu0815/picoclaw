package providers

import (
	"context"
	"fmt"
	"log"
)

// ExampleBasicUsage demonstrates the simplest way to use the Claude provider
func ExampleBasicUsage() {
	// Auto-detect token from environment, keychain, or auth package
	provider, err := NewClaudeProviderAuto()
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// Use the provider
	ctx := context.Background()
	messages := []Message{
		{Role: "user", Content: "Hello, Claude!"},
	}

	response, err := provider.Chat(ctx, messages, nil, provider.GetDefaultModel(), nil)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Println("Response: - claude_provider_example.go:28", response.Content)
}

// ExampleVerboseLogging shows how to enable debug logging
func ExampleVerboseLogging() {
	provider, err := NewClaudeProviderWithDynamicToken(TokenManagerConfig{
		Verbose: true, // Enable verbose logging
	})
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// You'll see logs like:
	// [TokenManager] Using ANTHROPIC_API_KEY from environment
	// or
	// [TokenManager] Retrieved API key from keychain
	// or
	// [TokenManager] Using credential from auth package

	ctx := context.Background()
	messages := []Message{
		{Role: "user", Content: "Explain quantum computing"},
	}

	response, err := provider.Chat(ctx, messages, nil, "claude-sonnet-4-5-20250929", nil)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Println("Response: - claude_provider_example.go:57", response.Content)
}

// ExampleWithToolCalls demonstrates using Claude with function calling
func ExampleWithToolCalls() {
	provider, err := NewClaudeProviderAuto()
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// Define a tool
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
					},
					"required": []interface{}{"location"},
				},
			},
		},
	}

	ctx := context.Background()
	messages := []Message{
		{Role: "user", Content: "What's the weather in San Francisco?"},
	}

	response, err := provider.Chat(ctx, messages, tools, provider.GetDefaultModel(), nil)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	// Check if Claude wants to use a tool
	if len(response.ToolCalls) > 0 {
		for _, toolCall := range response.ToolCalls {
			fmt.Printf("Tool call: %s(%v)\n - claude_provider_example.go:101", toolCall.Name, toolCall.Arguments)
		}
	} else {
		fmt.Println("Response: - claude_provider_example.go:104", response.Content)
	}
}

// ExampleCustomTokenSource demonstrates creating a custom token source
func ExampleCustomTokenSource() {
	// Custom token source that could rotate between multiple keys,
	// refresh from a secret manager, etc.
	customTokenSource := func() (string, error) {
		// Example: Load from custom secret manager
		// token, err := yourSecretManager.GetToken("anthropic")
		// return token, err

		// For this example, just return from environment
		return "sk-ant-api03-...", nil
	}

	provider := NewClaudeProviderWithTokenSource("initial-token", customTokenSource)

	ctx := context.Background()
	messages := []Message{
		{Role: "user", Content: "Hello!"},
	}

	response, err := provider.Chat(ctx, messages, nil, provider.GetDefaultModel(), nil)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Println("Response: - claude_provider_example.go:133", response.Content)
}

// ExampleTokenSourcePatterns shows different token source patterns
func ExampleTokenSourcePatterns() {
	// Pattern 1: Environment-only (for production)
	envProvider, _ := NewClaudeProviderWithDynamicToken(TokenManagerConfig{})

	// Pattern 2: Keychain-first (for local development on macOS)
	keychainProvider, _ := NewClaudeProviderWithDynamicToken(TokenManagerConfig{
		Verbose: true, // See where token comes from
	})

	// Pattern 3: Custom rotation
	rotatingTokenSource := func() (string, error) {
		// Implement your token rotation logic
		// This gets called on every API request
		tokens := []string{"token1", "token2", "token3"}
		// Pick one based on your logic (round-robin, load balancing, etc.)
		return tokens[0], nil
	}
	rotatingProvider := NewClaudeProviderWithTokenSource("", rotatingTokenSource)

	// Use any provider
	_ = envProvider
	_ = keychainProvider
	_ = rotatingProvider
}

// ExampleMigrationFromStatic shows migrating from static to dynamic tokens
func ExampleMigrationFromStatic() {
	// OLD WAY - static token
	// oldProvider := NewClaudeProvider("sk-ant-api03-...")

	// NEW WAY - dynamic token with automatic refresh
	newProvider, err := NewClaudeProviderAuto()
	if err != nil {
		log.Fatalf("Failed to create provider: %v", err)
	}

	// The API is exactly the same, but now token is managed automatically
	ctx := context.Background()
	messages := []Message{
		{Role: "user", Content: "Hello!"},
	}

	_, err = newProvider.Chat(ctx, messages, nil, newProvider.GetDefaultModel(), nil)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}
}
