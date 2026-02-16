# Telegram Channel Setup Guide

**Complete guide to integrating PicoClaw with Telegram Bot API**

---

## 📋 Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Configuration Methods](#configuration-methods)
- [Advanced Features](#advanced-features)
- [Troubleshooting](#troubleshooting)
- [Security Best Practices](#security-best-practices)
- [API Reference](#api-reference)

---

## 🎯 Overview

The Telegram channel integration allows PicoClaw to:

- ✅ **Receive messages** from Telegram users via Bot API
- ✅ **Send responses** with rich formatting (HTML, bold, italic, code blocks)
- ✅ **Process media** - photos, voice messages, audio, documents
- ✅ **Transcribe voice** messages using Groq Whisper API
- ✅ **Access control** with user/username allowlists
- ✅ **Proxy support** for regions with restricted access
- ✅ **Typing indicators** showing bot is processing
- ✅ **Message editing** for smoother conversation flow

**Architecture**: Uses long polling (no webhook required, easier for local development)

---

## 📦 Prerequisites

### 1. Create a Telegram Bot

1. **Open Telegram** and search for [@BotFather](https://t.me/botfather)
2. **Send** `/newbot` command
3. **Follow prompts**:
   - Choose a display name (e.g., "My PicoClaw Assistant")
   - Choose a username (must end with "bot", e.g., "mypicoclaw_bot")
4. **Save the token** - You'll get a message like:
   ```
   Use this token to access the HTTP API:
   123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ
   ```

### 2. Get Your User ID

**Method A: Use IDBot**
1. Search for [@userinfobot](https://t.me/userinfobot) on Telegram
2. Start a chat - it will reply with your User ID (e.g., `987654321`)

**Method B: Use Raw Updates (Advanced)**
1. Send a message to your bot
2. Visit: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
3. Look for `"from":{"id":987654321}` in the JSON response

### 3. Optional: Voice Transcription Setup

To enable voice message transcription, you need a Groq API key:

1. Visit [console.groq.com](https://console.groq.com)
2. Sign up/login and create an API key
3. Add to `.env`: `GROQ_API_KEY=gsk_...`

---

## 🚀 Quick Start

### Step 1: Configure Environment Variables

Edit your `.env` file:

```bash
# Telegram Bot Configuration
PICOCLAW_CHANNELS_TELEGRAM_ENABLED=true
PICOCLAW_CHANNELS_TELEGRAM_TOKEN=123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=987654321  # Your user ID

# Optional: Proxy (for regions with Telegram restrictions)
# PICOCLAW_CHANNELS_TELEGRAM_PROXY=socks5://127.0.0.1:1080

# Optional: Voice transcription
# GROQ_API_KEY=gsk_...
```

### Step 2: Start PicoClaw

```bash
# Build the project
cd /Users/ghu/aiworker/picoclaw
go build -o picoclaw cmd/picoclaw/main.go

# Run with Telegram channel enabled
./picoclaw

# Or run directly with go
go run cmd/picoclaw/main.go
```

### Step 3: Test the Bot

1. Open Telegram and search for your bot (@mypicoclaw_bot)
2. Send `/start` or any message
3. You should see "Thinking... 💭" followed by a response

**Example interaction**:
```
You: Hello!
Bot: Thinking... 💭
Bot: Hello! I'm PicoClaw, your AI assistant. How can I help you today?
```

---

## ⚙️ Configuration Methods

PicoClaw supports three configuration methods (in priority order):

### Method 1: Environment Variables (Recommended for Production)

```bash
# .env file
PICOCLAW_CHANNELS_TELEGRAM_ENABLED=true
PICOCLAW_CHANNELS_TELEGRAM_TOKEN=123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ
PICOCLAW_CHANNELS_TELEGRAM_PROXY=socks5://127.0.0.1:1080
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=987654321,123456789
```

**Advantages**:
- ✅ Works with Docker/Kubernetes
- ✅ Keeps secrets out of config files
- ✅ Easy to change without rebuilding

### Method 2: JSON Configuration File (Recommended for Development)

Edit `~/.picoclaw/config.json` or your custom config file:

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ",
      "proxy": "socks5://127.0.0.1:1080",
      "allow_from": ["987654321", "123456789"]
    }
  }
}
```

**Load custom config**:
```bash
./picoclaw --config /path/to/config.json
```

### Method 3: Hybrid (Environment Variables Override JSON)

Best practice: Use JSON for structure, environment variables for secrets:

**config.json**:
```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "allow_from": ["987654321"]
    }
  }
}
```

**.env**:
```bash
# Token from environment (more secure)
PICOCLAW_CHANNELS_TELEGRAM_TOKEN=123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ
```

---

## 🎨 Advanced Features

### 1. Access Control (Allowlist)

Control who can use your bot:

#### Allow Specific User IDs

```bash
# Single user
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=987654321

# Multiple users (comma-separated)
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=987654321,123456789,555666777
```

Or in JSON:
```json
{
  "channels": {
    "telegram": {
      "allow_from": ["987654321", "123456789"]
    }
  }
}
```

#### Allow by Username

```bash
# Allow specific usernames
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=alice,bob,charlie
```

#### Allow by User ID + Username

The system internally uses format `{user_id}|{username}`, so you can mix:

```bash
# Mix user IDs and usernames
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=987654321,alice,123456789
```

#### Public Bot (Allow All Users)

**⚠️ WARNING: Only for trusted environments!**

```bash
# Empty allowlist = allow everyone
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=
```

Or in JSON:
```json
{
  "channels": {
    "telegram": {
      "allow_from": []
    }
  }
}
```

### 2. Proxy Configuration

For regions where Telegram is restricted (e.g., China, Iran):

#### SOCKS5 Proxy

```bash
PICOCLAW_CHANNELS_TELEGRAM_PROXY=socks5://127.0.0.1:1080
```

#### HTTP Proxy

```bash
PICOCLAW_CHANNELS_TELEGRAM_PROXY=http://proxy.example.com:8080
```

#### Authenticated Proxy

```bash
PICOCLAW_CHANNELS_TELEGRAM_PROXY=socks5://username:password@127.0.0.1:1080
```

**Popular proxy tools**:
- **Clash**: `socks5://127.0.0.1:7891`
- **V2Ray**: `socks5://127.0.0.1:1080`
- **Shadowsocks**: `socks5://127.0.0.1:1080`

### 3. Voice Message Transcription

Enable automatic voice-to-text transcription:

#### Setup

```bash
# 1. Get Groq API key from console.groq.com
# 2. Add to .env
GROQ_API_KEY=gsk_...

# 3. Restart PicoClaw
```

#### How it Works

When a user sends a voice message:
1. **Download** - Voice file downloaded as `.ogg`
2. **Transcribe** - Sent to Groq Whisper API
3. **Process** - Transcribed text sent to agent
4. **Respond** - Bot replies based on transcription

**User sees**:
```
User: [sends voice message: "What's the weather today?"]
Bot: Thinking... 💭
Bot: [voice transcription: What's the weather today?]

     Based on your question, I'd need to know your location...
```

#### Timeout

Transcription has a 30-second timeout:
- If successful: `[voice transcription: {text}]`
- If failed: `[voice (transcription failed)]`
- If no Groq key: `[voice]`

### 4. Rich Message Formatting

PicoClaw automatically converts Markdown to Telegram HTML:

#### Supported Formatting

| Markdown | Telegram Display |
|----------|------------------|
| `**bold**` | **bold** |
| `__bold__` | **bold** |
| `_italic_` | *italic* |
| `~~strikethrough~~` | ~~strikethrough~~ |
| `` `code` `` | `code` |
| ` ```code block``` ` | Code block |
| `[link](url)` | Clickable link |
| `# Header` | Header (plain text) |
| `> Quote` | Quote (plain text) |

**Example bot response**:
```markdown
Here's how to **install** the package:

```bash
npm install picoclaw
```

Visit [documentation](https://example.com) for more info.
```

Rendered in Telegram as formatted HTML with clickable links and code blocks.

### 5. Media Processing

#### Photos

Users can send photos with optional captions:

```
User: [photo] What's in this image?
Bot: Analyzing the photo...
```

- **Format**: `.jpg` (highest resolution photo selected)
- **Size limit**: 20 MB (Telegram API limit)

#### Documents

```
User: [document.pdf] Summarize this document
Bot: Processing...
```

Supported formats: PDF, TXT, DOCX, etc.

#### Audio Files

```
User: [audio.mp3]
Bot: [audio]
```

Audio files are downloaded but not transcribed (unlike voice messages).

### 6. Typing Indicators

Bot sends "typing..." action while processing:

```go
// Automatically handled by the channel
c.bot.SendChatAction(ctx, telego.ChatActionTyping)
```

**User experience**:
1. User sends message
2. Bot shows "typing..." in chat
3. "Thinking... 💭" placeholder message appears
4. Final response replaces the placeholder

### 7. Group Chat Support

The bot works in both:
- **Private chats** (one-on-one)
- **Group chats** (if added to a group)

**Metadata includes**:
```json
{
  "is_group": "true",
  "message_id": "123",
  "user_id": "987654321",
  "username": "alice",
  "first_name": "Alice"
}
```

**⚠️ Note**: For groups, consider privacy settings on @BotFather:
- `/setprivacy` - Control when bot sees messages
- Disabled: Bot sees all messages
- Enabled: Bot only sees commands and @mentions

---

## 🛠️ Troubleshooting

### Issue 1: Bot Not Responding

**Symptoms**: Send messages, no response

**Check**:

1. **Bot is running**:
   ```bash
   # Look for this log
   [telegram] Telegram bot connected | username=mypicoclaw_bot
   ```

2. **User ID in allowlist**:
   ```bash
   # Check logs for
   [telegram] Message rejected by allowlist | user_id=987654321
   ```

   **Fix**: Add your user ID to `allow_from`

3. **Token is valid**:
   ```bash
   # Test with curl
   curl https://api.telegram.org/bot123456789:ABC.../getMe
   ```

   Should return bot info, not `{"ok":false,"error_code":401}`

### Issue 2: "Failed to create telegram bot"

**Error**:
```
Error: failed to create telegram bot: 401 Unauthorized
```

**Cause**: Invalid bot token

**Fix**:
1. Verify token in `.env` matches @BotFather token exactly
2. Check for extra spaces or quotes
3. Regenerate token via @BotFather if needed:
   - Send `/token` to @BotFather
   - Select your bot
   - Get new token

### Issue 3: Proxy Connection Failed

**Error**:
```
Error: failed to start long polling: dial tcp: lookup socks5: no such host
```

**Fix**:

1. **Check proxy format**:
   ```bash
   # Correct formats
   socks5://127.0.0.1:1080
   http://127.0.0.1:8080

   # Wrong (missing protocol)
   127.0.0.1:1080  # ❌
   ```

2. **Verify proxy is running**:
   ```bash
   # Test SOCKS5 proxy
   curl --proxy socks5://127.0.0.1:1080 https://api.telegram.org
   ```

3. **Try without proxy first**:
   ```bash
   PICOCLAW_CHANNELS_TELEGRAM_PROXY=
   ```

### Issue 4: Voice Transcription Not Working

**Symptoms**: Voice messages show as `[voice]` instead of transcription

**Check**:

1. **Groq API key set**:
   ```bash
   echo $GROQ_API_KEY
   # Should output: gsk_...
   ```

2. **Transcriber initialized**:
   ```bash
   # Look for log
   [telegram] Voice transcribed successfully | text=...

   # Or error
   [telegram] Voice transcription failed | error=...
   ```

3. **Groq quota**:
   - Free tier: 14,400 requests/day
   - Check at [console.groq.com](https://console.groq.com)

### Issue 5: HTML Parsing Failed

**Error in logs**:
```
[telegram] HTML parse failed, falling back to plain text
```

**Cause**: Invalid HTML characters in response (usually `<`, `>`, `&`)

**Automatic fix**: Bot automatically retries with plain text mode

**Manual fix**: Messages are auto-escaped, but if you're seeing this frequently, it indicates:
- Complex markdown in responses
- URLs with special characters
- Emojis in code blocks

### Issue 6: Updates Channel Closed

**Log message**:
```
[telegram] Updates channel closed, reconnecting...
```

**Cause**: Network interruption or long polling timeout

**Behavior**:
- Automatically attempts to reconnect
- No action needed (self-healing)

**If persists**:
1. Check network stability
2. Check Telegram API status: [telegram.org/status](https://telegram.org/status)
3. Try with/without proxy

---

## 🔒 Security Best Practices

### 1. Never Commit Tokens

**❌ Bad**:
```json
{
  "channels": {
    "telegram": {
      "token": "123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ"  // ❌ In version control
    }
  }
}
```

**✅ Good**:
```bash
# .env (add to .gitignore)
PICOCLAW_CHANNELS_TELEGRAM_TOKEN=123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ
```

### 2. Use Allowlists in Production

**❌ Bad** (Public bot):
```bash
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=
```

**✅ Good** (Restricted access):
```bash
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=987654321,123456789
```

### 3. Rotate Tokens Regularly

Via @BotFather:
1. Send `/token`
2. Select your bot
3. Choose "Revoke current token"
4. Update `.env` with new token

### 4. Monitor Logs

Enable debug logging to detect suspicious activity:

```bash
# Set log level to debug
PICOCLAW_LOG_LEVEL=debug
```

Watch for:
- Rejected messages from unknown users
- Unusual message patterns
- Failed authentication attempts

### 5. Rate Limiting

The Telegram API has built-in rate limits:
- **Private chats**: 1 msg/second per chat
- **Groups**: 20 msg/minute per group
- **Global**: 30 msg/second per bot

PicoClaw respects these automatically, but be aware if processing many users.

### 6. Webhook vs Long Polling

Current implementation uses **long polling** (pulls updates):

**Advantages**:
- ✅ No public IP/domain required
- ✅ Works behind NAT/firewall
- ✅ Perfect for development

**For production at scale**, consider webhooks:
- Telegram pushes updates to your server
- More efficient for high-volume bots
- Requires HTTPS endpoint

---

## 📚 API Reference

### Configuration Structure

```go
type TelegramConfig struct {
    Enabled   bool                // Enable/disable channel
    Token     string              // Bot token from @BotFather
    Proxy     string              // Optional proxy URL
    AllowFrom FlexibleStringSlice // User IDs or usernames
}
```

### Environment Variables

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `PICOCLAW_CHANNELS_TELEGRAM_ENABLED` | bool | No | `false` | Enable Telegram channel |
| `PICOCLAW_CHANNELS_TELEGRAM_TOKEN` | string | Yes | - | Bot token from @BotFather |
| `PICOCLAW_CHANNELS_TELEGRAM_PROXY` | string | No | - | Proxy URL (socks5:// or http://) |
| `PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM` | string | No | `[]` | Comma-separated user IDs/usernames |

### Message Metadata

Every incoming message includes metadata:

```go
metadata := map[string]string{
    "message_id": "123",        // Telegram message ID
    "user_id":    "987654321",  // Sender's user ID
    "username":   "alice",      // Sender's username (if set)
    "first_name": "Alice",      // Sender's first name
    "is_group":   "false",      // true if from group chat
}
```

### Supported Media Types

| Type | Extension | Transcription | Size Limit |
|------|-----------|---------------|------------|
| Photo | `.jpg` | No | 20 MB |
| Voice | `.ogg` | Yes (with Groq) | 20 MB |
| Audio | `.mp3` | No | 50 MB |
| Document | Various | No | 20 MB |

### Library Dependencies

```go
import (
    "github.com/mymmrac/telego"           // Telegram Bot API
    tu "github.com/mymmrac/telego/telegoutil" // Utility functions
)
```

---

## 🎓 Complete Examples

### Example 1: Personal Bot (Single User)

**.env**:
```bash
# LLM Provider
ANTHROPIC_API_KEY=sk-ant-xxx

# Telegram
PICOCLAW_CHANNELS_TELEGRAM_ENABLED=true
PICOCLAW_CHANNELS_TELEGRAM_TOKEN=123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=987654321
```

**Usage**:
```bash
go run cmd/picoclaw/main.go
```

### Example 2: Team Bot (Multiple Users)

**config.json**:
```json
{
  "agents": {
    "defaults": {
      "model": "claude-3-5-sonnet-20241022",
      "max_tokens": 4096
    }
  },
  "channels": {
    "telegram": {
      "enabled": true,
      "allow_from": [
        "987654321",  // Alice
        "123456789",  // Bob
        "555666777"   // Charlie
      ]
    }
  }
}
```

**.env**:
```bash
ANTHROPIC_API_KEY=sk-ant-xxx
PICOCLAW_CHANNELS_TELEGRAM_TOKEN=123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ
```

### Example 3: Bot with Voice + Proxy (China)

**.env**:
```bash
# LLM
ANTHROPIC_API_KEY=sk-ant-xxx

# Telegram (with proxy)
PICOCLAW_CHANNELS_TELEGRAM_ENABLED=true
PICOCLAW_CHANNELS_TELEGRAM_TOKEN=123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ
PICOCLAW_CHANNELS_TELEGRAM_PROXY=socks5://127.0.0.1:1080
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=987654321

# Voice transcription
GROQ_API_KEY=gsk_...
```

### Example 4: Docker Deployment

**docker-compose.yml**:
```yaml
version: '3.8'
services:
  picoclaw:
    image: picoclaw:latest
    environment:
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}
      - PICOCLAW_CHANNELS_TELEGRAM_ENABLED=true
      - PICOCLAW_CHANNELS_TELEGRAM_TOKEN=${TELEGRAM_BOT_TOKEN}
      - PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=${TELEGRAM_ALLOWED_USERS}
    restart: unless-stopped
```

**.env**:
```bash
ANTHROPIC_API_KEY=sk-ant-xxx
TELEGRAM_BOT_TOKEN=123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ
TELEGRAM_ALLOWED_USERS=987654321,123456789
```

**Run**:
```bash
docker-compose up -d
docker-compose logs -f picoclaw
```

---

## 🔗 Related Documentation

- **[System Architecture](../SYSTEM_ARCHITECTURE.md)** - Complete system architecture
- **[Architecture Diagrams](ARCHITECTURE_DIAGRAMS.md)** - Visual architecture diagrams
- **[Azure OpenAI Setup](../AZURE_END_TO_END_SUCCESS.md)** - Azure provider configuration
- **[Channel Implementation](../pkg/channels/telegram.go)** - Source code
- **[Telegram Bot API](https://core.telegram.org/bots/api)** - Official API documentation

---

## 📞 Support

### Getting Help

1. **Check logs**: Look for error messages and warnings
2. **Verify configuration**: Double-check all settings
3. **Test components**: Test bot token, proxy, API keys separately
4. **Review examples**: See working configurations above

### Common Log Messages

| Message | Meaning | Action |
|---------|---------|--------|
| `Starting Telegram bot (polling mode)...` | Bot starting | Normal |
| `Telegram bot connected` | Successfully connected | Normal |
| `Message rejected by allowlist` | User not allowed | Add to allowlist |
| `HTML parse failed, falling back` | Formatting issue | Automatic fallback |
| `Voice transcribed successfully` | Voice transcription worked | Normal |
| `Failed to get file` | File download failed | Check network |

---

**Last updated**: February 15, 2026
**Implementation**: [pkg/channels/telegram.go](../pkg/channels/telegram.go)
**Status**: ✅ Production Ready

🦞 **PicoClaw - Your AI Assistant on Telegram!**
