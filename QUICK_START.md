# PicoClaw - Quick Start

## 🚀 Get Started in 30 Seconds

```bash
# 1. Build
go build -o picoclaw ./cmd/picoclaw/

# 2. Set API key (choose one method)
export ANTHROPIC_API_KEY=sk-ant-api03-your-key-here
# OR let it auto-detect from keychain (macOS with Claude Code)

# 3. Run
./picoclaw agent -m "What is 2+2?"
```

## 📖 Common Commands

```bash
# Single question
picoclaw agent -m "Your question here"

# Interactive mode
picoclaw agent

# With debug logging
picoclaw agent -d -m "Test"

# Named session
picoclaw agent -s my-session

# Show version
picoclaw --version

# Show help
picoclaw --help
```

## ⚙️ Setup Authentication

### Method 1: Environment Variable (Recommended)
```bash
export ANTHROPIC_API_KEY=sk-ant-api03-your-actual-key
```

### Method 2: macOS Keychain (Automatic)
- If you have Claude Code installed, it works automatically!
- No setup needed

### Method 3: Configuration File
```bash
mkdir -p ~/.picoclaw
cat > ~/.picoclaw/config.json <<'EOF'
{
  "agents": {
    "defaults": {
      "provider": "anthropic",
      "model": "claude-sonnet-4-5-20250929"
    }
  },
  "providers": {
    "anthropic": {
      "auth_method": "token"
    }
  }
}
EOF
```

## 🎯 Examples

```bash
# Math
picoclaw agent -m "What is 15 * 23?"

# Code
picoclaw agent -m "Write a Go hello world"

# Explanation
picoclaw agent -m "Explain Docker in simple terms"

# Debug mode
picoclaw agent -d -m "Show me debug output"
```

## 🔧 Build Options

```bash
# Standard build
go build -o picoclaw ./cmd/picoclaw/

# Optimized (smaller binary)
go build -ldflags="-s -w" -o picoclaw ./cmd/picoclaw/

# With version info
go build -ldflags="-X main.version=1.0.0" -o picoclaw ./cmd/picoclaw/

# Install system-wide
go install ./cmd/picoclaw
```

## 🐛 Troubleshooting

```bash
# Check API key
echo $ANTHROPIC_API_KEY

# Check keychain (macOS)
security find-generic-password -s "Claude Code" -w

# Run with debug
picoclaw agent -d -m "Test"

# Fix build error (workspace not found)
cd cmd/picoclaw && cp -r ../../workspace . && cd ../..
```

## 📚 More Help

- Full guide: [COMMAND_LINE_GUIDE.md](COMMAND_LINE_GUIDE.md)
- Provider docs: [pkg/providers/README_UPDATES.md](pkg/providers/README_UPDATES.md)
- Examples: `pkg/providers/*_example.go`

## ✨ You're Ready!

```bash
picoclaw agent -m "Hello, Claude!"
```

🦞 **Happy coding!**
