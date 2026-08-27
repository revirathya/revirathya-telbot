package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/handler/container"
	"github.com/revirathya/revirathya-telbot/internal/handler/subflow"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
	ucCtg "github.com/revirathya/revirathya-telbot/internal/usecase/category"
	ucExp "github.com/revirathya/revirathya-telbot/internal/usecase/expense"
	ucUtel "github.com/revirathya/revirathya-telbot/internal/usecase/usercontext"
	"github.com/revirathya/revirathya-telbot/internal/utils"
)

type ConversationExpenseRecord struct {
	telbot  *telegram.TelegramBotHandler
	prepExp *container.ExpenseContainer
	utelUc  *ucUtel.UserTelegramContextUsecase
	ctgUc   *ucCtg.CategoryUsecase
	expUc   *ucExp.ExpenseUsecase
}

func NewConversationExpenseRecord(tbh *telegram.TelegramBotHandler, utelUc *ucUtel.UserTelegramContextUsecase, ctgUc *ucCtg.CategoryUsecase, expUc *ucExp.ExpenseUsecase) *ConversationExpenseRecord {
	return &ConversationExpenseRecord{
		telbot:  tbh,
		prepExp: &container.ExpenseContainer{},
		utelUc:  utelUc,
		ctgUc:   ctgUc,
		expUc:   expUc,
	}
}

func (c *ConversationExpenseRecord) Entry(ctx context.Context, user *domain.UserTelegramContext, chatID int64, messageID int, data string) string {
	// Set Text
	text := strings.Join([]string{
		"<b>Record Expense</b>",
		"Please select date:",
	}, "\n\n")

	// Set Reply Markup
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			utils.NewInlineTodayButton("Today", "date="),
			subflow.NewInlineCalendarButton("Pick a Date", nil),
		),
	)

	// Edit Message
	if _, err := c.telbot.EditMessage(chatID, messageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.utelUc.StartConversation(ctx, user.UserID, "expense-record", chatID, messageID); err != nil {
		slog.Error("failed to start conversation", "error", err)
	}
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseRecord) InteractionSelectCategory(ctx context.Context, user *domain.UserTelegramContext, data string) string {
	// Update Create Expense Context
	inpDate := strings.Split(data, "=")[1]
	c.prepExp.SetTransactionDate(inpDate)

	// Get Categories
	categories, _ := c.ctgUc.ListCategories(ctx, user.UserID)

	// Set Text
	text := strings.Join([]string{
		"<b>Record Expense</b>",
		c.prepExp.FormatDetail(),
		"Please select category:",
	}, "\n\n")

	// Set Reply Markup
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, ctg := range categories {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ctg.CategoryName,
				fmt.Sprintf("category=%s", ctg.ID),
			),
		))
	}
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// Edit Message
	if _, err := c.telbot.EditMessage(*user.InteractionChatID, *user.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.utelUc.ChangeState(ctx, user.UserID, "select-category"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseRecord) InteractionInputDescription(ctx context.Context, user *domain.UserTelegramContext, data string) string {
	// Update Create Expense Context
	categoryID := strings.Split(data, "=")[1]
	uuidCategoryID, err := uuid.Parse(categoryID)
	if err != nil {
		slog.Error("invalid category_id UUID format", "error", err)
	}

	category, err := c.ctgUc.GetCategoryByID(ctx, uuidCategoryID)
	if err != nil {
		slog.Error("failed to get category by id", "error", err)
	}
	c.prepExp.SetCategory(category.ID.String(), category.CategoryName)

	// Set text
	text := strings.Join([]string{
		"<b>Record Expense</b>",
		c.prepExp.FormatDetail(),
		"Please input description:",
	}, "\n\n")

	// Edit Message
	if _, err := c.telbot.EditMessage(*user.InteractionChatID, *user.InteractionMessageID, text, nil); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.utelUc.ChangeState(ctx, user.UserID, "input-description"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseRecord) InteractionInputAmount(ctx context.Context, user *domain.UserTelegramContext, text string) string {
	// Update Create Expense Context
	c.prepExp.SetDescription(text)

	// Set Text
	text = strings.Join([]string{
		"<b>Record Expense</b>",
		c.prepExp.FormatDetail(),
		"Please input amount:",
	}, "\n\n")

	// Set Reply Markup
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Input Amount", "numpad=0"),
		),
	)

	// Edit Message
	if _, err := c.telbot.EditMessage(*user.InteractionChatID, *user.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.utelUc.ChangeState(ctx, user.UserID, "input-amount"); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}
	return "text"
}

func (c *ConversationExpenseRecord) InteractionSelectConfirmation(ctx context.Context, user *domain.UserTelegramContext, data string) string {
	// Update Create Expense Context
	amount := strings.Split(data, "=")[1]
	c.prepExp.SetAmount(amount)

	// Set Text
	text := strings.Join([]string{
		"<b>Record Expense</b>",
		c.prepExp.FormatDetail(),
		"Do you want record this expense?",
	}, "\n\n")

	// Set Reply Markup
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Yes", "confirmation=yes"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("No", "confirmation=no"),
		),
	)

	// Edit Message
	if _, err := c.telbot.EditMessage(*user.InteractionChatID, *user.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Set User Context
	if _, err := c.utelUc.ChangeState(ctx, user.UserID, "select-confirmation"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", data)
}

func (c *ConversationExpenseRecord) Closure(ctx context.Context, user *domain.UserTelegramContext, data string) string {
	var text string
	confirmation := strings.Split(data, "=")[1]

	switch confirmation {
	case "yes":
		// Create Expense
		// Record expenses to storage
		expense, err := c.prepExp.ToDomain(user.UserID)
		if err != nil {
			slog.Error("failed to convert container to expense", "error", err)
		}

		_, err = c.expUc.CreateExpense(ctx, *expense)
		if err != nil {
			slog.Error("failed to create expense", "error", err)
		}

		// Set Text
		text = strings.Join([]string{
			"Successfully record following expense:",
			c.prepExp.FormatDetail(),
		}, "\n\n")
	case "no":
		// Set Text
		text = "Cancel record new expense"
	}

	// Edit Message
	if _, err := c.telbot.EditMessage(*user.InteractionChatID, *user.InteractionMessageID, text, nil); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Clear All Context
	c.prepExp.Clear()
	if _, err := c.utelUc.CloseConversation(ctx, user.UserID); err != nil {
		slog.Error("failed to close conversation", "error", err)
	}
	return fmt.Sprintf("callback:%s", data)
}
