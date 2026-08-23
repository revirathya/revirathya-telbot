package domain

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/revirathya/revirathya-telbot/internal/utils"
)

// Model: Expense
type Expense struct {
	ID            uuid.UUID
	Date          string
	CategoryID    string
	CategoryLabel string
	Description   string
	Amount        int64
}

func CreateExpense(date string, categoryID string, categoryLabel string, description string, amount int64) Expense {
	id := uuid.New()
	return Expense{
		ID:            id,
		Date:          date,
		CategoryID:    categoryID,
		CategoryLabel: categoryLabel,
		Description:   description,
		Amount:        amount,
	}
}

func (e *Expense) SetCategory(categoryID string, categories []Category) {
	for _, ctg := range categories {
		if categoryID == ctg.ID {
			e.CategoryID = categoryID
			e.CategoryLabel = ctg.Label
			break
		}
	}
}

func (e *Expense) FormatMessage() string {
	tokens := []string{}

	if e.Date != "" {
		tokens = append(tokens, fmt.Sprintf("<b>Date</b>: %s", e.Date))
	}

	if e.CategoryID != "" {
		tokens = append(tokens, fmt.Sprintf("<b>Category</b>: %s", e.CategoryLabel))
	}

	if e.Description != "" {
		tokens = append(tokens, fmt.Sprintf("<b>Description</b>: %s", e.Description))
	}

	if e.Amount != 0 {
		tokens = append(tokens, fmt.Sprintf("<b>Amount</b>: Rp %s", utils.FormatThousand(e.Amount)))
	}

	return strings.Join(tokens, "\n")
}

func (e *Expense) FormatInline() string {
	return fmt.Sprintf("[%s] Rp %s", e.CategoryLabel, utils.FormatThousand(e.Amount))
}

func (e *Expense) Clear() {
	e.Date = ""
	e.CategoryID = ""
	e.CategoryLabel = ""
	e.Description = ""
	e.Amount = 0
}
