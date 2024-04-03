package events

import (
	"context"
	"fmt"
	tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
	"github.com/nyanyamaga/finance-tracker-bot/app/storage"
	"log"
)

// TbAPI is an interface for telegram bot API, only subset of methods used
type TbAPI interface {
	GetUpdatesChan(config tbapi.UpdateConfig) tbapi.UpdatesChannel
	Send(c tbapi.Chattable) (tbapi.Message, error)
	Request(c tbapi.Chattable) (*tbapi.APIResponse, error)
	GetChat(config tbapi.ChatInfoConfig) (tbapi.Chat, error)
}

type TbKeyboards interface {
	GetMainKeyboard() tbapi.ReplyKeyboardMarkup
	GetCategoryKeyboard(userID int64) tbapi.InlineKeyboardMarkup
}

type UserStateRepository interface {
	Write(entry storage.UserStateInfo) error
	Read(userID int64) (*storage.UserStateInfo, error)
}

type CategoriesRepository interface {
	AddOrUpdateCategory(info storage.CategoryInfo) error
	ListCategories(userID int64) ([]storage.CategoryInfo, error)
}

type SpendingsRepository interface {
	AddSpending(info storage.SpendingInfo) error
	ListSpendings(userID int64) ([]storage.SpendingInfo, error)
}

type CommandHandler interface {
	HandleCommands(ctx context.Context, update tbapi.Update)
}

type MessageHandler interface {
	HandleMessages(ctx context.Context, update tbapi.Update)
}

type CallbackQueryHandler interface {
	HandleCallbackQuery(ctx context.Context, update tbapi.Update)
}

type StateManager interface {
	InitializeUserFSM(ctx context.Context, userID int64)
	SetIdleState(ctx context.Context, userID int64)
	TriggerStateChange(ctx context.Context, userID int64, action, value string) error
	GetCurrentState(ctx context.Context, userID int64) (*fsm.FSM, error)
}

// send a message to the telegram as markdown first and if failed - as plain text
func send(tbMsg tbapi.Chattable, tbAPI TbAPI) error {
	withParseMode := func(tbMsg tbapi.Chattable, parseMode string) tbapi.Chattable {
		switch msg := tbMsg.(type) {
		case tbapi.MessageConfig:
			msg.ParseMode = parseMode
			msg.DisableWebPagePreview = true
			return msg
		case tbapi.EditMessageTextConfig:
			msg.ParseMode = parseMode
			msg.DisableWebPagePreview = true
			return msg
		case tbapi.EditMessageReplyMarkupConfig:
			return msg
		}
		return tbMsg // don't touch other types
	}

	msg := withParseMode(tbMsg, tbapi.ModeMarkdown) // try markdown first
	if _, err := tbAPI.Send(msg); err != nil {
		log.Printf("[warn] failed to send message as markdown, %v", err)
		msg = withParseMode(tbMsg, "") // try plain text
		if _, err := tbAPI.Send(msg); err != nil {
			return fmt.Errorf("can't send message to telegram: %w", err)
		}
	}
	return nil
}
