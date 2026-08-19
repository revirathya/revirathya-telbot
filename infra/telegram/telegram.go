package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBotHandler struct {
	Bot     *tgbotapi.BotAPI
	Webhook tgbotapi.WebhookConfig
	Updates tgbotapi.UpdatesChannel
}

func GetBotWebhook(botToken string, webhookURL string) (*TelegramBotHandler, error) {
	// Authorize Bot
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}
	// bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Set Webhook
	wh, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		return nil, err
	}
	wh.DropPendingUpdates = true

	// Request Bot
	if _, err = bot.Request(wh); err != nil {
		return nil, err
	}

	// Get Info Bot
	info, err := bot.GetWebhookInfo()
	if err != nil {
		return nil, err
	}

	if info.LastErrorDate != 0 {
		err_msg := fmt.Errorf("Telegram callback failed: %s", info.LastErrorMessage)
		return nil, err_msg
	}

	// Get updates
	updates := bot.ListenForWebhook("/")

	// Return pointer
	return &TelegramBotHandler{
		Bot:     bot,
		Webhook: wh,
		Updates: updates,
	}, nil
}

func (tbh *TelegramBotHandler) SendMessage(chatID int64, text string, replyMarkup interface{}) (tgbotapi.Message, error) {
	msg := &tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: replyMarkup,
		},
		Text:                  text,
		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: false,
	}

	return tbh.Bot.Send(msg)
}

func (tbh *TelegramBotHandler) EditMessage(chatID int64, messageID int, text string, inlineKeyboardMarkup *tgbotapi.InlineKeyboardMarkup) (tgbotapi.Message, error) {
	editMsg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatID,
			MessageID:   messageID,
			ReplyMarkup: inlineKeyboardMarkup,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	}

	return tbh.Bot.Send(editMsg)
}
