package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/revirathya/revirathya-telbot/infra/telegram"
	"github.com/revirathya/revirathya-telbot/internal/config"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/handler"
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
		log.Fatal(err)
	}

	// Set User context
	user := domain.NewUserContext()
	expenseList := make([]domain.Expense, 0)

	// Data: Categories
	var CATEGORIES []domain.Category = []domain.Category{
		{ID: 1, Label: "Food - Daily"},
		{ID: 2, Label: "Food - Snacks & Coffee"},
		{ID: 3, Label: "Groceries"},
	}

	// Conversation: Handler
	convExpense := handler.NewConversationExpense(tbh, user)
	convExpenseRecord := handler.NewConversationExpenseRecord(tbh, user)
	convExpenseList := handler.NewConversationExpenseList(tbh, user)

	// API handler
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
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

					if data == "expense-record" {
						logHandleInput = convExpenseRecord.Entry(chatID, messageID)
					} else if data == "expense-list" {
						logHandleInput = convExpenseList.Entry(chatID, messageID)
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
					_ = reqBody.CallbackQuery.Message.Chat.ID
					_ = reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseRecord.InteractionSelectCategory(CATEGORIES, data)
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
						logHandleInput = convExpenseList.Closure(expenseList, data)
					}
				}
			}
		}

		slog.Info("user_context", "user.conversation", user.Conversation, "user.state", user.State, "log_handle_input", logHandleInput)
		slog.Debug("reqBody", "update_obj", reqBody)
	})

	http.ListenAndServe(":8080", r)
}
