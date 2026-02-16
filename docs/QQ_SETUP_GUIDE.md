# QQ Bot Setup Guide for PicoClaw

Complete guide to integrating PicoClaw with QQ (腾讯QQ) using the official QQ Bot API.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Detailed Setup](#detailed-setup)
- [Configuration](#configuration)
- [Features](#features)
- [Troubleshooting](#troubleshooting)
- [Security](#security)
- [FAQ](#faq)

---

## Overview

The QQ channel integration allows PicoClaw to communicate via QQ's official bot platform, supporting:

- **C2C Messages** (Private/Direct messages)
- **Group AT Messages** (Group messages where the bot is @mentioned)
- **WebSocket Connection** (Real-time message delivery)
- **Automatic Token Management** (Auto-refresh access tokens)

### Architecture

```
QQ Platform → WebSocket → PicoClaw Gateway → Agent → Response → QQ
```

The integration uses the official [tencent-connect/botgo](https://github.com/tencent-connect/botgo) SDK with OAuth2 token-based authentication.

---

## Prerequisites

1. **QQ Account**: A verified QQ account
2. **QQ Open Platform Access**: Developer account at [q.qq.com](https://q.qq.com)
3. **Bot Application**: Created bot with AppID and AppSecret
4. **PicoClaw Installed**: Version with QQ channel support

---

## Quick Start

### 1. Create Your QQ Bot

**Step 1: Register as a Developer**

1. Visit [QQ Open Platform](https://q.qq.com/#/)
2. Log in with your QQ account
3. Complete developer verification (requires identity verification for Chinese users)

**Step 2: Create a Bot Application**

1. Navigate to **Bot Management** (机器人管理)
2. Click **Create Bot** (创建机器人)
3. Fill in bot information:
   - **Bot Name**: Your bot's display name
   - **Bot Avatar**: Upload a profile picture
   - **Description**: Brief description of your bot
4. Submit for review (approval typically takes 1-3 business days)

**Step 3: Get Credentials**

Once approved:
1. Go to your bot's **Development Settings** (开发设置)
2. Find and copy:
   - **AppID** (also called BotAppID)
   - **AppSecret** (also called Token)

### 2. Configure PicoClaw

Edit `~/.picoclaw/config.json`:

```json
{
  "channels": {
    "qq": {
      "enabled": true,
      "app_id": "YOUR_BOT_APP_ID",
      "app_secret": "YOUR_BOT_APP_SECRET",
      "allow_from": []
    }
  }
}
```

### 3. Start the Gateway

```bash
picoclaw gateway
```

You should see:
```
[INFO] qq: Starting QQ bot (WebSocket mode)
[INFO] qq: Got WebSocket info {shards=1}
[INFO] qq: QQ bot started successfully
```

### 4. Test Your Bot

**Private Message Test:**
1. Find your bot in QQ (search by bot name or ID)
2. Send a message: "Hello"
3. The bot should respond

**Group Message Test:**
1. Add the bot to a QQ group
2. Send a message: "@YourBotName hello"
3. The bot should respond

---

## Detailed Setup

### Understanding QQ Bot Permissions

QQ bots have different permission scopes:

| Permission | Description | Required For |
|------------|-------------|--------------|
| **C2C Messages** | Direct/private messages | Private chats |
| **Group AT Messages** | Group messages with @mention | Group chats |
| **Group Messages** | All group messages (requires special approval) | Advanced use cases |

> **Note**: PicoClaw currently supports C2C and Group AT messages. For full group message access (without @mention), you need to apply for additional permissions through the QQ Open Platform.

### Bot Configuration Options

#### Basic Configuration

```json
{
  "channels": {
    "qq": {
      "enabled": true,
      "app_id": "102400XXXX",
      "app_secret": "xxxxxxxxxxxxxxxxxxxxx",
      "allow_from": []
    }
  }
}
```

#### Access Control (allow_from)

Control who can interact with your bot:

**Allow everyone** (default):
```json
"allow_from": []
```

**Restrict to specific users**:
```json
"allow_from": ["123456789", "987654321"]
```

To find a user's QQ ID:
- It's their QQ number
- You can view it in their QQ profile

### Environment Variables

You can also configure QQ channel via environment variables:

```bash
export PICOCLAW_CHANNELS_QQ_ENABLED=true
export PICOCLAW_CHANNELS_QQ_APP_ID="your_app_id"
export PICOCLAW_CHANNELS_QQ_APP_SECRET="your_app_secret"
export PICOCLAW_CHANNELS_QQ_ALLOW_FROM="123456789,987654321"
```

Then start the gateway:
```bash
picoclaw gateway
```

---

## Configuration

### Complete Example

Here's a full configuration example with QQ and other channels:

```json
{
  "agents": {
    "defaults": {
      "workspace": "~/.picoclaw/workspace",
      "model": "glm-4.7",
      "max_tokens": 8192,
      "temperature": 0.7,
      "max_tool_iterations": 20
    }
  },
  "channels": {
    "qq": {
      "enabled": true,
      "app_id": "102400XXXX",
      "app_secret": "xxxxxxxxxxxxxxxxxxxxx",
      "allow_from": []
    },
    "telegram": {
      "enabled": false,
      "token": "",
      "allow_from": []
    },
    "discord": {
      "enabled": false,
      "token": "",
      "allow_from": []
    }
  },
  "providers": {
    "zhipu": {
      "api_key": "your_zhipu_api_key",
      "api_base": "https://open.bigmodel.cn/api/paas/v4"
    }
  }
}
```

### Multi-Channel Setup

You can run multiple channels simultaneously:

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
      "token": "YOUR_TELEGRAM_BOT_TOKEN",
      "allow_from": ["YOUR_TELEGRAM_USER_ID"]
    }
  }
}
```

---

## Features

### 1. Private Messages (C2C)

Send direct messages to your bot:

```
User: "What's the weather today?"
Bot: "I'll check the weather for you..."
```

The bot receives and responds to private messages in real-time via WebSocket.

### 2. Group Messages with @mention

In QQ groups, mention your bot to get responses:

```
User: "@PicoClaw what is 2+2?"
Bot: "2+2 equals 4"
```

### 3. Message Deduplication

The QQ channel automatically handles duplicate messages:
- Maintains a cache of processed message IDs
- Auto-cleanup when cache exceeds 10,000 entries
- Prevents duplicate processing during reconnections

### 4. Automatic Token Refresh

OAuth2 tokens are automatically refreshed:
- Background refresh process
- No manual token management needed
- Seamless reconnection on token expiry

### 5. WebSocket Reliability

- Automatic reconnection on connection loss
- Heartbeat monitoring
- Graceful shutdown handling

---

## Troubleshooting

### Bot Not Starting

**Error**: `QQ app_id and app_secret not configured`

**Solution**: Check your configuration file:
```bash
cat ~/.picoclaw/config.json | grep -A 5 '"qq"'
```

Ensure `app_id` and `app_secret` are correctly set.

---

### WebSocket Connection Failed

**Error**: `failed to get websocket info`

**Possible causes**:
1. **Invalid credentials**: Verify your AppID and AppSecret
2. **Bot not approved**: Check if your bot is approved on QQ Open Platform
3. **Network issues**: Check your internet connection
4. **Firewall**: Ensure WebSocket connections are allowed

**Solution**:
```bash
# Test your credentials manually
curl -X POST "https://bots.qq.com/app/getAppAccessToken" \
  -H "Content-Type: application/json" \
  -d '{"appId":"YOUR_APP_ID","clientSecret":"YOUR_APP_SECRET"}'
```

---

### Bot Not Responding to Messages

**Check 1: Is the bot running?**
```bash
# Check logs
tail -f ~/.picoclaw/logs/picoclaw.log
```

Look for:
```
[INFO] qq: Received C2C message {sender=123456789, length=5}
```

**Check 2: Access control**

If you set `allow_from`, make sure your QQ ID is in the list:
```json
"allow_from": ["YOUR_QQ_NUMBER"]
```

**Check 3: Group messages**

In groups, you MUST @mention the bot:
```
✅ Correct: "@BotName hello"
❌ Wrong:   "hello" (without @mention)
```

---

### Token Refresh Failures

**Error**: `failed to start token refresh`

**Solution**:
1. Verify AppSecret is correct (not the token, but the secret)
2. Check if your bot's permissions are still active
3. Ensure system time is correct (OAuth2 requires accurate time)

```bash
# Check system time
date
# Sync time if needed (Linux)
sudo ntpdate -s time.nist.gov
```

---

### High Memory Usage

If you notice increasing memory usage over time:

**Cause**: Message ID cache growing too large

**Current behavior**: Auto-cleanup at 10,000 entries (removes oldest 5,000)

**Manual cleanup**: Restart the gateway periodically:
```bash
# Graceful restart
pkill -SIGTERM picoclaw
picoclaw gateway
```

---

## Security

### Best Practices

1. **Protect Your Credentials**
   - Never commit `config.json` to version control
   - Use environment variables in production
   - Rotate AppSecret regularly

2. **Use Access Control**
   ```json
   "allow_from": ["YOUR_QQ_NUMBER"]
   ```
   Only allow trusted users to interact with your bot.

3. **Monitor Logs**
   ```bash
   tail -f ~/.picoclaw/logs/picoclaw.log
   ```
   Watch for suspicious activity.

4. **Rate Limiting**
   The QQ platform has rate limits. Excessive requests may result in temporary bans.

### Credentials Storage

**Option 1: Configuration file** (Development)
```json
{
  "channels": {
    "qq": {
      "app_id": "102400XXXX",
      "app_secret": "xxxxxxxxxxxxxxxxxxxxx"
    }
  }
}
```

**Option 2: Environment variables** (Production)
```bash
# Add to ~/.bashrc or ~/.zshrc
export PICOCLAW_CHANNELS_QQ_APP_ID="102400XXXX"
export PICOCLAW_CHANNELS_QQ_APP_SECRET="xxxxxxxxxxxxxxxxxxxxx"
```

**Option 3: Secret management** (Enterprise)
- Use HashiCorp Vault
- AWS Secrets Manager
- Azure Key Vault

---

## FAQ

### Q: Do I need a verified QQ account?

**A**: Yes, creating a QQ bot requires identity verification through the QQ Open Platform. This is a requirement from Tencent.

---

### Q: Can I use the same bot on multiple PicoClaw instances?

**A**: No. Each QQ bot can only maintain one active WebSocket connection. Running multiple instances will cause conflicts:
```
[ERROR] qq: WebSocket session error {error=conflict: existing connection}
```

**Solution**: Use different bots for different instances.

---

### Q: How do I get my QQ user ID?

**A**: Your QQ user ID is your QQ number. You can find it in:
- QQ profile page
- QQ settings
- Your QQ account details

---

### Q: Does the bot work in all QQ groups?

**A**: The bot works in any group where:
1. It has been added as a member
2. It has proper permissions (set via QQ Open Platform)
3. Users @mention it in messages

For full group message access (without @mention), you need special approval from QQ.

---

### Q: What's the difference between AppID and BotAppID?

**A**: They are the same thing. The QQ Open Platform sometimes refers to it as "AppID" and sometimes as "BotAppID" - use whichever is displayed in your bot's settings.

---

### Q: Can I customize the bot's responses?

**A**: Yes! Customize via PicoClaw's workspace files:
- `~/.picoclaw/workspace/IDENTITY.md` - Bot personality
- `~/.picoclaw/workspace/SOUL.md` - Bot behavior
- `~/.picoclaw/workspace/USER.md` - User preferences

---

### Q: Does PicoClaw support QQ rich media?

**A**: Currently, PicoClaw supports:
- ✅ Text messages
- ❌ Images (roadmap)
- ❌ Files (roadmap)
- ❌ Voice (roadmap)
- ❌ Stickers (roadmap)

Check the [roadmap](https://github.com/sipeed/picoclaw/issues) for updates.

---

### Q: How do I debug connection issues?

**A**: Enable debug logging:

1. Check logs:
   ```bash
   tail -f ~/.picoclaw/logs/picoclaw.log | grep qq
   ```

2. Look for these key messages:
   ```
   [INFO] qq: Starting QQ bot (WebSocket mode)
   [INFO] qq: Got WebSocket info {shards=1}
   [INFO] qq: QQ bot started successfully
   [INFO] qq: Received C2C message {sender=..., length=...}
   ```

3. Common error patterns:
   ```
   [ERROR] qq: failed to get websocket info → Check credentials
   [ERROR] qq: WebSocket session error → Connection issue
   [WARN] qq: Received message with no sender ID → API issue
   ```

---

### Q: What regions does QQ bot support?

**A**: QQ bots primarily serve Chinese users. The QQ Open Platform requires:
- Chinese phone number for verification (in most cases)
- Identity verification for developers
- Compliance with Chinese internet regulations

International developers may face additional verification requirements.

---

## Additional Resources

- [QQ Open Platform](https://q.qq.com/#/)
- [QQ Bot Documentation](https://bot.q.qq.com/wiki/)
- [BotGo SDK GitHub](https://github.com/tencent-connect/botgo)
- [PicoClaw Documentation](https://github.com/sipeed/picoclaw)

---

## Summary

✅ **Quick Checklist**:

- [ ] Register at QQ Open Platform
- [ ] Create and verify bot application
- [ ] Get AppID and AppSecret
- [ ] Configure `~/.picoclaw/config.json`
- [ ] Start gateway: `picoclaw gateway`
- [ ] Test private messages
- [ ] Test group @mentions
- [ ] Set up access control (optional)
- [ ] Monitor logs for errors

**Need help?** Open an issue on [GitHub](https://github.com/sipeed/picoclaw/issues) or join the community discussion.

---

**🦐 PicoClaw - Let's Go! / 皮皮虾，我们走！**
