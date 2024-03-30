package keyboards

import (
	"fmt"
	tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (tbk *TbKeyboardProvider) GetCategoryKeyboard(userID int64) tbapi.InlineKeyboardMarkup {
	categories, err := tbk.Storage.ListCategories(userID)
	if err != nil {
		log.Printf("Error retrieving categories: %v", err)
		return tbapi.NewInlineKeyboardMarkup()
	}

	var rows [][]tbapi.InlineKeyboardButton
	for _, category := range categories {
		buttonText := category.Emoji + " " + category.Name
		callbackData := fmt.Sprintf("category_%d", category.ID)

		row := []tbapi.InlineKeyboardButton{tbapi.NewInlineKeyboardButtonData(buttonText, callbackData)}
		rows = append(rows, row)
	}

	keyboard := tbapi.NewInlineKeyboardMarkup(rows...)
	return keyboard
}
