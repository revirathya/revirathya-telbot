package subflow

import (
	"fmt"
	"log/slog"
	"math"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/revirathya/revirathya-telbot/infra/telegram"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/utils"
)

type Calendar struct {
	Telbot   *telegram.TelegramBotHandler
	User     *domain.UserContext
	Timezone *time.Location
}

type CalendarConfig struct {
	TotalDays int
	Padding   int
	TotalRows int
}

func NewCalendar(tbh *telegram.TelegramBotHandler, user *domain.UserContext) (*Calendar, error) {
	tz, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, err
	}

	return &Calendar{
		Telbot:   tbh,
		User:     user,
		Timezone: tz,
	}, nil
}

// Private
func (c *Calendar) getCalendarConfig(refDate time.Time) CalendarConfig {
	cfg := CalendarConfig{}

	// Total Days
	cfg.TotalDays = time.Date(refDate.Year(), refDate.Month()+1, 0, 0, 0, 0, 0, c.Timezone).Day()

	// Padding
	weekday := int(time.Date(refDate.Year(), refDate.Month(), 1, 0, 0, 0, 0, c.Timezone).Weekday())
	cfg.Padding = (((weekday - 1) % 7) + 7) % 7 // Set Mon = 0 instead of Sun = 0. Mitigate negative modulo

	// Total Rows
	cfg.TotalRows = int(math.Ceil(float64(cfg.TotalDays+cfg.Padding) / 7))

	return cfg
}

func (c *Calendar) newCalendarTitle(refDate time.Time) [][]tgbotapi.InlineKeyboardButton {
	header := make([][]tgbotapi.InlineKeyboardButton, 0)

	header = append(header, []tgbotapi.InlineKeyboardButton{
		utils.NewEmptyInlineKeyboardButton(refDate.Format("Jan 2006")),
	})
	header = append(header, []tgbotapi.InlineKeyboardButton{
		utils.NewEmptyInlineKeyboardButton("Mon"),
		utils.NewEmptyInlineKeyboardButton("Tue"),
		utils.NewEmptyInlineKeyboardButton("Wed"),
		utils.NewEmptyInlineKeyboardButton("Thu"),
		utils.NewEmptyInlineKeyboardButton("Fri"),
		utils.NewEmptyInlineKeyboardButton("Sat"),
		utils.NewEmptyInlineKeyboardButton("Sun"),
	})

	return header
}

func (c *Calendar) newCalendarDayButton(callbackPrefix string, year int, month time.Month, day int) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(
		fmt.Sprintf("%02d", day),
		fmt.Sprintf("%s%d-%02d-%02d", callbackPrefix, year, int(month), day),
	)
}

func (c *Calendar) newCalendarBody(callbackPrefix string, refDate time.Time) [][]tgbotapi.InlineKeyboardButton {
	cfg := c.getCalendarConfig(refDate)

	body := make([][]tgbotapi.InlineKeyboardButton, 0)
	row := make([]tgbotapi.InlineKeyboardButton, 0, 7)
	dayButton := tgbotapi.InlineKeyboardButton{}

	iter := 1
	padIter := 0
	for range cfg.TotalRows {
		row = nil
		for range 7 {
			// Create Button
			if (padIter < cfg.Padding) || (iter > cfg.TotalDays) {
				// Handle Start and End Pading button
				dayButton = utils.NewEmptyInlineKeyboardButton(" ")
				padIter += 1
			} else {
				// Handle Day button
				dayButton = c.newCalendarDayButton(callbackPrefix, refDate.Year(), refDate.Month(), iter)
				iter += 1
			}

			row = append(row, dayButton)
		}
		body = append(body, row)
	}

	return body
}

func (c *Calendar) newCalendarNavigation(refDate time.Time) []tgbotapi.InlineKeyboardButton {
	prevMonth := refDate.AddDate(0, -1, 0)
	nextMonth := refDate.AddDate(0, 1, 1)
	nav := []tgbotapi.InlineKeyboardButton{
		// tgbotapi.NewInlineKeyboardButtonData("<", fmt.Sprintf("calendar=%s", prevMonth.Format("2006-01"))),
		NewInlineCalendarButton("<", &prevMonth),
		utils.NewEmptyInlineKeyboardButton(" "),
		NewInlineCalendarButton(">", &nextMonth),
	}
	return nav
}

// Public
func (c *Calendar) Init(callbackPrefix string, refDateStr string) string {
	// Set Text
	text := "Please select date below:"

	// Build Calendar
	refDate, _ := time.ParseInLocation("2006-01", refDateStr, c.Timezone)
	calHeader := c.newCalendarTitle(refDate)
	calBody := c.newCalendarBody(callbackPrefix, refDate)
	calNav := c.newCalendarNavigation(refDate)

	// Set Reply Markup
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	rows = append(rows, calHeader...)
	rows = append(rows, calBody...)
	rows = append(rows, calNav)

	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// Edit Message
	if _, err := c.Telbot.EditMessage(c.User.InteractionChatID, c.User.InteractionMessageID, text, &inlineKeyboardMarkup); err != nil {
		slog.Error("failed to edit telegram message")
	}

	// Set User Context
	c.User.ChangeState("calendar")
	return "callback:calendar"
}

// Utility
func NewInlineCalendarButton(label string, refDate *time.Time) tgbotapi.InlineKeyboardButton {
	var dt time.Time
	if refDate == nil {
		dt = time.Now()
	} else {
		dt = *refDate
	}

	return tgbotapi.NewInlineKeyboardButtonData(
		label,
		fmt.Sprintf("calendar=%s", dt.Format("2006-01")),
	)
}
