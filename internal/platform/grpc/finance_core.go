package grpc

import (
	financecorev1 "github.com/revirathya/revirathya-telbot/gen/go/financecore/v1"
	"google.golang.org/grpc"
)

type FinanceCoreConn struct {
	Conn                *grpc.ClientConn
	Expense             financecorev1.ExpenseServiceClient
	Category            financecorev1.CategoryServiceClient
	UserTelegramContext financecorev1.UserTelegramContextServiceClient
}

func NewFinanceCoreConn(addr string) (*FinanceCoreConn, error) {
	conn, err := dial("finance-core", addr)
	if err != nil {
		return nil, err
	}

	return &FinanceCoreConn{
		Conn:                conn,
		Expense:             financecorev1.NewExpenseServiceClient(conn),
		Category:            financecorev1.NewCategoryServiceClient(conn),
		UserTelegramContext: financecorev1.NewUserTelegramContextServiceClient(conn),
	}, nil
}

func (f *FinanceCoreConn) Close() error {
	return f.Conn.Close()
}
