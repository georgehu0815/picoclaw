# Claude Provider - Dynamic Token Management

The Claude provider now supports dynamic token management with automatic fallback, similar to the TypeScript `token-manager.ts` implementation.

## Features
      "provider": "anthropic",
      "model": "claude-sonnet-4-5-20250929",
- **Multi-source token retrieval** with automatic fallback
- **macOS Keychain integration** (when running on macOS)
- **Environment variable support** (ANTHROPIC_API_KEY)
- **Existing auth package integration** as final fallback
- **Verbose logging** for debugging token retrieval

## Token Retrieval Priority

The provider checks for tokens in this order:

1. **Environment Variable** (`ANTHROPIC_API_KEY`) - Highest priority
2. **macOS Keychain** (if on macOS)
   - "Claude Code" service (contains API key)
   - "Claude Code-credentials" (MCP OAuth tokens)
   - "Claude Safe Storage" (session tokens)
3. **Auth Package** - Existing `picoclaw auth login` mechanism

## Usage Examples

### 1. Auto-Detection (Recommended)

```go
import "github.com/sipeed/picoclaw/pkg/providers"

// Automatically detects and uses available token source
provider, err := providers.NewClaudeProviderAuto()
if err != nil {
    log.Fatalf("Failed to create provider: %v", err)
}

// Use the provider
response, err := provider.Chat(ctx, messages, tools, model, options)
```

### 2. With Verbose Logging

```go
// Enable verbose logging to see where tokens are coming from
provider, err := providers.NewClaudeProviderWithDynamicToken(
    providers.TokenManagerConfig{
        Verbose: true,  // Enable debug logging
    },
)
if err != nil {
    log.Fatalf("Failed to create provider: %v", err)
}
```

### 3. With Specific macOS Account

```go
// Specify a macOS keychain account
provider, err := providers.NewClaudeProviderWithDynamicToken(
    providers.TokenManagerConfig{
        Verbose: true,
        Account: "user@example.com",  // Specific keychain account
    },
)
```

### 4. Custom Token Source (Advanced)

```go
// Create a custom token source function
customTokenSource := func() (string, error) {
    // Your custom token retrieval logic
    return "sk-ant-...", nil
}

provider := providers.NewClaudeProviderWithTokenSource("initial-token", customTokenSource)
```

### 5. Legacy Method (Static Token)

```go
// Use a static token (no automatic refresh)
provider := providers.NewClaudeProvider("sk-ant-api03-...")
```

## Environment Setup

### Option 1: Environment Variable

```bash
export ANTHROPIC_API_KEY="sk-ant-api03-..."
```

### Option 2: macOS Keychain (Claude Code)

If you have Claude Code installed, the provider will automatically use the API key stored in your macOS keychain:

```bash
# The provider looks for these keychain services:
# - "Claude Code" (API key)
# - "Claude Code-credentials" (MCP OAuth tokens)
# - "Claude Safe Storage" (session tokens)
```

### Option 3: Picoclaw Auth

```bash
picoclaw auth login --provider anthropic
```

## Token Refresh

When using `tokenSource`, the provider automatically refreshes the token on each API call:

```go
func (p *ClaudeProvider) Chat(ctx context.Context, ...) (*LLMResponse, error) {
    // Token is automatically refreshed here if tokenSource is set
    if p.tokenSource != nil {
        tok, err := p.tokenSource()
        if err != nil {
            return nil, fmt.Errorf("refreshing token: %w", err)
        }
        // Use refreshed token
    }
    // ... make API call
}
```

## Comparison with TypeScript Implementation

The Go implementation mirrors the TypeScript `token-manager.ts`:

| Feature | TypeScript | Go |
|---------|-----------|-----|
| Environment variable | `process.env.ANTHROPIC_API_KEY` | `os.Getenv("ANTHROPIC_API_KEY")` |
| macOS Keychain | `security find-generic-password` | `exec.Command("security", ...)` |
| Platform detection | `process.platform === 'darwin'` | `runtime.GOOS == "darwin"` |
| Verbose logging | `console.log('[TokenManager]')` | `fmt.Println("[TokenManager]")` |
| Fallback chain | ✅ | ✅ |
| Edge Runtime support | ✅ | N/A (Go native) |

## Error Handling

```go
provider, err := providers.NewClaudeProviderAuto()
if err != nil {
    // Possible errors:
    // - No token found in any source
    // - Keychain access denied (on macOS)
    // - Invalid token format
    // - Auth package errors
    log.Fatalf("Token retrieval failed: %v", err)
}
```

## Best Practices

1. **Use `NewClaudeProviderAuto()`** for most cases - it handles everything automatically
2. **Enable verbose logging** during development to understand token sources
3. **Set environment variable** in production for better security and portability
4. **Use keychain** for local development on macOS (seamless integration with Claude Code)
5. **Implement retry logic** in your application for production use

## Migration Guide

### From Old Code

```go
// Old way
provider := providers.NewClaudeProvider("sk-ant-...")
```

### To New Code

```go
// New way - automatic token management
provider, err := providers.NewClaudeProviderAuto()
if err != nil {
    log.Fatal(err)
}
```

## Security Considerations

- **Keychain Access**: Only works on macOS, gracefully falls back on other platforms
- **Token Exposure**: Tokens are retrieved on-demand, not stored in memory longer than necessary
- **Environment Variables**: Standard method for production deployments
- **No Token Caching**: Fresh token retrieved on each API call when using `tokenSource`

## Troubleshooting

### "No credentials found"

Enable verbose mode to see what's happening:

```go
provider, err := providers.NewClaudeProviderWithDynamicToken(
    providers.TokenManagerConfig{Verbose: true},
)
```

### Keychain Access Issues on macOS

```bash
# Verify keychain contains the API key
security find-generic-password -s "Claude Code" -w

# If empty, set it manually
security add-generic-password -s "Claude Code" -a "anthropic" -w "sk-ant-..."
```

### Cross-Platform Compatibility

The provider automatically detects the platform and only attempts keychain access on macOS. On Linux/Windows, it uses environment variables or the auth package.
