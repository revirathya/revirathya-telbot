package subflow

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
	"github.com/revirathya/revirathya-telbot/internal/utils"
)

type Numpad struct {
	telbot *telegram.TelegramBotHandler
}

func NewNumpad(tbh *telegram.TelegramBotHandler) *Numpad {
	return &Numpad{
		telbot: tbh,
	}
}

// Private
func (n *Numpad) newNumpadAmount(currAmount string) ([]tgbotapi.InlineKeyboardButton, error) {
	amountInt, err := strconv.ParseInt(currAmount, 10, 64)
	if err != nil {
		return nil, err
	}

	numpadAmount := tgbotapi.NewInlineKeyboardRow(
		utils.NewEmptyInlineKeyboardButton(
			utils.FormatThousand(amountInt),
		),
	)
	return numpadAmount, nil
}

func (n *Numpad) newNumpadNumbers(currAmount string) [][]tgbotapi.InlineKeyboardButton {
	body := make([][]tgbotapi.InlineKeyboardButton, 0)

	var iter int64 = 1
	var callbackData string
	row := make([]tgbotapi.InlineKeyboardButton, 0)
	for range 3 {
		row = nil
		for range 3 {
			if currAmount != "0" {
				callbackData = fmt.Sprintf("numpad=%s%d", currAmount, iter)
			} else {
				callbackData = fmt.Sprintf("numpad=%d", iter)
			}
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(
				strconv.FormatInt(iter, 10),
				callbackData,
			))
			iter += 1
		}
		body = append(body, row)
	}

	return body
}

func (n *Numpad) newZeroUtility(callbackPrefix string, currAmount string) [][]tgbotapi.InlineKeyboardButton {
	npRemove := "numpad=0"
	if len(currAmount) > 1 {
		npRemove = fmt.Sprintf("numpad=%s", currAmount[:len(currAmount)-1])
	}

	utility := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("000", fmt.Sprintf("numpad=%s000", currAmount)),
			tgbotapi.NewInlineKeyboardButtonData("0", fmt.Sprintf("numpad=%s0", currAmount)),
			tgbotapi.NewInlineKeyboardButtonData("C", npRemove),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("CE", "numpad=0"),
			tgbotapi.NewInlineKeyboardButtonData("OK", fmt.Sprintf("%s%s", callbackPrefix, currAmount)),
			utils.NewEmptyInlineKeyboardButton(" "),
		),
	}
	return utility
}

// Public
func (n *Numpad) Init(ctx context.Context, user *domain.UserTelegramContext, callbackPrefix string, currAmount string, data string) string {
	// Set Text
	text := "Please input amount below:"

	// Build Numpad
	npAmount, err := n.newNumpadAmount(currAmount)
	if err != nil {
		return ""
	}

	npNumbers := n.newNumpadNumbers(currAmount)
	npZeroUtility := n.newZeroUtility(callbackPrefix, currAmount)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	rows = append(rows, npAmount)
	rows = append(rows, npNumbers...)
	rows = append(rows, npZeroUtility...)

	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// Edit Message
	if _, err := n.telbot.EditMessage(*user.InteractionChatID, *user.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	return fmt.Sprintf("callback:%s", data)
}
