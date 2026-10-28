package clients

import (
	"log/slog"

	financecorev1 "github.com/revirathya/revirathya-finance-core/gen/grpc/financecore/v1"
	"google.golang.org/grpc"

	// Internal
	appGrpc "github.com/revirathya/revirathya-telbot/internal/platform/grpc"
)

type FinanceClient struct {
	Conn                *grpc.ClientConn
	Expense             financecorev1.ExpenseServiceClient
	Category            financecorev1.CategoryServiceClient
	UserTelegramContext financecorev1.UserTelegramContextServiceClient
}

func NewFinanceClient(addr string, apiKey string) (*FinanceClient, error) {
	conn, err := appGrpc.Dial("finance-core", addr, apiKey)
	if err != nil {
		return nil, err
	}

	return &FinanceClient{
		Conn:                conn,
		Expense:             financecorev1.NewExpenseServiceClient(conn),
		Category:            financecorev1.NewCategoryServiceClient(conn),
		UserTelegramContext: financecorev1.NewUserTelegramContextServiceClient(conn),
	}, nil
}

func (f *FinanceClient) Close() {
	if err := f.Conn.Close(); err != nil {
		slog.Error("failed to close finance connection", "error", err)
	}
}
