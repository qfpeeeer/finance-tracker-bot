package events

import (
	"context"
	"encoding/json"
	"fmt"
	tbapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/looplab/fsm"
	"github.com/nyanyamaga/finance-tracker-bot/app/storage"
	"log"
)

type BotStateManager struct {
	TbAPI       TbAPI
	TbKeyboards TbKeyboards
	UserState   UserStateRepository
	Categories  CategoriesRepository
	UserFSMs    map[int64]*fsm.FSM
	UserValues  map[int64]string
}

func NewBotStateManager(tbAPI TbAPI, tbKeyboards TbKeyboards, usRepository UserStateRepository, cRepository CategoriesRepository) *BotStateManager {
	return &BotStateManager{
		TbAPI:       tbAPI,
		TbKeyboards: tbKeyboards,
		UserState:   usRepository,
		Categories:  cRepository,
		UserFSMs:    make(map[int64]*fsm.FSM),
		UserValues:  make(map[int64]string),
	}
}

func (sm *BotStateManager) InitializeUserFSM(ctx context.Context, userID int64) {
	sm.UserFSMs[userID] = fsm.NewFSM(
		"Idle",
		fsm.Events{
			{Name: "ChooseAddSpending", Src: []string{"Idle"}, Dst: "AwaitingCategorySelection"},
			{Name: "ChooseAddCategory", Src: []string{"Idle"}, Dst: "AwaitingNewCategoryName"},
			{Name: "CategorySelected", Src: []string{"AwaitingCategorySelection"}, Dst: "AwaitingAmountInput"},
			{Name: "AmountEntered", Src: []string{"AwaitingAmountInput"}, Dst: "Idle"},
			{Name: "NewCategoryNameEntered", Src: []string{"AwaitingNewCategoryName"}, Dst: "AwaitingNewCategoryEmoji"},
			{Name: "NewCategoryEmojiEntered", Src: []string{"AwaitingNewCategoryEmoji"}, Dst: "AwaitingSaveCategoryName"},
			{Name: "SaveNewCategory", Src: []string{"AwaitingSaveCategoryName"}, Dst: "Idle"},
		},
		fsm.Callbacks{
			"leave_state":                     func(ctx context.Context, e *fsm.Event) { sm.leaveState(e, userID) },
			"enter_Idle":                      func(ctx context.Context, e *fsm.Event) { sm.promptEnterIdle(userID) },
			"enter_AwaitingCategorySelection": func(ctx context.Context, e *fsm.Event) { sm.promptCategorySelection(userID) },
			"enter_AwaitingAmountInput":       func(ctx context.Context, e *fsm.Event) { sm.promptAmountInput(userID) },
			"enter_AwaitingNewCategoryName":   func(ctx context.Context, e *fsm.Event) { sm.promptNewCategoryName(userID) },
			"enter_AwaitingNewCategoryEmoji":  func(ctx context.Context, e *fsm.Event) { sm.promptNewCategoryEmoji(userID) },
			"enter_AwaitingSaveCategoryName":  func(ctx context.Context, e *fsm.Event) { sm.promptSaveNewCategory(ctx, userID) },
		},
	)

	initialState := storage.UserStateInfo{
		UserID:   userID,
		State:    "Idle",
		DataJSON: "{}",
	}
	if err := sm.UserState.Write(initialState); err != nil {
		log.Printf("[error] Failed to create initial state for user %d: %v", userID, err)
	}
}

func (sm *BotStateManager) TriggerStateChange(ctx context.Context, userID int64, action, value string) error {
	var err error

	if _, exists := sm.UserFSMs[userID]; !exists {
		sm.InitializeUserFSM(ctx, userID)
	}
	userFSM := sm.UserFSMs[userID]
	sm.UserValues[userID] = value

	if userFSM.Can(action) {
		err = userFSM.Event(ctx, action, value)
	} else {
		err = fmt.Errorf("can't trigger %s event", action)
	}

	if err != nil {
		return err
	}
	return nil
}

func (sm *BotStateManager) GetCurrentState(ctx context.Context, userID int64) (*fsm.FSM, error) {
	if _, exists := sm.UserFSMs[userID]; !exists {
		return nil, fmt.Errorf("user %d has no state machine", userID)
	}
	userFSM := sm.UserFSMs[userID]

	return userFSM, nil
}

func (sm *BotStateManager) SetIdleState(ctx context.Context, userID int64) {
	sm.InitializeUserFSM(ctx, userID)
}

