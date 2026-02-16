# QQ Bot Quick Reference

One-page cheat sheet for PicoClaw QQ integration.

## Setup (5 minutes)

### 1. Create Bot
```
1. Visit: https://q.qq.com
2. Create bot → Get AppID + AppSecret
3. Wait for approval (1-3 days)
```

### 2. Configure
```json
{
  "channels": {
    "qq": {
      "enabled": true,
      "app_id": "YOUR_APP_ID",
      "app_secret": "YOUR_APP_SECRET",
      "allow_from": []
    }
  }
}
```

### 3. Start
```bash
picoclaw gateway
```

---

## Configuration Options

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `enabled` | boolean | Enable QQ channel | `true` |
| `app_id` | string | Bot AppID from QQ Open Platform | `"102400XXXX"` |
| `app_secret` | string | Bot AppSecret (Token) | `"xxxxxxxxxxxx"` |
| `allow_from` | array | Allowed QQ user IDs (empty = all) | `["123456789"]` |

---

## Environment Variables

```bash
export PICOCLAW_CHANNELS_QQ_ENABLED=true
export PICOCLAW_CHANNELS_QQ_APP_ID="102400XXXX"
export PICOCLAW_CHANNELS_QQ_APP_SECRET="xxxxxxxxxxxx"
export PICOCLAW_CHANNELS_QQ_ALLOW_FROM="123456789,987654321"
```

---

## Usage

### Private Messages
```
User → Bot: "Hello"
Bot → User: "Hi! How can I help you?"
```

### Group Messages
```
User in group: "@BotName what is 2+2?"
Bot in group: "2+2 equals 4"
```

> **Note**: In groups, you MUST @mention the bot!

---

## Access Control

### Allow Everyone
```json
"allow_from": []
```

### Restrict to Specific Users
```json
"allow_from": ["123456789", "987654321"]
```

Find your QQ ID: It's your QQ number (visible in QQ profile)

---

## Troubleshooting

### Bot Not Starting
```bash
# Check config
cat ~/.picoclaw/config.json | grep -A 5 '"qq"'

# Check logs
tail -f ~/.picoclaw/logs/picoclaw.log | grep qq
```

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `app_id and app_secret not configured` | Missing credentials | Add AppID/AppSecret to config |
| `failed to get websocket info` | Invalid credentials | Verify AppID/AppSecret |
| Bot not responding | Not @mentioned in group | Use `@BotName message` |
| `conflict: existing connection` | Multiple instances | Stop other instances |

---

## Features

| Feature | Supported |
|---------|-----------|
| Private messages (C2C) | ✅ |
| Group @mentions | ✅ |
| Text messages | ✅ |
| Images | ❌ (roadmap) |
| Files | ❌ (roadmap) |
| Voice | ❌ (roadmap) |
| Auto token refresh | ✅ |
| Message deduplication | ✅ |
| WebSocket reconnection | ✅ |

---

## Security Checklist

- [ ] Never commit `config.json` to git
- [ ] Use `allow_from` to restrict access
- [ ] Use environment variables in production
- [ ] Monitor logs for suspicious activity
- [ ] Rotate AppSecret regularly

---

## Quick Commands

```bash
# Start gateway
picoclaw gateway

# Check status
picoclaw status

# View logs
tail -f ~/.picoclaw/logs/picoclaw.log

# Test configuration
cat ~/.picoclaw/config.json | jq '.channels.qq'
```

---

## Example: Multi-Channel Setup

Run QQ + Telegram simultaneously:

```json
{
  "channels": {
    "qq": {
      "enabled": true,
      "app_id": "YOUR_QQ_APP_ID",
      "app_secret": "YOUR_QQ_APP_SECRET",
      "allow_from": []
    },
    "telegram": {
      "enabled": true,
      "token": "YOUR_TELEGRAM_TOKEN",
      "allow_from": ["YOUR_TELEGRAM_ID"]
    }
  }
}
```

---

## Resources

- **Full Guide**: [QQ Setup Guide](QQ_SETUP_GUIDE.md)
- **QQ Platform**: https://q.qq.com
- **Bot Docs**: https://bot.q.qq.com/wiki/
- **PicoClaw**: https://github.com/sipeed/picoclaw

---

## Need Help?

1. Check [QQ Setup Guide](QQ_SETUP_GUIDE.md)
2. Review logs: `tail -f ~/.picoclaw/logs/picoclaw.log`
3. Open issue: https://github.com/sipeed/picoclaw/issues

---

**🦐 皮皮虾，我们走！**
