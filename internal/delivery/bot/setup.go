package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/auth"
	"github.com/revirathya/revirathya-telbot/internal/bot"
	"github.com/revirathya/revirathya-telbot/internal/bot/subflow"
	"github.com/revirathya/revirathya-telbot/internal/bot/updates"
	"github.com/revirathya/revirathya-telbot/internal/delivery"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
	"github.com/revirathya/revirathya-telbot/internal/usecase"
)

func RegisterHandler(
	ctx context.Context,
	srvLayer *delivery.ServerLayer,
	tbh *telegram.TelegramBotHandler,
	uc usecase.Usecase,
) error {
	// Set handler
	botHdl := bot.BootstrapBotHandler(tbh, uc)

	// Set calendar subflow
	sfCalendar, err := subflow.NewCalendar(tbh, uc.UserTelegramContext)
	if err != nil {
		return fmt.Errorf("failed when init calendar: %w", err)
	}

	// Set  numpad subflow
	sfNumpad := subflow.NewNumpad(tbh)
	srvLayer.Mux.HandleFunc(
		"POST /webhook",
		handleWebhookFunc(ctx, botHdl, uc, sfCalendar, sfNumpad),
	)
	return nil
}

func handleWebhookFunc(
	ctx context.Context,
	botHdl *bot.BotHandler,
	uc usecase.Usecase,
	sfCalendar *subflow.Calendar,
	sfNumpad *subflow.Numpad,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var tgUpt tgbotapi.Update
		if err := json.NewDecoder(r.Body).Decode(&tgUpt); err != nil {
			slog.Error("invalid body request: %s", "error", err.Error())
		}
		/*
			-- Hierarchies --
			Conversation (default: root)
			-> State (default: init)
				-> TextHandler
				-> CommandHandler
				-> CallbackHandler
		*/
		var logHandleInput string

		// Load context with telegram ID
		ctx = auth.WithTelegramID(ctx, tgUpt)

		// Get User Update
		upt, err := updates.GetUser(ctx, tgUpt, bot.UserLookupFunc(uc))
		if err != nil {
			slog.Error("get user", "error", err)
			return
		}

		// Cancel Command
		if (upt.Incoming.Mode == "command") && (upt.Incoming.Data == "cancel") {
			botHdl.CommandCancel.Command(ctx, upt)
		}

		// Handling Inputs (Update)
		slog.Info("body",
			"mode", upt.Incoming.Mode,
			"key", upt.Incoming.Key,
			"data", upt.Incoming.Data,
		)

		switch upt.Interaction.Conversation {
		case "root":
			// Conversation: Root
			if upt.Incoming.Mode == "command" {
				// Handler: Command
				switch upt.Incoming.Data {
				case "expense":
					logHandleInput = botHdl.Expense.Entry(ctx, upt)
				}
			}
		case "expense":
			// Conversation: Expense
			switch upt.Interaction.State {
			case "init":
				// State: init
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Data {
					case "expense-record":
						logHandleInput = botHdl.ExpenseRecord.Entry(ctx, upt)
					case "expense-list":
						logHandleInput = botHdl.ExpenseList.Entry(ctx, upt)
					case "expense-summary":
						logHandleInput = botHdl.ExpenseSummary.Entry(ctx, upt)
					}
				}
			}
		case "expense-record":
			// Conversation: Expense Record
			switch upt.Interaction.State {
			case "init":
				// State: init
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "date":
						logHandleInput = botHdl.ExpenseRecord.InteractionSelectCategory(ctx, upt)
					case "calendar":
						logHandleInput = sfCalendar.Init(ctx, "date", upt)
					}
				}
			case "calendar":
				// State: calendar
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "date":
						logHandleInput = botHdl.ExpenseRecord.InteractionSelectCategory(ctx, upt)
					case "calendar":
						logHandleInput = sfCalendar.Init(ctx, "date", upt)
					}
				}
			case "select-category":
				// State: select-category
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "category":
						logHandleInput = botHdl.ExpenseRecord.InteractionInputDescription(ctx, upt)
					}
				}
			case "input-description":
				// State: input-description
				if upt.Incoming.Mode == "text" {
					logHandleInput = botHdl.ExpenseRecord.InteractionInputAmount(ctx, upt)
				}
			case "input-amount":
				// State: input-amount
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "amount":
						logHandleInput = botHdl.ExpenseRecord.InteractionSelectConfirmation(ctx, upt)
					case "numpad":
						logHandleInput = sfNumpad.Init(ctx, "amount", upt)
					}
				}
			case "select-confirmation":
				// State: select-confirmation
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "confirmation":
						logHandleInput = botHdl.ExpenseRecord.Closure(ctx, upt)
					}
				}
			}
		case "expense-list":
			// Conversation: Expense List
			switch upt.Interaction.State {
			case "init":
				// State: init
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "date":
						logHandleInput = botHdl.ExpenseList.SelectExpense(ctx, upt)
					case "calendar":
						logHandleInput = sfCalendar.Init(ctx, "date", upt)
					}
				}
			case "calendar":
				// State: calendar
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "date":
						logHandleInput = botHdl.ExpenseList.SelectExpense(ctx, upt)
					case "calendar":
						logHandleInput = sfCalendar.Init(ctx, "date=", upt)
					}
				}
			case "select-expense":
				// State: select-expense
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "expense":
						logHandleInput = botHdl.ExpenseList.PreviewExpense(ctx, upt)
					case "done":
						logHandleInput = botHdl.ExpenseList.Closure(ctx, upt)
					case "calendar":
						logHandleInput = sfCalendar.Init(ctx, "date", upt)
					}
				}
			case "preview-expense":
				// State: preview-expense
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "delete":
						logHandleInput = botHdl.ExpenseList.ConfirmationDeleteExpense(ctx, upt)
					case "done":
						logHandleInput = botHdl.ExpenseList.Closure(ctx, upt)
					case "date":
						logHandleInput = botHdl.ExpenseList.SelectExpense(ctx, upt)
					}
				}
			case "delete-confirmation":
				// State: delete-confirmation
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "confirmation":
						if upt.Incoming.Data == "yes" {
							logHandleInput = botHdl.ExpenseList.AfterDeleteExpense(ctx, upt)
						} else {
							logHandleInput = botHdl.ExpenseList.PreviewExpense(ctx, upt)
						}
					}
				}
			case "after-delete":
				// State: after-delete
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "done":
						logHandleInput = botHdl.ExpenseList.Closure(ctx, upt)
					case "calendar":
						logHandleInput = sfCalendar.Init(ctx, "date", upt)
					}
				}
			}
		case "expense-summary":
			// Conversation: Expense Summary
			switch upt.Interaction.State {
			case "init":
				// State: init
				if upt.Incoming.Mode == "callback" {
					if (upt.Incoming.Key == "summary") && (upt.Incoming.Data == "daily-expense") {
						logHandleInput = botHdl.ExpenseSummary.DailyExpenseFilterDate(ctx, upt)
					}
				}
			case "daily-expense-filter-date":
				// State: init
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "date":
						logHandleInput = botHdl.ExpenseSummary.DailyExpenseReport(ctx, upt)
					case "calendar":
						logHandleInput = sfCalendar.Init(ctx, "date", upt)
					}
				}
			case "daily-expense-report":
				// State: init
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "done":
						logHandleInput = botHdl.ExpenseSummary.Closure(ctx, upt)
					case "calendar":
						logHandleInput = sfCalendar.Init(ctx, "date", upt)
					}
				}
			case "calendar":
				// State: calendar
				if upt.Incoming.Mode == "callback" {
					switch upt.Incoming.Key {
					case "date":
						logHandleInput = botHdl.ExpenseSummary.DailyExpenseReport(ctx, upt)
					case "calendar":
						logHandleInput = sfCalendar.Init(ctx, "date", upt)
					}
				}
			}
		}

		slog.Info("user_context",
			"user.conversation", upt.Interaction.Conversation,
			"user.state", upt.Interaction.State,
			"log_handle_input", logHandleInput,
		)
	}
}
