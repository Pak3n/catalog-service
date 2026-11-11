package rcpostgres

import (
	"context"

	"github.com/uptrace/bun"
)

type contextKeyTx struct{}

var txKey contextKeyTx

func getTxFromContext(ctx context.Context) bun.Tx {
	if tx, ok := ctx.Value(txKey).(bun.Tx); ok {
		return tx
	}
	return bun.Tx{}
}

func setTxToContext(ctx context.Context, tx bun.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}
