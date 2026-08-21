package utils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const EMPTY_CALLBACK string = "em"

func NewEmptyInlineKeyboardButton(text string) tgbotapi.InlineKeyboardButton {
	emCallback := EMPTY_CALLBACK
	return tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &emCallback,
	}
}
