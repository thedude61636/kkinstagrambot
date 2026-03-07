# KKInstagram Bot

A lightweight, zero-configuration Telegram Bot written in Go that automatically converts `instagram.com` links into `kkinstagram.com` links.

This bot is designed to be extremely resource-efficient, securely compiled down to a minimal binary, and easily deployed via Docker to any VPS (like Dokploy).

## Features

- **Passive Conversion:** Add the bot to any group (with Privacy Mode disabled) and it will automatically monitor chat for `instagram.com` URLs, replying to the user with the corrected `kkinstagram.com` URL so the embed previews load properly.
- **Direct Messages:** DM the bot an Instagram link and it will convert it immediately.
- **Inline Queries:** Type `@yourbotname https://instagram.com/p/something` in *any* Telegram chat box, and a popup will appear allowing you to send the converted link instantly without adding the bot to the group.

## Local Testing Setup

Since the app is written in Go, you'll need the Go compiler installed on your machine.

1. Clone the repository.
2. Run `go mod download` to fetch the telegram API library.
3. Start the bot by providing your token as an environment variable:

```powershell
# Windows
$env:BOT_TOKEN="your_token_here"; go run main.go

# Linux/Mac
BOT_TOKEN="your_token_here" go run main.go
```

## Deployment (Dokploy)

This repository includes a multi-stage `Dockerfile`. When deployed, it produces a tiny, zero-dependency `scratch` image that usually consumes less than 15MB of RAM on your VPS.

### Deploying via Dokploy
1. Push this repository to GitHub.
2. Go to your Dokploy Dashboard and click **Create Application**.
3. Point Dokploy to your GitHub repository and select the **main** branch.
4. Set the **Build Type** to **Dockerfile**.
5. Add your Telegram Bot Token in the Dokploy **Environment Variables** tab:
   - **Key:** `BOT_TOKEN`
   - **Value:** `your-bot-token-from-botfather`
6. Click **Deploy**.

## Telegram BotFather Configuration

To maximize the bot's features, ensure you configure it correctly in Telegram using `@BotFather`:

1. **Get your Token:** Send `/newbot` to BotFather and follow the prompts.
2. **Enable Inline Queries:** Send `/setinline` to BotFather, select your bot, and provide a placeholder text like "Search Instagram URLs...".
3. **Disable Privacy Mode (mportant for Groups):** Send `/setprivacy` to BotFather, select your bot, and choose **Disable**. This allows the bot to passively read all messages in a group looking for Instagram links, instead of only reading messages that start with a `/`.

## Technologies

- Go 1.24+
- `github.com/go-telegram/bot`
- Multi-stage Docker Builds

---
*Built with ❤️ in 10 minutes using Gemini and Antigravity.*
