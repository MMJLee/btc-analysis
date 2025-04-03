package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"mjlee.dev/btc-analysis/internal/util"
)

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx pgx.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

func CreateStagingTable() {

}

func DropStagingTable() {

}

func InsertFromStagingTable() {

}

func (q *Queries) CreateStagingTable(ctx context.Context, ticker string, arg util.Candle) {
	create_table_query := `
		CREATE TABLE IF NOT EXISTS staging_candle_one_minute (LIKE candle_one_minute INCLUDING ALL);
	`
	_, err := q.db.Exec(ctx, create_table_query)
	if err != nil {
		log.Panicf("Error: Repository-CreateStagingTable-Exec: %v", err)
	}
}
