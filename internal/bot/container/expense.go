package container

import (
	"fmt"
	"strings"
	"time"

	"github.com/revirathya/revirathya-telbot/internal/domain"
	"github.com/revirathya/revirathya-telbot/internal/utils"
	"github.com/shopspring/decimal"
)

type ExpenseContainer struct {
	TransactionDt InputDate
	CategoryID    InputString
	CategoryName  InputString
	Description   InputString
	Amount        InputInt
}

func NewFromDomain(exp *domain.Expense) ExpenseContainer {
	ec := ExpenseContainer{}
	ec.SetTransactionDt(exp.TransactionDt)
	ec.SetCategory(exp.CategoryID.String(), exp.CategoryName)
	ec.SetDescription(exp.Description)
	ec.SetAmount(exp.Amount)
	return ec
}

// Setter
func (ec *ExpenseContainer) SetTransactionDtFromString(val string) error {
	tm, err := time.Parse("2006-01-02", val)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}
	ec.SetTransactionDt(tm)
	return nil
}

func (ec *ExpenseContainer) SetTransactionDt(tm time.Time) {
	ec.TransactionDt.IsFill = true
	ec.TransactionDt.Value = tm
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

func (ec *ExpenseContainer) SetAmountFromString(val string) error {
	dc, err := decimal.NewFromString(val)
	if err != nil {
		return err
	}

	ec.SetAmount(dc)
	return nil
}

func (ec *ExpenseContainer) SetAmount(dc decimal.Decimal) {
	ec.Amount.IsFill = true
	ec.Amount.Value = dc.IntPart()
}

// Clear
func (ec *ExpenseContainer) Clear() {
	ec.TransactionDt = InputDate{}
	ec.CategoryID = InputString{}
	ec.CategoryName = InputString{}
	ec.Description = InputString{}
	ec.Amount = InputInt{}
}

// Format
func (ec *ExpenseContainer) FormatDetail() string {
	rows := make([]string, 0)

	if ec.TransactionDt.IsFill {
		rows = append(rows, fmt.Sprintf("<b>Date</b>: %s", ec.TransactionDt.Value.Format("2006-01-02")))
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
func (ec *ExpenseContainer) ToDomain() (domain.ExpenseForm, error) {
	categoryID, err := utils.StringToUUID("category_id", ec.CategoryID.Value)
	if err != nil {
		return domain.ExpenseForm{}, err
	}

	return domain.ExpenseForm{
		TransactionDt: ec.TransactionDt.Value,
		CategoryID:    categoryID,
		Description:   ec.Description.Value,
		Amount:        ec.Amount.Value,
	}, nil
}
