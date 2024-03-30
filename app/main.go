package main

import (
	"context"
	"fmt"
	tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/nyanyamaga/finance-tracker-bot/app/events"
	"github.com/nyanyamaga/finance-tracker-bot/app/keyboards"
	"github.com/nyanyamaga/finance-tracker-bot/app/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var revision = "local"

func main() {
	fmt.Printf("finance-tracker-bot %s\n", revision)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Printf("[warn] interrupt signal")
		cancel()
	}()

	if err := execute(ctx); err != nil {
		log.Printf("[error] %v", err)
		os.Exit(1)
	}
}

func execute(ctx context.Context) error {
	dataFilePath := os.Getenv("DATA_FILE_PATH")
	telegramToken := os.Getenv("TELEGRAM_TOKEN")

	dataDB, err := storage.NewSqliteDB(dataFilePath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite database: %v", err)
	}
	defer func(dataDB *sqlx.DB) {
		err = dataDB.Close()
		if err != nil {
			log.Printf("[warn] error closing sqlite database: %v", err)
		}
	}(dataDB)

	categoryDB, err := storage.NewCategory(dataDB)
	if err != nil {
		return fmt.Errorf("failed to initialize category storage: %v", err)
	}

	userStateDB, err := storage.NewUserState(dataDB)
	if err != nil {
		return fmt.Errorf("failed to initialize user state storage: %v", err)
	}

	tbAPI, err := tbapi.NewBotAPI(telegramToken)
	if err != nil {
		return fmt.Errorf("can't make telegram bot, %w", err)
	}
	tbAPI.Debug = false

	botKeyboardProvider := keyboards.NewTbKeyboardProvider(categoryDB)
	botStateManager := events.NewBotStateManager(tbAPI, botKeyboardProvider, userStateDB, categoryDB)

	commandHandler := &events.BotCommandHandler{
		TbAPI:        tbAPI,
		TbKeyboards:  botKeyboardProvider,
		StateManager: botStateManager,
	}

	messageHandler := &events.BotMessageHandler{
		TbAPI:        tbAPI,
		StateManager: botStateManager,
	}

	callbackQueryHandler := &events.BotCallbackQueryHandler{}

	listener := events.TelegramListener{
		TbAPI:                tbAPI,
		CommandHandler:       commandHandler,
		MessageHandler:       messageHandler,
		CallbackQueryHandler: callbackQueryHandler,
	}

	err = listener.StartListening(ctx)
	if err != nil {
		return fmt.Errorf("failed to start listening: %w", err)
	}

	return nil
}
