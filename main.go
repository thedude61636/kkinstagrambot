package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// githubLink is the repository URL
const githubLink = "https://github.com/thedude61636/kkinstagrambot"

// welcomeMessage is sent when /start is used
const welcomeMessage = `Hello! I am the KKInstagram Bot. 🤖

I automatically convert instagram.com links into kksave.com links so they embed nicely in Telegram! 

Here is how to use me:
1. Add me to a group, and I'll automatically reply to any Instagram links with the converted link.
2. Send an Instagram link directly to me here, and I'll convert it.
3. Type @kkinstagrambot (or my actual username) in ANY chat followed by a link, and I'll let you send the converted link instantly! 

Source Code: ` + githubLink

var instaRegex = regexp.MustCompile(`https?://(www\.)?instagram\.com[^\s]*`)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN environment variable is required")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot is running...")
	b.Start(ctx)
}

func convertLink(foundLink string) string {
	newLink := strings.Replace(foundLink, "instagram.com", "kksave.com", 1)
	if idx := strings.Index(newLink, "?"); idx != -1 {
		newLink = newLink[:idx]
	}
	return newLink
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Handle inline queries
	if update.InlineQuery != nil {
		handleInlineQuery(ctx, b, update.InlineQuery)
		return
	}

	// Ignore non-Message updates
	if update.Message == nil {
		return
	}

	// Ignore non-text messages
	if update.Message.Text == "" {
		return
	}

	// Handle /start command
	if strings.HasPrefix(update.Message.Text, "/start") {
		// Ignore if it's not a private message or direct mentioned command, but simple prefix check is usually fine here
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   welcomeMessage,
			LinkPreviewOptions: &models.LinkPreviewOptions{
				IsDisabled: bot.True(),
			},
		})
		if err != nil {
			log.Printf("Failed to send start message: %v", err)
		}
		return
	}

	// Check for Instagram links
	foundLink := instaRegex.FindString(update.Message.Text)
	if foundLink != "" {
		// Convert to kkinstagram and strip query params
		newLink := convertLink(foundLink)

		// Reply to user
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   newLink,
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})

		if err != nil {
			log.Printf("Failed to send reply: %v", err)
		}
	}
}

func handleInlineQuery(ctx context.Context, b *bot.Bot, query *models.InlineQuery) {
	// Check if the query contains an Instagram link
	foundLink := instaRegex.FindString(query.Query)

	if foundLink == "" {
		// Answer empty if no link is found
		_, _ = b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
			InlineQueryID: query.ID,
			Results:       []models.InlineQueryResult{},
		})
		return
	}

	// Convert the link and strip query params
	newLink := convertLink(foundLink)

	// Create a single result article
	result := &models.InlineQueryResultArticle{
		ID:    query.ID,
		Title: "Send KKSave Link",
		InputMessageContent: &models.InputTextMessageContent{
			MessageText: newLink,
		},
		Description: newLink,
	}

	_, err := b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
		InlineQueryID: query.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       []models.InlineQueryResult{result},
	})
	if err != nil {
		log.Printf("Failed to answer inline query: %v", err)
	}
}
