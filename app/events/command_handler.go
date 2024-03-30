package events

import (
	"context"
	tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type BotCommandHandler struct {
	TbAPI        TbAPI
	TbKeyboards  TbKeyboards
	StateManager StateManager // Add StateManager to the command handler
}

func (h *BotCommandHandler) HandleCommands(ctx context.Context, update tbapi.Update) {
	userID := update.Message.From.ID

	if update.Message.Command() == "start" {
		h.StateManager.SetIdleState(ctx, userID)

		msg := tbapi.NewMessage(update.Message.Chat.ID, "Welcome! Choose an option.")
		msg.ReplyMarkup = h.TbKeyboards.GetMainKeyboard()

		if _, err := h.TbAPI.Send(msg); err != nil {
			log.Printf("[warn] error sending welcome message: %v", err)
		}
	}
}
