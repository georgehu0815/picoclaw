# PicoClaw Build Scripts Guide

Quick reference for build and run scripts.

---

## 📜 Available Scripts

### 1. `build-and-run.sh` - Main Build & Run Script

**Default usage** (builds and runs test):
```bash
./build-and-run.sh
```

**Custom commands**:
```bash
# Agent mode with custom message
./build-and-run.sh agent -m "Your question here"

# Gateway mode (for Telegram, Discord, etc.)
./build-and-run.sh gateway

# Version info
./build-and-run.sh version

# Help
./build-and-run.sh --help
```

**What it does**:
1. ✅ Copies workspace directory for `go:embed`
2. ✅ Builds the binary (`go build`)
3. ✅ Runs with specified command (default: agent test)

---

### 2. `quick-test.sh` - Fast Test (No Rebuild)

**Usage**:
```bash
# Default test (What is 2+2?)
./quick-test.sh

# Custom message
./quick-test.sh "Explain quantum computing in one sentence"
```

**What it does**:
- ⚡ Runs existing binary WITHOUT rebuilding
- 🚀 Faster for quick tests
- ❌ Fails if binary doesn't exist

---

### 3. `run-gateway.sh` - Gateway Mode

**Usage**:
```bash
# Start gateway with all configured channels
./run-gateway.sh

# Gateway with config file
./run-gateway.sh --config custom-config.json
```

**What it does**:
- 🌐 Starts PicoClaw in gateway mode
- 📱 Enables Telegram, Discord, Slack, etc.
- 🔄 Auto-builds if binary doesn't exist

---

## 🚀 Quick Start Workflow

### First Time Setup

```bash
# 1. Build and test
./build-and-run.sh

# 2. Test with custom question
./build-and-run.sh agent -m "What is the capital of France?"

# 3. Start gateway for chat apps
./run-gateway.sh
```

### Daily Development

```bash
# After code changes
./build-and-run.sh

# Quick test without rebuild
./quick-test.sh "test message"

# Rebuild and test
./build-and-run.sh agent -m "new test"
```

---

## 📋 Common Tasks

### Build Only (No Run)

```bash
go build -o picoclaw cmd/picoclaw/main.go
```

### Run Without Build

```bash
./picoclaw agent -m "your message"
```

### Clean Build

```bash
# Remove old binary
rm picoclaw

# Remove embedded workspace
rm -rf cmd/picoclaw/workspace

# Rebuild
./build-and-run.sh
```

### Check Binary Size

The script shows binary size after build:
```
✓ Build successful
   Binary size: 26M
```

---

## 🎯 Examples

### Example 1: Simple Question

```bash
./build-and-run.sh agent -m "What is Rust?"
```

### Example 2: Code Generation

```bash
./build-and-run.sh agent -m "Write a Python function to calculate factorial"
```

### Example 3: Multiple Commands

```bash
# Build
./build-and-run.sh

# Test 1
./quick-test.sh "What is AI?"

# Test 2
./quick-test.sh "Explain blockchain"

# Gateway mode
./run-gateway.sh
```

### Example 4: Custom Model/Provider

Edit `~/.picoclaw/config.json` first, then:
```bash
./build-and-run.sh agent -m "test with new config"
```

---

## 🔧 Script Details

### `build-and-run.sh` Features

- ✅ Automatic workspace preparation
- ✅ Build error detection
- ✅ Binary size display
- ✅ Colored output for readability
- ✅ Flexible command passing
- ✅ Default test mode

### Color Coding

- 🔵 **Blue**: Headers and info
- 🟢 **Green**: Success messages
- 🟡 **Yellow**: Actions in progress
- 🔴 **Red**: Errors

---

## 🐛 Troubleshooting

### "Binary not found"

```bash
# Run full build
./build-and-run.sh
```

### "Permission denied"

```bash
# Make scripts executable
chmod +x *.sh
```

### Build fails

```bash
# Check Go installation
go version

# Verify you're in project directory
pwd  # Should be /Users/ghu/aiworker/picoclaw

# Check dependencies
go mod tidy
```

### "workspace not found"

```bash
# The script auto-copies it, but if needed:
cp -r workspace cmd/picoclaw/workspace
```

---

## 📁 Project Structure

```
picoclaw/
├── build-and-run.sh       # Main build & run script
├── quick-test.sh          # Fast test (no rebuild)
├── run-gateway.sh         # Gateway mode
├── picoclaw               # Binary (created after build)
├── cmd/picoclaw/
│   ├── main.go
│   └── workspace/         # Copied by script for go:embed
├── workspace/             # Original workspace files
└── ~/.picoclaw/
    └── config.json        # User configuration
```

---

## ⚙️ Configuration

Scripts use configuration from:

1. **`~/.picoclaw/config.json`** - Main configuration
2. **`.env`** - Environment variables (optional)

**Example config** (`~/.picoclaw/config.json`):
```json
{
  "agents": {
    "defaults": {
      "provider": "anthropic",
      "model": "claude-3-haiku-20240307",
      "max_tokens": 4096,
      "temperature": 0.7
    }
  },
  "providers": {
    "anthropic": {
      "auth_method": "token"
    }
  }
}
```

---

## 🎨 Customization

### Add Custom Script

Create `my-script.sh`:
```bash
#!/bin/bash
./build-and-run.sh agent -m "My custom workflow"
```

Make executable:
```bash
chmod +x my-script.sh
```

### Modify Default Test

Edit `build-and-run.sh`, line ~49:
```bash
./picoclaw agent -m "Your custom default test"
```

---

## 📚 Related Documentation

- **[TELEGRAM_SETUP_GUIDE.md](docs/TELEGRAM_SETUP_GUIDE.md)** - Telegram integration
- **[COMMAND_LINE_GUIDE.md](COMMAND_LINE_GUIDE.md)** - CLI reference
- **[AZURE_END_TO_END_SUCCESS.md](AZURE_END_TO_END_SUCCESS.md)** - Azure setup
- **[SYSTEM_ARCHITECTURE.md](SYSTEM_ARCHITECTURE.md)** - System design

---

## 🔗 Quick Links

| Task | Command |
|------|---------|
| **Build & Test** | `./build-and-run.sh` |
| **Quick Test** | `./quick-test.sh "message"` |
| **Gateway Mode** | `./run-gateway.sh` |
| **Custom Message** | `./build-and-run.sh agent -m "message"` |
| **Help** | `./picoclaw --help` |

---

**Created**: February 15, 2026
**Status**: ✅ Production Ready
🦞 **PicoClaw - Ultra-Efficient AI Assistant**
