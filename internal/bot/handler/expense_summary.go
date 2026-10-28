package handler

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/bot/subflow"
	"github.com/revirathya/revirathya-telbot/internal/bot/updates"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
	"github.com/revirathya/revirathya-telbot/internal/utils"
	"github.com/shopspring/decimal"
)

type ReportOption struct {
	Code  string
	Label string
}

type ConversationExpenseSummary struct {
	telbot              *telegram.TelegramBotHandler
	reportOptions       []ReportOption
	userTelegramContext userTelegramContextUsecase
	expense             expenseUsecase
}

func NewConversationExpenseSummary(
	tbh *telegram.TelegramBotHandler,
	userTelegramContext userTelegramContextUsecase,
	expense expenseUsecase,
) *ConversationExpenseSummary {
	return &ConversationExpenseSummary{
		telbot: tbh,
		reportOptions: []ReportOption{
			{
				Code:  "daily-expense",
				Label: "Daily Expense",
			},
		},
		userTelegramContext: userTelegramContext,
		expense:             expense,
	}
}

// Private
func (c *ConversationExpenseSummary) genDailyExpense(expenses []domain.Expense) string {
	var totalExpense = decimal.New(0, 0)

	existingCategory := make([]string, 0)
	for _, e := range expenses {
		totalExpense = totalExpense.Add(e.Amount)
		if !slices.Contains(existingCategory, e.CategoryName) {
			existingCategory = append(existingCategory, e.CategoryName)
		}
	}

	var idx int
	var categoryAmountReportRows []string
	var categoryReport string

	slices.Sort(existingCategory)
	if len(existingCategory) > 0 {
		categoryAmount := make([]decimal.Decimal, len(existingCategory))
		for _, e := range expenses {
			idx = slices.IndexFunc(existingCategory, func(c string) bool {
				return c == e.CategoryName
			})
			categoryAmount[idx] = categoryAmount[idx].Add(e.Amount)
		}

		for i, c := range existingCategory {
			categoryAmountReportRows = append(
				categoryAmountReportRows,
				fmt.Sprintf("[%s] Rp %s", c, utils.FormatThousand(categoryAmount[i].IntPart())),
			)
		}

		categoryReport = strings.Join(categoryAmountReportRows, "\n")
	} else {
		categoryReport = ""
	}

	// Report
	reportRows := make([]string, 0)
	reportRows = append(reportRows, fmt.Sprintf("Total Expense: Rp %s", utils.FormatThousand(totalExpense.IntPart())))
	if len(existingCategory) > 0 {
		reportRows = append(reportRows, categoryReport)
	}

	report := strings.Join(reportRows, "\n\n")
	return report
}

// Public
func (c *ConversationExpenseSummary) Entry(ctx context.Context, upt updates.UserUpdate) string {
	// Set Text
	text := "Please select from options of summary reports below"

	// Set Reply Markup
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(c.reportOptions))
	for _, ro := range c.reportOptions {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ro.Label,
				fmt.Sprintf("summary=%s", ro.Code),
			),
		))
	}

	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// Edit Message
	if _, err := c.telbot.EditMessage(upt.Incoming.ChatID, upt.Incoming.MessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.userTelegramContext.StartConversation(ctx, "expense-summary", upt.Incoming.ChatID, upt.Incoming.MessageID); err != nil {
		slog.Error("failed to start conversation", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}

func (c *ConversationExpenseSummary) DailyExpenseFilterDate(ctx context.Context, upt updates.UserUpdate) string {
	// Set Text
	text := "Please select expense date:"

	// Set Inline Keyboard
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			utils.NewInlineTodayButton("Today", "date="),
			subflow.NewInlineCalendarButton("Pick a Date", nil),
		),
	)

	// Edit Message
	if _, err := c.telbot.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.userTelegramContext.ChangeState(ctx, "daily-expense-filter-date"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}

func (c *ConversationExpenseSummary) DailyExpenseReport(ctx context.Context, upt updates.UserUpdate) string {
	// Get Expense
	trxDate, err := time.Parse("2006-01-02", upt.Incoming.Data)
	if err != nil {
		slog.Error("invalid transaction_dt Date format", "error", err)
	}

	expenses, err := c.expense.ListExpenses(ctx, trxDate)
	if err != nil {
		slog.Error("failed to list expenses", "error", err)
	}

	// Set Text
	text := strings.Join([]string{
		fmt.Sprintf("Daily Expense - %s", trxDate.Format("2006-01-02")),
		c.genDailyExpense(expenses),
	}, "\n\n")

	// Set Inline Keyboard
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			subflow.NewInlineCalendarButton("Pick Other Date", nil),
			tgbotapi.NewInlineKeyboardButtonData("Done", "done"),
		),
	)

	// Edit Message
	if _, err := c.telbot.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.userTelegramContext.ChangeState(ctx, "daily-expense-report"); err != nil {
		slog.Error("failed to change state", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}

func (c *ConversationExpenseSummary) Closure(ctx context.Context, upt updates.UserUpdate) string {
	// Set Text
	text := "Closing expense summary"

	// Edit Message
	if _, err := c.telbot.EditMessage(upt.Interaction.ChatID, upt.Interaction.MessageID, text, nil); err != nil {
		slog.Error("failed to edit telegram message", "error", err)
	}

	// Update User Context
	if _, err := c.userTelegramContext.CloseConversation(ctx); err != nil {
		slog.Error("failed to close conversation", "error", err)
	}
	return fmt.Sprintf("callback:%s", upt.Incoming.Data)
}
