# PicoClaw - Command Line Guide

Complete guide for building and running PicoClaw from the command line.

## 📋 Prerequisites

- **Go**: Version 1.25.7 or later
- **Git**: For cloning the repository
- **macOS, Linux, or Windows**: Cross-platform support

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw
```

### 2. Build the Project

```bash
# Build the binary
go build -o picoclaw ./cmd/picoclaw/

# Verify the build
./picoclaw --version
```

### 3. Run Your First Command

```bash
# Simple question
./picoclaw agent -m "What is 2+2?"

# Expected output:
# 🦞 4
```

## 🔧 Building PicoClaw

### Standard Build

```bash
# Build for your current platform
go build -o picoclaw ./cmd/picoclaw/

# The binary will be created in the current directory
ls -lh picoclaw
```

### Optimized Build (Smaller Binary)

```bash
# Build with optimizations
go build -ldflags="-s -w" -o picoclaw ./cmd/picoclaw/

# Further compress with UPX (optional)
upx --best --lzma picoclaw
```

### Cross-Platform Builds

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o picoclaw-linux ./cmd/picoclaw/

# Build for macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o picoclaw-macos-intel ./cmd/picoclaw/

# Build for macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o picoclaw-macos-arm ./cmd/picoclaw/

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o picoclaw.exe ./cmd/picoclaw/
```

### Build with Version Information

```bash
# Set version info during build
go build \
  -ldflags="-X main.version=1.0.0 -X main.gitCommit=$(git rev-parse --short HEAD) -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o picoclaw ./cmd/picoclaw/

# Check version
./picoclaw --version
```

## 🏃 Running PicoClaw

### Using Go Run (Development)

```bash
# Run without building
go run ./cmd/picoclaw agent -m "Hello, Claude!"

# With debug logging
go run ./cmd/picoclaw agent -d -m "Debug mode test"
```

### Using Built Binary

```bash
# Run the built binary
./picoclaw agent -m "Your question here"

# Interactive mode
./picoclaw agent

# With session management
./picoclaw agent -s my-session
```

### Install System-Wide

```bash
# Install to GOPATH/bin
go install ./cmd/picoclaw

# Or copy to /usr/local/bin (requires sudo)
sudo cp picoclaw /usr/local/bin/

# Now run from anywhere
picoclaw agent -m "Hello from anywhere!"
```

## ⚙️ Configuration

### 1. Set Up Authentication

#### Option A: Environment Variable (Recommended)

```bash
# Set API key
export ANTHROPIC_API_KEY=sk-ant-api03-your-actual-key-here

# Test
picoclaw agent -m "Test message"
```

#### Option B: macOS Keychain (Automatic)

If you have Claude Code or Agency installed, PicoClaw will automatically use the keychain:

```bash
# Verify keychain has token
security find-generic-password -s "Claude Code" -w

# No additional setup needed!
picoclaw agent -m "Using keychain token"
```

#### Option C: Configuration File

```bash
# Create config directory
mkdir -p ~/.picoclaw

# Create configuration file
cat > ~/.picoclaw/config.json <<'EOF'
{
  "agents": {
    "defaults": {
      "workspace": "~/.picoclaw/workspace",
      "restrict_to_workspace": true,
      "provider": "anthropic",
      "model": "claude-sonnet-4-5-20250929",
      "max_tokens": 8192,
      "temperature": 0.7,
      "max_tool_iterations": 20
    }
  },
  "providers": {
    "anthropic": {
      "auth_method": "token"
    }
  }
}
EOF

# Test
picoclaw agent -m "Using config file"
```

### 2. Environment Variables

```bash
# Provider settings
export PICOCLAW_AGENTS_DEFAULTS_PROVIDER=anthropic
export PICOCLAW_AGENTS_DEFAULTS_MODEL=claude-sonnet-4-5-20250929

# Azure OpenAI (optional)
export AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com
export AZURE_OPENAI_DEPLOYMENT=gpt-4o
export AZURE_OPENAI_API_VERSION=2024-02-15-preview
export AZURE_OPENAI_SCOPE=https://cognitiveservices.azure.com/.default

# Verbose logging
export AZURE_OPENAI_VERBOSE=true
```

