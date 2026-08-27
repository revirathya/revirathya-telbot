package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
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
	grpcRepo "github.com/revirathya/revirathya-telbot/internal/repository/grpc"

	ucCtg "github.com/revirathya/revirathya-telbot/internal/usecase/category"
	ucExp "github.com/revirathya/revirathya-telbot/internal/usecase/expense"
	ucUtel "github.com/revirathya/revirathya-telbot/internal/usecase/usercontext"
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

	// Set gRPC Servers
	fincoreConn, err := grpcCon.NewFinanceCoreConn(cfg.FinanceCoreGRPC)
	if err != nil {
		slog.Error("failed to connect finance core grpc server", "error", err)
		os.Exit(1)
	}
	defer fincoreConn.Close()

	// Set Repository
	ctgRepo := grpcRepo.NewCategoryRepositoryImpl(fincoreConn)
	expRepo := grpcRepo.NewExpenseRepositoryImpl(fincoreConn)
	utelRepo := grpcRepo.NewUserTelegramContextRepositoryImpl(fincoreConn)

	// Set Usecase
	ctgUc := ucCtg.NewCategoryUsecase(ctgRepo)
	expUc := ucExp.NewExpenseUsecase(expRepo)
	utelUc := ucUtel.NewUserTelegramContextUsecase(utelRepo)

	// Conversation: Handler
	convExpense := handler.NewConversationExpense(tbh, utelUc)
	convExpenseRecord := handler.NewConversationExpenseRecord(tbh, utelUc, ctgUc, expUc)
	convExpenseList := handler.NewConversationExpenseList(tbh, utelUc, expUc)
	convExpenseSummary := handler.NewConversationExpenseSummary(tbh, utelUc, expUc)

	// Subflow: Calendar
	sfCalendar, err := subflow.NewCalendar(tbh, utelUc)
	if err != nil {
		slog.Error("failed when init calendar: %s", "error", err.Error())
		return
	}

	// Subflow: Numpad
	sfNumpad := subflow.NewNumpad(tbh)

	// Temp
	ctx := context.Background()

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

		// Get User
		var user *domain.UserTelegramContext
		var telegramChatID string
		if reqBody.Message != nil {
			telegramChatID = strconv.FormatInt(reqBody.Message.From.ID, 10)
		} else if reqBody.CallbackQuery != nil {
			telegramChatID = strconv.FormatInt(reqBody.CallbackQuery.From.ID, 10)
		}

		user, err = utelUc.GetUserContextByTelegramID(ctx, telegramChatID)
		if err != nil {
			slog.Error("failed to get user", "error", err)
		}

		// Cancel Command
		if (reqBody.Message != nil) && (reqBody.Message.Text == "/cancel") {
			handler.CommandCancel(tbh, user, ctx, utelUc, reqBody.Message.Chat.ID, "cancel")
		}

		// Handling Inputs (Update)
		switch *user.Conversation {
		case "root":
			// Conversation: Root
			if reqBody.Message != nil {
				// Handler: Command & Text
				chatID := reqBody.Message.Chat.ID
				text := reqBody.Message.Text

				if text[0] == '/' {
					// Handler: Command
					command := text[1:]
					switch command {
					case "expense":
						logHandleInput = convExpense.Entry(ctx, user, chatID)
					}
				}
			}
		case "expense":
			// Conversation: Expense
			switch *user.State {
			case "init":
				// State: init
				if reqBody.CallbackQuery != nil {
					// Handler: CallbackQuery
					chatID := reqBody.CallbackQuery.Message.Chat.ID
					messageID := reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					switch data {
					case "expense-record":
						logHandleInput = convExpenseRecord.Entry(ctx, user, chatID, messageID, data)
					case "expense-list":
						logHandleInput = convExpenseList.Entry(ctx, user, chatID, messageID, data)
					case "expense-summary":
						logHandleInput = convExpenseSummary.Entry(ctx, user, chatID, messageID, data)
					}
				}
			}
		case "expense-record":
			// Conversation: Expense Record
			switch *user.State {
			case "init":
				// State: init
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseRecord.InteractionSelectCategory(ctx, user, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init(ctx, user, "date=", refDate, data)
					}
				}
			case "calendar":
				// State: calendar
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseRecord.InteractionSelectCategory(ctx, user, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init(ctx, user, "date=", refDate, data)
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
						logHandleInput = convExpenseRecord.InteractionInputDescription(ctx, user, data)
					}
				}
			case "input-description":
				// State: input-description
				if reqBody.Message != nil {
					// Handler: Command & Text
					text := reqBody.Message.Text

					logHandleInput = convExpenseRecord.InteractionInputAmount(ctx, user, text)
				}
			case "input-amount":
				// State: input-amount
				if reqBody.CallbackQuery != nil {
					// Handler: CallbackQuery
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "amount=") {
						logHandleInput = convExpenseRecord.InteractionSelectConfirmation(ctx, user, data)
					} else if strings.HasPrefix(data, "numpad=") {
						currAmount := strings.Split(data, "=")[1]
						logHandleInput = sfNumpad.Init(ctx, user, "amount=", currAmount, data)
					}
				}
			case "select-confirmation":
				// State: select-confirmation
				if reqBody.CallbackQuery != nil {
					// Handler: CallbackQuery
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "confirmation=") {
						logHandleInput = convExpenseRecord.Closure(ctx, user, data)
					}
				}
			}
		case "expense-list":
			// Conversation: Expense List
			switch *user.State {
			case "init":
				// State: init
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					_ = reqBody.CallbackQuery.Message.Chat.ID
					_ = reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseList.SelectExpense(ctx, user, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init(ctx, user, "date=", refDate, data)
					}
				}
			case "calendar":
				// State: calendar
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseList.SelectExpense(ctx, user, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init(ctx, user, "date=", refDate, data)
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
						logHandleInput = convExpenseList.PreviewExpense(ctx, user, data)
					} else if data == "done" {
						logHandleInput = convExpenseList.Closure(ctx, user, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init(ctx, user, "date=", refDate, data)
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
						logHandleInput = convExpenseList.ConfirmationDeleteExpense(ctx, user, data)
					} else if data == "done" {
						logHandleInput = convExpenseList.Closure(ctx, user, data)
					} else if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseList.SelectExpense(ctx, user, data)
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
						logHandleInput = convExpenseList.AfterDeleteExpense(ctx, user, data)
					case "confirmation=no":
						logHandleInput = convExpenseList.PreviewExpense(ctx, user, data)
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
						logHandleInput = convExpenseList.Closure(ctx, user, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init(ctx, user, "date=", refDate, data)
					}
				}
			}
		case "expense-summary":
			// Conversation: Expense Summary
			switch *user.State {
			case "init":
				// State: init
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					_ = reqBody.CallbackQuery.Message.Chat.ID
					_ = reqBody.CallbackQuery.Message.MessageID
					data := reqBody.CallbackQuery.Data

					if data == "summary=daily-expense" {
						logHandleInput = convExpenseSummary.DailyExpenseFilterDate(ctx, user, data)
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
						logHandleInput = convExpenseSummary.DailyExpenseReport(ctx, user, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init(ctx, user, "date=", refDate, data)
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
						logHandleInput = convExpenseSummary.Closure(ctx, user, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init(ctx, user, "date=", refDate, data)
					}
				}
			case "calendar":
				// State: calendar
				if reqBody.CallbackQuery != nil {
					// Handler: Callback Query
					data := reqBody.CallbackQuery.Data

					if strings.HasPrefix(data, "date=") {
						logHandleInput = convExpenseSummary.DailyExpenseReport(ctx, user, data)
					} else if strings.HasPrefix(data, "calendar=") {
						refDate := strings.Split(data, "=")[1]
						logHandleInput = sfCalendar.Init(ctx, user, "date=", refDate, data)
					}
				}
			}
		}

		slog.Info("user_context", "user.conversation", user.Conversation, "user.state", user.State, "log_handle_input", logHandleInput)
		slog.Debug("reqBody", "update_obj", reqBody)
	})

	http.ListenAndServe(":8080", r)
}
