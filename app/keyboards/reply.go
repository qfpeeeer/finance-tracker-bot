package keyboards

import tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Action identifiers
const (
	ActionAddSpending         = "ADD_SPENDING"
	ActionNewSpendingCategory = "NEW_SPENDING_CATEGORY"
)

// ActionMessages maps action identifiers to user-facing text.
var ActionMessages = map[string]string{
	ActionAddSpending:         "Add spending",
	ActionNewSpendingCategory: "New spending category",
}

// GetMainKeyboard generates the main keyboard with dynamic actions.
func (tbk *TbKeyboardProvider) GetMainKeyboard() tbapi.ReplyKeyboardMarkup {
	return tbapi.ReplyKeyboardMarkup{
		Keyboard: [][]tbapi.KeyboardButton{
			{{Text: ActionMessages[ActionAddSpending]}},
			{{Text: ActionMessages[ActionNewSpendingCategory]}},
		},
		ResizeKeyboard: true,
	}
}
