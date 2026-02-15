# Codex Provider - Azure OpenAI with Managed Identity Support

The Codex provider now supports Azure OpenAI with Managed Identity authentication, based on the TypeScript `azure-openai-models.ts` reference implementation.

## Features
      "provider": "codex",
      "model": "gpt-5.2-chat",
- **Azure OpenAI Integration** - Full support for Azure OpenAI endpoints
- **Managed Identity Authentication** - Both system-assigned and user-assigned managed identities
- **Multi-source Authentication** with automatic fallback:
  1. Azure Managed Identity (when configured)
  2. OpenAI OAuth (with automatic token refresh)
  3. API Key authentication
- **Environment-based Configuration** - Easy setup via environment variables
- **Automatic Provider Detection** - Auto-selects Azure or OpenAI based on configuration

## Azure Configuration

### Environment Variables

Based on `azure-openai-models.ts`, the following environment variables are required for Azure OpenAI:

```bash
# Required Azure OpenAI Configuration
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com
AZURE_OPENAI_DEPLOYMENT=your-deployment-name
AZURE_OPENAI_API_VERSION=2024-02-15-preview
AZURE_OPENAI_SCOPE=https://cognitiveservices.azure.com/.default

# Optional: User-Assigned Managed Identity
AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID=your-managed-identity-client-id

# Optional: Enable verbose logging
AZURE_OPENAI_VERBOSE=true
```

### Configuration Structure

```go
type AzureConfig struct {
    Endpoint           string // Azure OpenAI endpoint URL
    Deployment         string // Azure OpenAI deployment name
    APIVersion         string // Azure OpenAI API version
    Scope              string // Azure OpenAI scope for authentication
    ManagedIdentityID  string // Client ID for user-assigned managed identity
    UseManagedIdentity bool   // Enable managed identity authentication
    Verbose            bool   // Enable debug logging
}
```

## Usage Examples

### 1. Auto-Detection (Recommended)

The provider automatically detects Azure configuration from environment variables:

```go
import "github.com/sipeed/picoclaw/pkg/providers"

// Automatically detects Azure or OpenAI configuration
provider, err := providers.NewCodexProviderAuto()
if err != nil {
    log.Fatalf("Failed to create provider: %v", err)
}

// Use the provider (same API regardless of backend)
ctx := context.Background()
messages := []providers.Message{
    {Role: "user", Content: "Hello!"},
}

response, err := provider.Chat(ctx, messages, nil, "gpt-4o", nil)
if err != nil {
    log.Fatalf("Chat failed: %v", err)
}

fmt.Println("Response:", response.Content)
```

### 2. Explicit Azure Configuration

```go
// Load Azure config from environment
azureConfig, err := providers.LoadAzureConfigFromEnv()
if err != nil {
    log.Fatalf("Failed to load Azure config: %v", err)
}

// Create provider with Azure config
provider, err := providers.NewCodexProviderWithAzure(azureConfig, "")
if err != nil {
    log.Fatalf("Failed to create provider: %v", err)
}

// Use the provider
response, err := provider.Chat(ctx, messages, nil, "gpt-4o", nil)
```

### 3. Manual Azure Configuration

```go
// Create Azure config manually
azureConfig := &providers.AzureConfig{
    Endpoint:           "https://your-resource.openai.azure.com",
    Deployment:         "gpt-4o",
    APIVersion:         "2024-02-15-preview",
    Scope:              "https://cognitiveservices.azure.com/.default",
    ManagedIdentityID:  "your-client-id", // Optional
    UseManagedIdentity: true,
    Verbose:            true,
}

provider, err := providers.NewCodexProviderWithAzure(azureConfig, "")
if err != nil {
    log.Fatalf("Failed to create provider: %v", err)
}
```

### 4. With Verbose Logging

```bash
# Enable verbose logging
export AZURE_OPENAI_VERBOSE=true
```

