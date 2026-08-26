package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/revirathya/revirathya-telbot/internal/config"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/handler"
	"github.com/revirathya/revirathya-telbot/internal/handler/subflow"
	grpcCon "github.com/revirathya/revirathya-telbot/internal/platform/grpc"
	"github.com/revirathya/revirathya-telbot/internal/platform/telegram"
)

// Main
func main() {
	// Set logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Get Telegram Bot
	tbh, err := telegram.GetBotWebhook(cfg.TelegramBotToken, cfg.TelegramWebhookURL)
	if err != nil {
		slog.Error("failed to get webhook", "error", err)
		os.Exit(1)
	}

	// Set User context
	user := domain.NewUserContext()
	expenseList := make([]domain.Expense, 0)

	// Data: Categories
	var CATEGORIES []domain.Category = []domain.Category{
		{ID: "1", Label: "Food - Daily"},
		{ID: "2", Label: "Food - Snacks & Coffee"},
		{ID: "3", Label: "Groceries"},
	}

	// Set gRPC Servers
	fincoreConn, err := grpcCon.NewFinanceCoreConn(cfg.FinanceCoreGRPC)
	if err != nil {
		slog.Error("failed to connect finance core grpc server", "error", err)
		os.Exit(1)
	}
	defer fincoreConn.Close()

	// Conversation: Handler
	convExpense := handler.NewConversationExpense(tbh, user)
	convExpenseRecord := handler.NewConversationExpenseRecord(tbh, user)
	convExpenseList := handler.NewConversationExpenseList(tbh, user)
	convExpenseSummary := handler.NewConversationExpenseSummary(tbh, user)

	// Subflow: Calendar
	sfCalendar, err := subflow.NewCalendar(tbh, user)
	if err != nil {
		slog.Error("failed when init calendar: %s", "error", err.Error())
		return
	}

	// Subflow: Numpad
	sfNumpad := subflow.NewNumpad(tbh, user)

	// API handler
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
		var reqBody tgbotapi.Update
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			slog.Error("invalid body request: %s", "error", err.Error())
			return
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

		// Handling Updates from Telegram Bot
		var chatID int64
		var text string
		var command string

		if (reqBody.Message != nil) && (reqBody.Message.Text == "/cancel") {
			handler.CommandCancel(tbh, user, reqBody.Message.Chat.ID)
		}

		// Handling Inputs (Update)
		switch user.Conversation {
		case "root":
			// Conversation: Root
			if reqBody.Message != nil {
				// Handler: Command & Text
				chatID = reqBody.Message.Chat.ID
				text = reqBody.Message.Text

				if text[0] == '/' {
					// Handler: Command
					command = text[1:]
					switch command {
					case "expense":
						logHandleInput = convExpense.Entry(chatID)
					}
				}
			}
		case "expense":
			// Conversation: Expense
			switch user.State {
			case "init":
				// State: init
				if reqBody.CallbackQuery != nil {
					// Handler: CallbackQuery
					chatID := reqBody.CallbackQuery.Message.Chat.ID
					messageID := reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					switch data {
					case "expense-record":
						logHandleInput = convExpenseRecord.Entry(chatID, messageID)
					case "expense-list":
						logHandleInput = convExpenseList.Entry(chatID, messageID)
					case "expense-summary":
						logHandleInput = convExpenseSummary.Entry(chatID, messageID)
					}
				}
			}
		case "expense-record":
			// Conversation: Expense Record
			switch user.State {
			case "init":
				// State: init
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseRecord.InteractionSelectCategory(CATEGORIES, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init("date=", refDate)
					}
				}
			case "calendar":
				// State: calendar
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseRecord.InteractionSelectCategory(CATEGORIES, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init("date=", refDate)
					}
				}
			case "select-category":
				// State: select-category
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					_ = reqBody.CallbackQuery.Message.Chat.ID
					_ = reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "category=") {
						logHandleInput = convExpenseRecord.InteractionInputDescription(CATEGORIES, data)
					}
				}
			case "input-description":
				// State: input-description
				if reqBody.Message != nil {
					// Handler: Command & Text
					text := reqBody.Message.Text

					logHandleInput = convExpenseRecord.InteractionInputAmount(text)
				}
			case "input-amount":
				// State: input-amount
				if reqBody.CallbackQuery != nil {
					// Handler: CallbackQuery
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "amount=") {
						logHandleInput = convExpenseRecord.InteractionSelectConfirmation(data)
					} else if strings.HasPrefix(data, "numpad=") {
						currAmount := strings.Split(data, "=")[1]
						logHandleInput = sfNumpad.Init("amount=", currAmount)
					}
				}
			case "select-confirmation":
				// State: select-confirmation
				if reqBody.CallbackQuery != nil {
					// Handler: CallbackQuery
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "confirmation=") {
						logHandleInput = convExpenseRecord.Closure(&expenseList, data)
					}
				}
			}
		case "expense-list":
			// Conversation: Expense List
			switch user.State {
			case "init":
				// State: init
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					_ = reqBody.CallbackQuery.Message.Chat.ID
					_ = reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseList.SelectExpense(expenseList, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init("date=", refDate)
					}
				}
			case "calendar":
				// State: calendar
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseList.SelectExpense(expenseList, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init("date=", refDate)
					}
				}
			case "select-expense":
				// State: select-expense
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					_ = reqBody.CallbackQuery.Message.Chat.ID
					_ = reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "expense=") {
						logHandleInput = convExpenseList.PreviewExpense(expenseList, data, false)
					} else if data == "done" {
						logHandleInput = convExpenseList.Closure()
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init("date=", refDate)
					}
				}
			case "preview-expense":
				// State: preview-expense
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					_ = reqBody.CallbackQuery.Message.Chat.ID
					_ = reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					if data == "delete" {
						logHandleInput = convExpenseList.ConfirmationDeleteExpense()
					} else if data == "done" {
						logHandleInput = convExpenseList.Closure()
					} else if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseList.SelectExpense(expenseList, data)
					}
				}
			case "delete-confirmation":
				// State: delete-confirmation
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					_ = reqBody.CallbackQuery.Message.Chat.ID
					_ = reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					switch data {
					case "confirmation=yes":
						logHandleInput = convExpenseList.AfterDeleteExpense(expenseList)
					case "confirmation=no":
						logHandleInput = convExpenseList.PreviewExpense(nil, data, true)
					}
				}
			case "after-delete":
				// State: after-delete
				if reqBody.CallbackQuery != nil {
					// handler: Callback Query
					_ = reqBody.CallbackQuery.Message.Chat.ID
					_ = reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					if data == "done" {
						logHandleInput = convExpenseList.Closure()
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init("date=", refDate)
					}
				}
			}
		case "expense-summary":
			// Conversation: Expense Summary
			switch user.State {
			case "init":
				// State: init
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					_ = reqBody.CallbackQuery.Message.Chat.ID
					_ = reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					if data == "summary=daily-expense" {
						logHandleInput = convExpenseSummary.DailyExpenseFilterDate()
					}
				}
			case "daily-expense-filter-date":
				// State: init
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					_ = reqBody.CallbackQuery.Message.Chat.ID
					_ = reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseSummary.DailyExpenseReport(expenseList, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init("date=", refDate)
					}
				}
			case "daily-expense-report":
				// State: init
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					_ = reqBody.CallbackQuery.Message.Chat.ID
					_ = reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					if data == "done" {
						logHandleInput = convExpenseSummary.Closure()
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init("date=", refDate)
					}
				}
			case "calendar":
				// State: calendar
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseSummary.DailyExpenseReport(expenseList, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init("date=", refDate)
					}
				}
			}
		}

		slog.Info("user_context", "user.conversation", user.Conversation, "user.state", user.State, "log_handle_input", logHandleInput)
		slog.Debug("reqBody", "update_obj", reqBody)
	})

	http.ListenAndServe(":8080", r)
}
