package subflow

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revirathya/revirathya-telbot/infra/telegram"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/utils"
)

type Calendar struct {
	Telbot *telegram.TelegramBotHandler
	User   *domain.UserContext
}

func NewCalendar(tbh *telegram.TelegramBotHandler, user *domain.UserContext) *Calendar {
	return &Calendar{
		Telbot: tbh,
		User:   user,
	}
}

func (c *Calendar) Init(callbackPrefix string) string {
	// Set Text
	text := "Please select date below:"

	// Build Calendar
	calHeader := [][]tgbotapi.InlineKeyboardButton{{{Text: "Aug 2026", CallbackData: utils.EmptyCallback()}}}
	calBody := [][]tgbotapi.InlineKeyboardButton{
		{
			{Text: " ", CallbackData: utils.EmptyCallback()},
			{Text: " ", CallbackData: utils.EmptyCallback()},
			{Text: " ", CallbackData: utils.EmptyCallback()},
			{Text: " ", CallbackData: utils.EmptyCallback()},
			{Text: " ", CallbackData: utils.EmptyCallback()},
			{Text: "01", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-01")},
			{Text: "02", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-02")},
		},
		{
			{Text: "03", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-03")},
			{Text: "04", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-04")},
			{Text: "05", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-05")},
			{Text: "06", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-06")},
			{Text: "07", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-07")},
			{Text: "08", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-08")},
			{Text: "09", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-09")},
		},
		{
			{Text: "10", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-10")},
			{Text: "11", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-11")},
			{Text: "12", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-12")},
			{Text: "13", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-13")},
			{Text: "14", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-14")},
			{Text: "15", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-15")},
			{Text: "16", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-16")},
		},
		{
			{Text: "17", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-17")},
			{Text: "18", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-18")},
			{Text: "19", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-19")},
			{Text: "20", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-20")},
			{Text: "21", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-21")},
			{Text: "22", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-22")},
			{Text: "23", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-23")},
		},
		{
			{Text: "24", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-24")},
			{Text: "25", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-25")},
			{Text: "26", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-26")},
			{Text: "27", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-27")},
			{Text: "28", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-28")},
			{Text: "29", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-29")},
			{Text: "30", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-30")},
		},
		{
			{Text: "31", CallbackData: utils.ToStrPtr(callbackPrefix + "2026-08-31")},
			{Text: " ", CallbackData: utils.EmptyCallback()},
			{Text: " ", CallbackData: utils.EmptyCallback()},
			{Text: " ", CallbackData: utils.EmptyCallback()},
			{Text: " ", CallbackData: utils.EmptyCallback()},
			{Text: " ", CallbackData: utils.EmptyCallback()},
			{Text: " ", CallbackData: utils.EmptyCallback()},
		},
	}
	calNav := [][]tgbotapi.InlineKeyboardButton{
		{
			{Text: "<", CallbackData: utils.ToStrPtr("date-picker=prev")},
			{Text: " ", CallbackData: utils.EmptyCallback()},
			{Text: ">", CallbackData: utils.ToStrPtr("date-picker=next")},
		},
	}

	// Set Reply Markup
	inlineKeyboard := make([][]tgbotapi.InlineKeyboardButton, 0)
	inlineKeyboard = append(inlineKeyboard, calHeader...)
	inlineKeyboard = append(inlineKeyboard, calBody...)
	inlineKeyboard = append(inlineKeyboard, calNav...)

	inlineKeyboardMarkup := &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}

	// Edit Message
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Set User Context
	c.User.ChangeState("calendar")
	return "callback:calendar"
}
