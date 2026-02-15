# Claude Provider - Dynamic Token Management Implementation

## Summary

Successfully implemented dynamic token management in the Claude provider (`claude_provider.go`), based on the TypeScript `token-manager.ts` reference implementation.

## What Changed

### 1. New Data Structures

Added three new types to support token management:

```go
type TokenManagerConfig struct {
    Verbose bool    // Enable debug logging
    Account string  // macOS keychain account (optional)
}

type ClaudeCredentials struct {
    APIKey         string
    MCPOAuthTokens map[string]interface{}
    SessionToken   string
}
```

### 2. Enhanced ClaudeProvider

Updated the provider struct to include configuration:

```go
type ClaudeProvider struct {
    client      *anthropic.Client
    tokenSource func() (string, error)
    config      TokenManagerConfig  // NEW
}
```

### 3. New Constructor Functions

Added two new constructors for automatic token management:

- `NewClaudeProviderWithDynamicToken(config)` - Full control with config
- `NewClaudeProviderAuto()` - Simple auto-detection (recommended)

### 4. Token Retrieval Functions

Implemented complete token management system:

#### `createDynamicTokenSource(config)`
Main token source with 3-level fallback:
1. Environment variable `ANTHROPIC_API_KEY`
2. macOS keychain (if on macOS)
3. Auth package (existing mechanism)

#### `getClaudeCredentialsFromKeychain(config)`
Retrieves credentials from macOS keychain services:
- "Claude Code" - Primary API key location
- "Claude Code-credentials" - MCP OAuth tokens
- "Claude Safe Storage" - Session tokens

#### `getKeychainPassword(service, account)`
Low-level function to execute `security` command on macOS

#### `extractAPIKeyFromMCPCredentials(data)`
Parses MCP credentials JSON to extract API keys

## Key Features

### ✅ Multi-Source Token Support
- Environment variables (production)
- macOS keychain (local development)
- Auth package (existing users)

### ✅ Platform Detection
- Automatically detects macOS for keychain access
- Graceful fallback on Linux/Windows

### ✅ Verbose Logging
```go
provider, err := NewClaudeProviderWithDynamicToken(TokenManagerConfig{
    Verbose: true,
})
// Output:
// [TokenManager] Using ANTHROPIC_API_KEY from environment
// [TokenManager] Retrieved API key from keychain
// [TokenManager] Using credential from auth package
```

### ✅ Backward Compatible
- Existing code using `NewClaudeProvider(token)` continues to work
- Existing `NewClaudeProviderWithTokenSource()` unchanged

### ✅ Token Refresh
- Token automatically refreshed on each API call when using `tokenSource`
- No manual refresh required

## Usage Examples

### Simple Auto-Detection
```go
provider, err := providers.NewClaudeProviderAuto()
response, err := provider.Chat(ctx, messages, tools, model, options)
```

### With Verbose Logging
```go
provider, err := providers.NewClaudeProviderWithDynamicToken(
    providers.TokenManagerConfig{Verbose: true},
)
```

### With Specific Account
```go
provider, err := providers.NewClaudeProviderWithDynamicToken(
    providers.TokenManagerConfig{
        Account: "user@example.com",
        Verbose: true,
    },
)
```

## Files Created/Modified

### Modified
- ✏️ `claude_provider.go` - Core implementation with dynamic token management

### Created
- 📄 `CLAUDE_PROVIDER_USAGE.md` - Complete usage documentation
- 📄 `claude_provider_example.go` - Working code examples
- 📄 `IMPLEMENTATION_SUMMARY.md` - This file

## Testing

All code compiles successfully:
```bash
$ go build ./pkg/providers/
# Build successful!
```

## Comparison: TypeScript vs Go

| Feature | token-manager.ts | claude_provider.go |
|---------|-----------------|-------------------|
| Environment vars | ✅ `process.env` | ✅ `os.Getenv` |
| macOS Keychain | ✅ `execSync("security")` | ✅ `exec.Command("security")` |
| Platform detection | ✅ `process.platform` | ✅ `runtime.GOOS` |
| Verbose logging | ✅ `console.log` | ✅ `fmt.Println` |
| Multiple services | ✅ 3 services | ✅ 3 services |
| MCP OAuth tokens | ✅ Supported | ✅ Supported |
| Edge Runtime | ✅ Detected | N/A |
| Fallback chain | ✅ Full | ✅ Full |

## Migration Guide

### From Static Token
```go
// Before
provider := providers.NewClaudeProvider("sk-ant-...")

// After
provider, err := providers.NewClaudeProviderAuto()
if err != nil {
    log.Fatal(err)
}
```

### From Custom Token Source
```go
// Before
tokenSource := func() (string, error) {
    return auth.GetCredential("anthropic")
}
provider := providers.NewClaudeProviderWithTokenSource("", tokenSource)

// After - automatic fallback handles this
provider, err := providers.NewClaudeProviderAuto()
```

## Environment Setup

### Option 1: Environment Variable (Recommended for Production)
```bash
export ANTHROPIC_API_KEY="sk-ant-api03-..."
```

### Option 2: macOS Keychain (For Local Development)
If you have Claude Code installed, it's automatic. Otherwise:
```bash
security add-generic-password -s "Claude Code" -a "anthropic" -w "sk-ant-..."
```

### Option 3: Picoclaw Auth (Existing Method)
```bash
picoclaw auth login --provider anthropic
```

## Security Notes

- ✅ Tokens retrieved on-demand, not cached in memory
- ✅ Keychain access requires user authorization on macOS
- ✅ Platform detection prevents keychain errors on non-macOS
- ✅ Multiple fallback sources increase reliability
- ✅ Verbose mode for debugging (disabled by default)

## Next Steps

1. **Update existing code** to use `NewClaudeProviderAuto()`
2. **Set environment variable** in production: `ANTHROPIC_API_KEY`
3. **Test on macOS** with Claude Code keychain integration
4. **Enable verbose logging** during development to verify token sources

## Questions?

See:
- 📖 [CLAUDE_PROVIDER_USAGE.md](./CLAUDE_PROVIDER_USAGE.md) - Detailed usage guide
- 💻 [claude_provider_example.go](./claude_provider_example.go) - Working examples
- 🔍 [token-manager.ts](./token-manager.ts) - Original TypeScript reference
