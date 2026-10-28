package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/bot/container"
	"github.com/revirathya/revirathya-telbot/internal/bot/subflow"
	"github.com/revirathya/revirathya-telbot/internal/bot/updates"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
	"github.com/revirathya/revirathya-telbot/internal/utils"
)

type ConversationExpenseRecord struct {
	telbot              *telegram.TelegramBotHandler
	prepExp             *container.ExpenseContainer
	userTelegramContext userTelegramContextUsecase
	category            categoryUsecase
	expense             expenseUsecase
}

func NewConversationExpenseRecord(
	tbh *telegram.TelegramBotHandler,
	userTelegramContext userTelegramContextUsecase,
	category categoryUsecase,
	expense expenseUsecase,
) *ConversationExpenseRecord {
	return &ConversationExpenseRecord{
		telbot:              tbh,
		prepExp:             &container.ExpenseContainer{},
		userTelegramContext: userTelegramContext,
		category:            category,
		expense:             expense,
	}
}

func (c *ConversationExpenseRecord) Entry(ctx context.Context, upt updates.UserUpdate) string {
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
	msg, err := c.telbot.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, &inlineKeyboardMarkup)
	if err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.userTelegramContext.StartConversation(ctx, "expense-record", msg.Chat.ID, msg.MessageID); err != nil {
		slog.Error("failed to start conversation", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}

func (c *ConversationExpenseRecord) InteractionSelectCategory(ctx context.Context, upt updates.UserUpdate) string {
	// Update Create Expense Context
	if err := c.prepExp.SetTransactionDtFromString(upt.Incoming.Data); err != nil {
		slog.Error("failed to set transaction_dt", "error", err)
	}

	// Get Categories
	categories, err := c.category.ListCategories(ctx)
	if err != nil {
		slog.Error("failed to list categories")
	}

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
	if _, err := c.telbot.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.userTelegramContext.ChangeState(ctx, "select-category"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}

func (c *ConversationExpenseRecord) InteractionInputDescription(ctx context.Context, upt updates.UserUpdate) string {
	// Update Create Expense Context
	categoryID, err := uuid.Parse(upt.Incoming.Data)
	if err != nil {
		slog.Error("invalid category_id UUID format", "error", err)
	}

	category, err := c.category.GetCategoryByID(ctx, categoryID)
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
	if _, err := c.telbot.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, nil); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.userTelegramContext.ChangeState(ctx, "input-description"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}

func (c *ConversationExpenseRecord) InteractionInputAmount(ctx context.Context, upt updates.UserUpdate) string {
	// Update Create Expense Context
	c.prepExp.SetDescription(upt.Incoming.Data)

	// Set Text
	text := strings.Join([]string{
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
	if _, err := c.telbot.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.userTelegramContext.ChangeState(ctx, "input-amount"); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}
	return "text"
}

func (c *ConversationExpenseRecord) InteractionSelectConfirmation(ctx context.Context, upt updates.UserUpdate) string {
	// Update Create Expense Context
	if err := c.prepExp.SetAmountFromString(upt.Incoming.Data); err != nil {
		slog.Error("failed to set amount", "error", err)
	}

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
	if _, err := c.telbot.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Set User Context
	if _, err := c.userTelegramContext.ChangeState(ctx, "select-confirmation"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}

func (c *ConversationExpenseRecord) Closure(ctx context.Context, upt updates.UserUpdate) string {
	confirmation := upt.Incoming.Data

	var text string
	switch confirmation {
	case "yes":
		// Create Expense
		// Record expenses to storage
		expense, err := c.prepExp.ToDomain()
		if err != nil {
			slog.Error("failed to convert container to expense", "error", err)
		}

		_, err = c.expense.CreateExpense(ctx, expense, uuid.New())
		if err != nil {
			slog.Error("failed to create expense", "error", err)
		}

		// Set Text
		text = strings.Join([]string{
			"Successfully record following expense:",
			c.prepExp.FormatDetail(),
		}, "\n\n")

		// Edit Message
		if _, err := c.telbot.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, nil); err != nil {
			slog.Error("failed to edit telegram message", "error", err)
		}
	case "no":
		// Cancel recording
		text = "Cancel record new expense"
		if _, err := c.telbot.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, nil); err != nil {
			slog.Error("failed to edit telegram message", "error", err)
		}
	}

	// Clear All Context
	c.prepExp.Clear()
	if _, err := c.userTelegramContext.CloseConversation(ctx); err != nil {
		slog.Error("failed to close conversation", "error", err)
	}

	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}