```go
provider, err := providers.NewCodexProviderAuto()
// Output:
// [CodexProvider] Using Azure OpenAI configuration
// [CodexProvider] Attempting Azure Managed Identity authentication
// [AzureManagedIdentity] Retrieved token for scope: https://cognitiveservices.azure.com/.default
```

## Azure Managed Identity Setup

### Prerequisites

To use Azure Managed Identity authentication, you need to install the Azure SDK:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
go get github.com/Azure/azure-sdk-for-go/sdk/azcore
```

Then uncomment the implementation in `createAzureManagedIdentityTokenSource()` in [codex_provider.go](codex_provider.go).

### System-Assigned Managed Identity

1. Enable system-assigned managed identity on your Azure resource (VM, App Service, etc.)
2. Grant the managed identity access to your Azure OpenAI resource
3. Set environment variables (without `MANAGED_IDENTITY_CLIENT_ID`)

```bash
export AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com
export AZURE_OPENAI_DEPLOYMENT=gpt-4o
export AZURE_OPENAI_API_VERSION=2024-02-15-preview
export AZURE_OPENAI_SCOPE=https://cognitiveservices.azure.com/.default
```

### User-Assigned Managed Identity

1. Create a user-assigned managed identity in Azure
2. Grant it access to your Azure OpenAI resource
3. Set environment variables (including `MANAGED_IDENTITY_CLIENT_ID`)

```bash
export AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com
export AZURE_OPENAI_DEPLOYMENT=gpt-4o
export AZURE_OPENAI_API_VERSION=2024-02-15-preview
export AZURE_OPENAI_SCOPE=https://cognitiveservices.azure.com/.default
export AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID=your-client-id
```

## Authentication Priority

The provider uses the following authentication priority:

```
1. Azure Managed Identity (if AZURE_OPENAI_* env vars are set)
   ├─ User-Assigned (if MANAGED_IDENTITY_CLIENT_ID is set)
   └─ System-Assigned (default)

2. OpenAI OAuth (with automatic token refresh)
   └─ Uses refresh token if expired

3. API Key (from auth package)
   └─ Static API key authentication
```

## Comparison: TypeScript vs Go

| Feature | azure-openai-models.ts | codex_provider.go |
|---------|----------------------|-------------------|
| Environment config | ✅ `getRequiredEnv()` | ✅ `LoadAzureConfigFromEnv()` |
| Managed Identity | ✅ Implied | ✅ `createAzureManagedIdentityTokenSource()` |
| User-assigned MI | ✅ CLIENT_ID | ✅ `ManagedIdentityID` |
| System-assigned MI | ✅ Default | ✅ Default |
| Scope configuration | ✅ `AZURE_OPENAI_SCOPE` | ✅ `Scope` field |
| Verbose logging | ❌ | ✅ `Verbose` flag |
| Auto-detection | ❌ | ✅ `NewCodexProviderAuto()` |

## Environment File Example

Create a `.env` file in your project root:

```bash
# Azure OpenAI Configuration (required for Azure)
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com
AZURE_OPENAI_DEPLOYMENT=gpt-4o
AZURE_OPENAI_API_VERSION=2024-02-15-preview
AZURE_OPENAI_SCOPE=https://cognitiveservices.azure.com/.default

# Optional: User-Assigned Managed Identity
AZURE_OPENAI_MANAGED_IDENTITY_CLIENT_ID=12345678-1234-1234-1234-123456789abc

# Optional: Enable verbose logging
AZURE_OPENAI_VERBOSE=true
```

## Migration Guide

### From Standard OpenAI to Azure OpenAI

```go
// Before (OpenAI)
provider := providers.NewCodexProvider(token, accountID)

// After (Azure OpenAI with auto-detection)
provider, err := providers.NewCodexProviderAuto()
if err != nil {
    log.Fatal(err)
}
```

### From Static Token to Managed Identity

```go
// Before (static token)
provider := providers.NewCodexProvider("your-api-key", "")

