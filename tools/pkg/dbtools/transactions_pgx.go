package dbtools

import (
	"context"
	errors "github.com/commerce/tools/pkg/helpers"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// IQuery интерфейс для выполнения запросов к БД
type IQuery interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// TxRepository интерфейс для расширения репозиториев
type TxRepository interface {
	WithTx(ctx context.Context, handler func(ctx context.Context) error) error
	GetDb(ctx context.Context) IQuery
	GetNative() *pgxpool.Pool
}

type txKey string

const (
	TX_CTX_KEY txKey = "TX_CTX_KEY"
)

// txRepository структура для добавления возможности выполнения транзакций на уровень выше репозитория
type txRepository struct {
	db *pgxpool.Pool
}

// NewTxRepository обертка над репозиторием, для выполнения транзакций
func NewTxRepository(db *pgxpool.Pool) TxRepository {
	return &txRepository{
		db: db,
	}
}

// getTx возвращает транзакцию и флаг указывающий ее наличие
func (t *txRepository) getTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(TX_CTX_KEY).(pgx.Tx)
	return tx, ok
}

// WithTx создание транзакции и передача его в контекст
func (t *txRepository) WithTx(ctx context.Context, handler func(ctx context.Context) error) error {

	var err error
	tx, txExists := t.getTx(ctx)

	if !txExists {
		// создаем новую транзакцию, только если ее ещё нет в контексте
		// init TX
		tx, err = t.db.Begin(ctx)
		if err != nil {
			return errors.Wrap("init transaction error:", err)
		}
		defer func(ctx context.Context) { _ = tx.Rollback(ctx) }(ctx)

		// set TX to context
		ctx = context.WithValue(ctx, TX_CTX_KEY, tx)
	}

	if err := handler(ctx); err != nil {
		return errors.Wrap("exec transaction error:", err)
	}

	if !txExists {
		// завершаем только ту транзакцию, которую создали
		ctx = context.WithValue(ctx, TX_CTX_KEY, nil)

		// commit TX
		if err := tx.Commit(ctx); err != nil {
			return errors.Wrap("commit transaction error:", err)
		}
	}

	return nil
}

// GetDb возвращает сущность для выполнения запросов в БД
func (t *txRepository) GetDb(ctx context.Context) IQuery {
	tx, ok := ctx.Value(TX_CTX_KEY).(pgx.Tx)
	if !ok {
		return t.db
	}
	return tx
}

// GetNative возвращает указатель на пул соединений БД
func (t *txRepository) GetNative() *pgxpool.Pool {
	return t.db
}