func (sm *BotStateManager) leaveState(e *fsm.Event, userID int64) {
	var (
		updatedData map[string]interface{}
		err         error
	)

	if e.Dst != "Idle" {
		updatedData, err = sm.getStateData(userID)

		if updatedData == nil {
			updatedData = make(map[string]interface{})
		}

		updatedData[e.Event] = sm.UserValues[userID]
	}

	dataJSON, err := json.Marshal(updatedData)
	if err != nil {
		log.Printf("[error] Failed to marshal updated state data to JSON for user %d: %v", userID, err)
		return
	}

	stateInfo := storage.UserStateInfo{
		UserID:   userID,
		State:    e.Event,
		DataJSON: string(dataJSON),
	}

	if err := sm.UserState.Write(stateInfo); err != nil {
		log.Printf("[error] Failed to save updated user state for user %d: %v", userID, err)
		return
	}

	log.Printf("[info] User %d entered state %s, data: %s", userID, e.Dst, dataJSON)
	delete(sm.UserValues, userID)
}

func (sm *BotStateManager) promptEnterIdle(userID int64) {
	err := sm.sendBotResponse(userID, "Choose an option:", sm.TbKeyboards.GetMainKeyboard())
	if err != nil {
		log.Printf("[warn] error sending main message: %v", err)
	}
}

func (sm *BotStateManager) promptCategorySelection(userID int64) {
	text := "Please select a category:"
	keyboard := sm.TbKeyboards.GetCategoryKeyboard(userID)

	err := sm.sendBotResponse(userID, text, &keyboard)
	if err != nil {
		log.Printf("[warn] error sending category selection prompt: %v", err)
		return
	}
}

func (sm *BotStateManager) promptAmountInput(userID int64) {
	text := "Please enter the amount:"

	err := sm.sendBotResponse(userID, text, nil)
	if err != nil {
		log.Printf("[warn] error sending amount prompt: %v", err)
		return
	}
}

func (sm *BotStateManager) promptNewCategoryName(userID int64) {
	text := "Please enter the name of the new category:"

	err := sm.sendBotResponse(userID, text, nil)
	if err != nil {
		log.Printf("[warn] error sending new category name prompt: %v", err)
		return
	}
}

func (sm *BotStateManager) promptNewCategoryEmoji(userID int64) {
	text := "Please enter the emoji for the new category:"

	err := sm.sendBotResponse(userID, text, nil)
	if err != nil {
		log.Printf("[warn] error sending new category emoji prompt: %v", err)
		return
	}
}

func (sm *BotStateManager) promptSaveNewCategory(ctx context.Context, userID int64) {
	stateData, err := sm.getStateData(userID)
	if err != nil {
		log.Printf("[warn] error fetching state data: %v", err)
		return
	}

	category := storage.CategoryInfo{
		UserID: userID,
		Name:   stateData["NewCategoryNameEntered"].(string),
		Emoji:  stateData["NewCategoryEmojiEntered"].(string),
	}

	err = sm.Categories.AddOrUpdateCategory(category)
	if err != nil {
		log.Printf("[warn] error saving new category: %v", err)
		return
	}

	text := "Category saved!"
	err = sm.sendBotResponse(userID, text, nil)
	if err != nil {
		log.Printf("[warn] error sending new category save prompt: %v", err)
		return
	}

	if err := sm.UserFSMs[userID].Event(ctx, "SaveNewCategory"); err != nil {
		log.Printf("[error] Failed to transition to Idle state for user %d: %v", userID, err)
	}
}

func (sm *BotStateManager) sendBotResponse(chatID int64, text string, keyboard interface{}) error {
	tbMsg := tbapi.NewMessage(chatID, text)
	tbMsg.ParseMode = tbapi.ModeMarkdown
	tbMsg.DisableWebPagePreview = true
	tbMsg.ReplyMarkup = keyboard

	if keyboard == nil {
		removeKeyboard := tbapi.NewRemoveKeyboard(true)
		tbMsg.ReplyMarkup = removeKeyboard
	}

	if err := send(tbMsg, sm.TbAPI); err != nil {
		return fmt.Errorf("can't send message to telegram %s, %d: %w", text, chatID, err)
	}
	return nil
}

func (sm *BotStateManager) getStateData(userID int64) (map[string]interface{}, error) {
	currentStateInfo, err := sm.UserState.Read(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current state data for user %d: %w", userID, err)
	}

	data, err := unmarshalUserData(currentStateInfo.DataJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return data, nil
}
func unmarshalUserData(jsonData string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	if jsonData != "" {
		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}
		return data, nil
	}
	return make(map[string]interface{}), nil
}
