package domain

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
)

// Model: Expense
type Expense struct {
	Date          string
	CategoryID    int64
	CategoryLabel string
	Description   string
	Amount        int64
}

func (e *Expense) SetCategory(categoryID int64, categories []Category) {
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

	if e.CategoryID != 0 {
		tokens = append(tokens, fmt.Sprintf("<b>Category</b>: %s", e.CategoryLabel))
	}

	if e.Description != "" {
		tokens = append(tokens, fmt.Sprintf("<b>Description</b>: %s", e.Description))
	}

	if e.Amount != 0 {
		tokens = append(tokens, fmt.Sprintf("<b>Amount</b>: %s", humanize.Comma(e.Amount)))
	}

	return strings.Join(tokens, "\n")
}

func (e *Expense) FormatInline() string {
	return fmt.Sprintf("[%s] %s: Rp %s", e.CategoryLabel, e.Description, humanize.Comma(e.Amount))
}

func (e *Expense) Clear() {
	e.Date = ""
	e.CategoryID = 0
	e.CategoryLabel = ""
	e.Description = ""
	e.Amount = 0
}
