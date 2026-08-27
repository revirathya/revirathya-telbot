package container

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/utils"
)

type ExpenseContainer struct {
	TransactionDate InputDate
	CategoryID      InputString
	CategoryName    InputString
	Description     InputString
	Amount          InputInt
}

func NewFromDomain(exp *domain.Expense) ExpenseContainer {
	ec := ExpenseContainer{}
	ec.SetTransactionDate(exp.TransactionDate.Format("2006-01-02"))
	ec.SetCategory(exp.CategoryID.String(), exp.CategoryName)

	var desc string = ""
	if exp.Description != nil {
		desc = *exp.Description
	}
	ec.SetDescription(desc)

	ec.SetAmount(strconv.FormatInt(exp.Amount.IntPart(), 10))
	return ec
}

// Setter
func (ec *ExpenseContainer) SetTransactionDate(val string) error {
	dt, err := time.Parse("2006-01-02", val)
	if err != nil {
		return err
	}

	ec.TransactionDate.IsFill = true
	ec.TransactionDate.Value = dt
	return nil
}

func (ec *ExpenseContainer) SetCategory(valID string, valName string) {
	ec.CategoryID.IsFill = true
	ec.CategoryID.Value = valID

	ec.CategoryName.IsFill = true
	ec.CategoryName.Value = valName

}

func (ec *ExpenseContainer) SetDescription(val string) {
	ec.Description.IsFill = true
	ec.Description.Value = val
}

func (ec *ExpenseContainer) SetAmount(val string) error {
	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return err
	}

	ec.Amount.IsFill = true
	ec.Amount.Value = num
	return nil
}

// Clear
func (ec *ExpenseContainer) Clear() {
	ec.TransactionDate = InputDate{}
	ec.CategoryID = InputString{}
	ec.CategoryName = InputString{}
	ec.Description = InputString{}
	ec.Amount = InputInt{}
}

// Format
func (ec *ExpenseContainer) FormatDetail() string {
	rows := make([]string, 0)

	if ec.TransactionDate.IsFill {
		rows = append(rows, fmt.Sprintf("<b>Date</b>: %s", ec.TransactionDate.Value.Format("2006-01-02")))
	}

	if ec.CategoryID.IsFill {
		rows = append(rows, fmt.Sprintf("<b>Category</b>: %s", ec.CategoryName.Value))
	}

	if ec.Description.IsFill {
		descVal := ec.Description.Value
		if descVal == "" {
			descVal = "(No Description)"
		}
		rows = append(rows, fmt.Sprintf("<b>Description</b>: %s", descVal))
	}

	if ec.Amount.IsFill {
		amountFmt := utils.FormatThousand(ec.Amount.Value)
		rows = append(rows, fmt.Sprintf("<b>Amount</b>: Rp %s", amountFmt))
	}

	return strings.Join(rows, "\n")
}

func (ec *ExpenseContainer) FormatInline() string {
	return fmt.Sprintf("[%s] Rp %s", ec.CategoryName.Value, utils.FormatThousand(ec.Amount.Value))
}

// Export
func (ec *ExpenseContainer) ToDomain(userID uuid.UUID) (*domain.CreateExpense, error) {
	uuidCategoryID, err := uuid.Parse(ec.CategoryID.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid category id UUID format: %v", userID)
	}

	desc := &ec.Description.Value
	if *desc == "" {
		desc = nil
	}

	return &domain.CreateExpense{
		UserID:          userID,
		TransactionDate: ec.TransactionDate.Value,
		CategoryID:      uuidCategoryID,
		Description:     desc,
		Amount:          ec.Amount.Value,
	}, nil
}
