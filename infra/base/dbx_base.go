package base

import (
	"context"
	"database/sql"
	"github.com/tietang/dbx"
)

const TX = "tx"

type BaseDao struct {
	TX *sql.Tx
}

func (d *BaseDao) SetTx(tx *sql.Tx) {
	d.TX = tx
}

type txFunc func(*dbx.TxRunner) error

// 事务执行
func Tx(fn txFunc) error {
	return TxContext(context.Background(), fn)
}

// 事务执行
func TxContext(ctx context.Context, fn txFunc) error {
	return DbxDatabase().Tx(fn)
}

// 将 runner 绑定到上下文
func WithValueContext(parent context.Context, runner *dbx.TxRunner) context.Context {
	return context.WithValue(parent, TX, runner)
}

func ExecuteContext(ctx context.Context, fn txFunc) error {
	tx := ctx.Value(TX).(*dbx.TxRunner)
	return fn(tx)
}