// After (managed identity)
// 1. Set environment variables
// 2. Use auto-detection
provider, err := providers.NewCodexProviderAuto()
```

## Error Handling

### Missing Azure SDK

If you see this error:

```
Azure Managed Identity support requires Azure SDK packages. Please add:
  go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
  go get github.com/Azure/azure-sdk-for-go/sdk/azcore
```

Solution:
1. Install the Azure SDK packages
2. Uncomment the implementation in `createAzureManagedIdentityTokenSource()`

### Missing Environment Variables

If you see this error:

```
missing required Azure OpenAI environment variables: [AZURE_OPENAI_ENDPOINT AZURE_OPENAI_DEPLOYMENT]
Please set them in your .env file. See .env.example for reference
```

Solution: Set all required environment variables in your `.env` file.

### Authentication Failed

If managed identity authentication fails, the provider automatically falls back to:
1. OpenAI OAuth (if available)
2. API Key authentication (from auth package)

Enable verbose logging to see the fallback chain:

```bash
export AZURE_OPENAI_VERBOSE=true
```

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/sipeed/picoclaw/pkg/providers"
)

func main() {
    // Auto-detect Azure or OpenAI configuration
    provider, err := providers.NewCodexProviderAuto()
    if err != nil {
        log.Fatalf("Failed to create provider: %v", err)
    }

    // Prepare messages
    ctx := context.Background()
    messages := []providers.Message{
        {
            Role:    "system",
            Content: "You are a helpful coding assistant.",
        },
        {
            Role:    "user",
            Content: "Write a function to calculate fibonacci numbers in Go",
        },
    }

    // Call the API
    response, err := provider.Chat(ctx, messages, nil, "gpt-4o", map[string]interface{}{
        "max_tokens":  1000,
        "temperature": 0.7,
    })
    if err != nil {
        log.Fatalf("Chat failed: %v", err)
    }

    // Print response
    fmt.Println("Response:", response.Content)
    fmt.Printf("Usage: %d input tokens, %d output tokens\n",
        response.Usage.PromptTokens,
        response.Usage.CompletionTokens,
    )
}
```

## Security Best Practices

1. **Use Managed Identity** in Azure environments (most secure)
2. **Never commit** `.env` files with credentials to version control
3. **Rotate tokens** regularly when using API keys
4. **Use RBAC** to grant minimal permissions to managed identities
5. **Enable verbose logging** only in development, not production
6. **Use system-assigned MI** when possible (simpler management)

## Troubleshooting

### Check Configuration

```go
// Load and validate Azure config
azureConfig, err := providers.LoadAzureConfigFromEnv()
if err != nil {
    fmt.Printf("Error: %v\n", err)
}
if azureConfig != nil {
    fmt.Printf("Endpoint: %s\n", azureConfig.Endpoint)
    fmt.Printf("Deployment: %s\n", azureConfig.Deployment)
    fmt.Printf("Managed Identity: %v\n", azureConfig.UseManagedIdentity)
}
```

### Test Managed Identity

Run on an Azure VM or App Service with managed identity enabled:

```bash
# Set verbose logging
export AZURE_OPENAI_VERBOSE=true

# Run your application
go run main.go
```

Expected output:
```
[CodexProvider] Using Azure OpenAI configuration
[CodexProvider] Attempting Azure Managed Identity authentication
[AzureManagedIdentity] Retrieved token for scope: https://cognitiveservices.azure.com/.default
[CodexProvider] Successfully authenticated with Azure Managed Identity
```

## Reference

- TypeScript implementation: [azure-openai-models.ts](azure-openai-models.ts)
- Go implementation: [codex_provider.go](codex_provider.go)
- Azure SDK for Go: https://github.com/Azure/azure-sdk-for-go
- Azure Managed Identity: https://docs.microsoft.com/azure/active-directory/managed-identities-azure-resources/
