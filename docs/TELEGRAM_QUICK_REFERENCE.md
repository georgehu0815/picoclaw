# Telegram Channel - Quick Reference

**One-page cheat sheet for PicoClaw Telegram integration**

---

## ⚡ 30-Second Setup

```bash
# 1. Get bot token from @BotFather
# 2. Get your user ID from @userinfobot
# 3. Add to .env:

PICOCLAW_CHANNELS_TELEGRAM_ENABLED=true
PICOCLAW_CHANNELS_TELEGRAM_TOKEN=123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=987654321

# 4. Run
go run cmd/picoclaw/main.go
```

---

## 🔧 Environment Variables

| Variable | Example | Required |
|----------|---------|----------|
| `PICOCLAW_CHANNELS_TELEGRAM_ENABLED` | `true` | ✅ |
| `PICOCLAW_CHANNELS_TELEGRAM_TOKEN` | `123456789:ABC...` | ✅ |
| `PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM` | `987654321,123456789` | Recommended |
| `PICOCLAW_CHANNELS_TELEGRAM_PROXY` | `socks5://127.0.0.1:1080` | Optional |

---

## 📝 JSON Configuration

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ",
      "proxy": "socks5://127.0.0.1:1080",
      "allow_from": ["987654321", "alice"]
    }
  }
}
```

---

## 🎯 Common Tasks

### Get Bot Token
1. Message [@BotFather](https://t.me/botfather)
2. Send `/newbot`
3. Follow prompts
4. Copy token: `123456789:ABC...`

### Get User ID
- Method 1: [@userinfobot](https://t.me/userinfobot)
- Method 2: `https://api.telegram.org/bot<TOKEN>/getUpdates`

### Allow Multiple Users
```bash
# Comma-separated
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=987654321,123456789,alice,bob
```

### Use Proxy (China/Iran)
```bash
# SOCKS5
PICOCLAW_CHANNELS_TELEGRAM_PROXY=socks5://127.0.0.1:1080

# HTTP
PICOCLAW_CHANNELS_TELEGRAM_PROXY=http://127.0.0.1:8080
```

### Enable Voice Transcription
```bash
GROQ_API_KEY=gsk_...
```

---

## ✨ Features

| Feature | Status | Notes |
|---------|--------|-------|
| Text messages | ✅ | Full markdown support |
| Photos | ✅ | Highest resolution |
| Voice messages | ✅ | Auto-transcribe with Groq |
| Audio files | ✅ | Downloaded, not transcribed |
| Documents | ✅ | All formats |
| Typing indicator | ✅ | Shows "typing..." |
| Message editing | ✅ | Smooth UX |
| Group chats | ✅ | Works in groups |
| Allowlist | ✅ | User ID or username |
| Proxy | ✅ | SOCKS5/HTTP |

---

## 🐛 Troubleshooting

| Problem | Solution |
|---------|----------|
| Bot not responding | Check `allow_from` includes your user ID |
| "401 Unauthorized" | Verify token is correct |
| Proxy connection failed | Check proxy format: `socks5://host:port` |
| Voice not transcribed | Add `GROQ_API_KEY` to .env |
| "HTML parse failed" | Automatic fallback to plain text (no action needed) |

---

## 📊 Supported Formats

### Markdown → Telegram HTML

| Markdown | Telegram |
|----------|----------|
| `**bold**` | **bold** |
| `_italic_` | *italic* |
| `` `code` `` | `code` |
| ` ```block``` ` | Code block |
| `[link](url)` | Clickable link |

---

## 🔒 Security Checklist

- [ ] Bot token in `.env` (not committed to git)
- [ ] `allow_from` configured (not empty)
- [ ] `.env` in `.gitignore`
- [ ] Debug logs disabled in production
- [ ] Regular token rotation via @BotFather

---

## 📦 Example Configurations

### Personal Bot
```bash
PICOCLAW_CHANNELS_TELEGRAM_ENABLED=true
PICOCLAW_CHANNELS_TELEGRAM_TOKEN=123456789:ABC...
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=987654321
```

### Team Bot
```bash
PICOCLAW_CHANNELS_TELEGRAM_ENABLED=true
PICOCLAW_CHANNELS_TELEGRAM_TOKEN=123456789:ABC...
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=987654321,123456789,555666777
```

### Bot with Proxy + Voice
```bash
PICOCLAW_CHANNELS_TELEGRAM_ENABLED=true
PICOCLAW_CHANNELS_TELEGRAM_TOKEN=123456789:ABC...
PICOCLAW_CHANNELS_TELEGRAM_PROXY=socks5://127.0.0.1:1080
PICOCLAW_CHANNELS_TELEGRAM_ALLOW_FROM=987654321
GROQ_API_KEY=gsk_...
```

---

## 📚 Full Documentation

See [TELEGRAM_SETUP_GUIDE.md](TELEGRAM_SETUP_GUIDE.md) for:
- Detailed setup instructions
- Advanced features
- Complete troubleshooting guide
- Security best practices
- API reference

---

**Quick Links**:
- [@BotFather](https://t.me/botfather) - Create bot
- [@userinfobot](https://t.me/userinfobot) - Get user ID
- [Telegram Bot API](https://core.telegram.org/bots/api) - Official docs
- [Groq Console](https://console.groq.com) - Voice transcription API

---

**Last updated**: February 15, 2026
**Status**: ✅ Production Ready
🦞 **PicoClaw**
