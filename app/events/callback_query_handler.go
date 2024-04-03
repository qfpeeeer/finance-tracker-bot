package events

import (
	"context"
	"fmt"
	tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type BotCallbackQueryHandler struct {
	StateManager StateManager
}

func (h *BotCallbackQueryHandler) HandleCallbackQuery(ctx context.Context, update tbapi.Update) {
	var err error

	if update.CallbackQuery == nil {
		log.Println("Received an update that is not a callback query")
		return
	}

	userID := update.CallbackQuery.From.ID
	callbackData := update.CallbackQuery.Data

	log.Printf("[info] handling callback query: user %d, data %s", userID, callbackData)

	currentState, err := h.StateManager.GetCurrentState(ctx, userID)
	if err != nil {
		log.Printf("Error retrieving current state for user %d: %v", userID, err)
		return
	}

	nextStates := currentState.AvailableTransitions()
	if len(nextStates) == 0 {
		err = fmt.Errorf("no available transitions from current state")
	} else if len(nextStates) > 1 {
		err = fmt.Errorf("more than one available transition from current state")
	}

	err = h.StateManager.TriggerStateChange(ctx, userID, nextStates[0], callbackData)

	if err != nil {
		log.Printf("[warn] error triggering state change: %v", err)
	}
}
