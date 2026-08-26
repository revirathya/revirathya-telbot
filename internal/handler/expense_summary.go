package handler

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/handler/subflow"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
	"github.com/revirathya/revirathya-telbot/internal/utils"
)

type ReportOption struct {
	Code  string
	Label string
}

type ConversationExpenseSummary struct {
	Telbot        *telegram.TelegramBotHandler
	User          *domain.UserContext
	ReportOptions []ReportOption
}

func NewConversationExpenseSummary(tbh *telegram.TelegramBotHandler, user *domain.UserContext) *ConversationExpenseSummary {
	return &ConversationExpenseSummary{
		Telbot: tbh,
		User:   user,
		ReportOptions: []ReportOption{
			{
				Code:  "daily-expense",
				Label: "Daily Expense",
			},
		},
	}
}

// Private
func (c *ConversationExpenseSummary) genDailyExpense(expenseList []domain.Expense, date string) string {
	var totalExpense int64 = 0

	existingCategory := make([]string, 0)
	for _, e := range expenseList {
		totalExpense += e.Amount
		if !slices.Contains(existingCategory, e.CategoryLabel) {
			existingCategory = append(existingCategory, e.CategoryLabel)
		}
	}

	var idx int
	var categoryAmountReportRows []string
	var categoryReport string

	slices.Sort(existingCategory)
	if len(existingCategory) > 0 {
		categoryAmount := make([]int64, len(existingCategory))
		for _, e := range expenseList {
			idx = slices.IndexFunc(existingCategory, func(c string) bool {
				return c == e.CategoryLabel
			})
			categoryAmount[idx] += e.Amount
		}

		for i, c := range existingCategory {
			categoryAmountReportRows = append(
				categoryAmountReportRows,
				fmt.Sprintf("[%s] Rp %s", c, utils.FormatThousand(categoryAmount[i])),
			)
		}

		categoryReport = strings.Join(categoryAmountReportRows, "\n")
	} else {
		categoryReport = ""
	}

	// Report
	reportRows := make([]string, 0)
	reportRows = append(reportRows, fmt.Sprintf("Total Expense: Rp %s", utils.FormatThousand(totalExpense)))
	if len(existingCategory) > 0 {
		reportRows = append(reportRows, categoryReport)
	}

	report := strings.Join(reportRows, "\n\n")
	return report
}

// Public
func (c *ConversationExpenseSummary) Entry(chatID int64, messageID int) string {
	// Set Text
	text := "Please select from options of summary reports below"

	// Set Reply Markup
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(c.ReportOptions))
	for _, ro := range c.ReportOptions {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ro.Label,
				fmt.Sprintf("summary=%s", ro.Code),
			),
		))
	}

	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// Edit Message
	if _, err := c.Telbot.EditMessage(chatID, c.User.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Update User Context
	c.User.StartConversation("expense-summary", chatID, messageID)
	return "callback:expense-summary"
}

func (c *ConversationExpenseSummary) DailyExpenseFilterDate() string {
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
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Update User Context
	c.User.ChangeState("daily-expense-filter-date")
	return "callback:report=daily-expense"
}

func (c *ConversationExpenseSummary) DailyExpenseReport(expenseList []domain.Expense, data string) string {
	date := strings.Split(data, "=")[1]

	// Set Text
	text := strings.Join([]string{
		fmt.Sprintf("Daily Expense - %s", date),
		c.genDailyExpense(expenseList, date),
	}, "\n\n")

	// Set Inline Keyboard
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			subflow.NewInlineCalendarButton("Pick Other Date", nil),
			tgbotapi.NewInlineKeyboardButtonData("Done", "done"),
		),
	)

	// Edit Message
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Update User Context
	c.User.ChangeState("daily-expense-report")
	return "callback:report=daily-expense"
}

func (c *ConversationExpenseSummary) Closure() string {
	// Set Text
	text := "Closing expense summary"

	// Edit Message
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, nil); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Update User Context
	c.User.CloseConversation()
	return "callback:done"
}
