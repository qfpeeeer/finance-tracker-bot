package events

import (
	"context"
	"fmt"
	tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nyanyamaga/finance-tracker-bot/app/keyboards"
	"log"
)

type BotMessageHandler struct {
	TbAPI        TbAPI
	StateManager StateManager
}

func (h *BotMessageHandler) HandleMessages(ctx context.Context, update tbapi.Update) {
	var err error

	userID := update.Message.From.ID
	messageText := update.Message.Text

	switch messageText {
	case keyboards.ActionMessages[keyboards.ActionAddSpending]:
		err = h.StateManager.TriggerStateChange(ctx, userID, "ChooseAddSpending", "")
	case keyboards.ActionMessages[keyboards.ActionNewSpendingCategory]:
		err = h.StateManager.TriggerStateChange(ctx, userID, "ChooseAddCategory", "")
	default:
		currentState, stateErr := h.StateManager.GetCurrentState(ctx, userID)
		if stateErr != nil {
			err = fmt.Errorf("failed to get current state: %v", stateErr)
			break
		}

		nextStates := currentState.AvailableTransitions()
		if len(nextStates) == 0 {
			err = fmt.Errorf("no available transitions from current state")
			break
		} else if len(nextStates) > 1 {
			err = fmt.Errorf("more than one available transition from current state")
			break
		}

		err = h.StateManager.TriggerStateChange(ctx, userID, nextStates[0], messageText)
	}

	if err != nil {
		log.Printf("[warn] error triggering state change: %v", err)
	}
}
