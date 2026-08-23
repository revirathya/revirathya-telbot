package utils

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const emptyCallback string = "em"

func NewEmptyInlineKeyboardButton(text string) tgbotapi.InlineKeyboardButton {
	emCallback := emptyCallback
	return tgbotapi.InlineKeyboardButton{
		Text:         text,
		CallbackData: &emCallback,
	}
}

func NewInlineTodayButton(label string, prefix string) tgbotapi.InlineKeyboardButton {
	dt := time.Now()
	return tgbotapi.NewInlineKeyboardButtonData(
		label,
		fmt.Sprintf("%s%s", prefix, dt.Format("2006-01-02")),
	)
}