## 📖 Command Reference

### Agent Commands

```bash
# Single message
picoclaw agent -m "Your message here"

# Interactive mode
picoclaw agent

# With session
picoclaw agent -s session-name

# With debug logging
picoclaw agent -d -m "Debug message"

# Combined options
picoclaw agent -d -s my-session -m "Debug with session"
```

### Gateway Commands

```bash
# Start gateway server
picoclaw gateway

# Custom host and port
PICOCLAW_GATEWAY_HOST=0.0.0.0 PICOCLAW_GATEWAY_PORT=8080 picoclaw gateway
```

### Authentication Commands

```bash
# Login to provider
picoclaw auth login --provider anthropic

# Check auth status
picoclaw auth status

# Logout
picoclaw auth logout --provider anthropic
```

### Skills Commands

```bash
# List installed skills
picoclaw skills list

# Install a skill
picoclaw skills install

# Remove a skill
picoclaw skills remove <skill-name>

# List available built-in skills
picoclaw skills list-builtin

# Install built-in skills
picoclaw skills install-builtin
```

### Other Commands

```bash
# Show version
picoclaw --version

# Show help
picoclaw --help
picoclaw agent --help

# Run cron jobs
picoclaw cron

# Check status
picoclaw status

# Migration
picoclaw migrate
```

## 💡 Usage Examples

### Basic Questions

```bash
# Simple math
picoclaw agent -m "What is 15 * 23?"

# Code generation
picoclaw agent -m "Write a Python function to calculate fibonacci"

# Explanation
picoclaw agent -m "Explain how Docker works"
```

### Interactive Session

```bash
# Start interactive mode
picoclaw agent

# Then type your questions:
You: What is Go?
🦞 [Response about Go programming language]

You: Write a hello world in Go
🦞 [Code example]

# Exit with: exit, quit, or Ctrl+C
```

### Session Management

```bash
# Start named session
picoclaw agent -s learning-go

# Ask questions (context preserved)
picoclaw agent -s learning-go -m "What are goroutines?"
picoclaw agent -s learning-go -m "Show me an example"

# Each session maintains its own conversation history
```

### Debug Mode

```bash
# Enable debug logging
picoclaw agent -d -m "Test with debug info"

# Output will include detailed logs:
# [DEBUG] Token manager checking keychain...
# [DEBUG] Found API key in Claude Code keychain
# [INFO] Agent initialized...
```

## 🔍 Troubleshooting

### Build Errors

```bash
# If you get "pattern workspace: no matching files found"
cd cmd/picoclaw
cp -r ../../workspace .
cd ../..
go build -o picoclaw ./cmd/picoclaw/
```

### Authentication Issues

```bash
# Check if API key is set
echo $ANTHROPIC_API_KEY

# Check keychain (macOS)
security find-generic-password -s "Claude Code" -w

# Test with verbose logging
picoclaw agent -d -m "Test"
```

### Permission Errors

```bash
# Make binary executable
chmod +x picoclaw

# If copying to /usr/local/bin fails
sudo cp picoclaw /usr/local/bin/
```

## 🚀 Advanced Usage

### Using with Scripts

```bash
#!/bin/bash
# ask-claude.sh - Simple Claude wrapper

QUESTION="$*"
if [ -z "$QUESTION" ]; then
    echo "Usage: $0 <question>"
    exit 1
fi

picoclaw agent -m "$QUESTION"
```

### Piping Input

```bash
# Ask about code
cat myfile.go | picoclaw agent -m "Explain this code"

# Process output
picoclaw agent -m "List 5 random numbers" | grep -oE '[0-9]+'
```

### Background Processing

```bash
# Run in background
picoclaw gateway &

# Check if running
ps aux | grep picoclaw

# Stop
pkill picoclaw
```

