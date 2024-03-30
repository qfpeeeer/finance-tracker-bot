package events

import (
	"context"
	tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type BotCallbackQueryHandler struct{}

func (h *BotCallbackQueryHandler) HandleCallbackQuery(ctx context.Context, update tbapi.Update) {
	log.Println("Handle callback query")
}
