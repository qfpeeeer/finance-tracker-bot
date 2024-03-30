package keyboards

import "github.com/nyanyamaga/finance-tracker-bot/app/storage"

type TbKeyboardProvider struct {
	Storage *storage.Category
}

func NewTbKeyboardProvider(storage *storage.Category) *TbKeyboardProvider {
	return &TbKeyboardProvider{Storage: storage}
}
