package grpc

import "google.golang.org/grpc"

type FinanceCoreConn struct {
	Conn *grpc.ClientConn
}

func NewFinanceCoreConn(addr string) (*FinanceCoreConn, error) {
	conn, err := dial("finance-core", addr)
	if err != nil {
		return nil, err
	}

	return &FinanceCoreConn{
		Conn: conn,
	}, nil
}

func (f *FinanceCoreConn) Close() error {
	return f.Conn.Close()
}