### Using with Docker

```bash
# Create Dockerfile
cat > Dockerfile <<'EOF'
FROM golang:1.25-alpine
WORKDIR /app
COPY . .
RUN go build -o picoclaw ./cmd/picoclaw/
CMD ["./picoclaw", "agent"]
EOF

# Build
docker build -t picoclaw .

# Run
docker run -it -e ANTHROPIC_API_KEY=sk-ant-... picoclaw agent -m "Hello"
```

## 📊 Performance Tips

### Faster Builds

```bash
# Use build cache
go build -o picoclaw ./cmd/picoclaw/

# Parallel compilation (default)
go build -p 8 -o picoclaw ./cmd/picoclaw/

# Skip tests during development
go build -tags skiptest -o picoclaw ./cmd/picoclaw/
```

### Smaller Binaries

```bash
# Strip debug info
go build -ldflags="-s -w" -o picoclaw ./cmd/picoclaw/

# Check size
ls -lh picoclaw

# Compare with regular build
go build -o picoclaw-regular ./cmd/picoclaw/
ls -lh picoclaw*
```

## 🛠️ Development Workflow

### Build and Test Loop

```bash
# Watch for changes and rebuild (requires entr)
ls cmd/picoclaw/*.go pkg/**/*.go | entr -r go run ./cmd/picoclaw agent -m "Test"

# Manual loop
while true; do
    go build -o picoclaw ./cmd/picoclaw/
    ./picoclaw agent -m "Test"
    read -p "Press enter to rebuild..."
done
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run specific package
go test ./pkg/providers/

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Format and Lint

```bash
# Format code
go fmt ./...

# Run linter (if installed)
golangci-lint run

# Check for issues
go vet ./...
```

## 📦 Distribution

### Create Release Binary

```bash
# Build optimized binary
VERSION=1.0.0
go build \
  -ldflags="-s -w -X main.version=$VERSION" \
  -o picoclaw ./cmd/picoclaw/

# Create tarball
tar -czf picoclaw-$VERSION-$(uname -s)-$(uname -m).tar.gz picoclaw

# Create checksums
shasum -a 256 picoclaw-*.tar.gz > checksums.txt
```

### Multi-Platform Release

```bash
#!/bin/bash
# build-release.sh

VERSION=1.0.0
PLATFORMS="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64"

for PLATFORM in $PLATFORMS; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    OUTPUT=picoclaw-$VERSION-$GOOS-$GOARCH

    if [ $GOOS = "windows" ]; then
        OUTPUT+='.exe'
    fi

    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build \
      -ldflags="-s -w -X main.version=$VERSION" \
      -o $OUTPUT ./cmd/picoclaw/
done

echo "Done! Built $(ls picoclaw-* | wc -l) binaries"
```

## 🎯 Best Practices

1. **Always use version control for config files**
   ```bash
   git add .picoclaw/config.json
   ```

2. **Use environment variables in production**
   ```bash
   export ANTHROPIC_API_KEY=sk-ant-...
   ```

3. **Enable debug mode during development**
   ```bash
   alias picoclaw-debug='picoclaw agent -d'
   ```

4. **Keep your API keys secure**
   ```bash
   # Never commit .env files
   echo ".env" >> .gitignore
   ```

5. **Use named sessions for context**
   ```bash
   picoclaw agent -s project-name
   ```

## 📚 Additional Resources

- **Documentation**: See `pkg/providers/README_UPDATES.md`
- **Examples**: Check `pkg/providers/*_example.go`
- **Configuration**: See `.env.example`
- **Issues**: https://github.com/sipeed/picoclaw/issues

## ✨ Summary

**Quick Commands:**
```bash
# Build
go build -o picoclaw ./cmd/picoclaw/

# Run
./picoclaw agent -m "Your question"

# Interactive
./picoclaw agent

# Install
go install ./cmd/picoclaw
```

**That's it!** You're now ready to use PicoClaw from the command line. 🦞
