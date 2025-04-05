package repository

import (
	"context"
	"log"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmjlee/btc-analysis/internal/util"
)

type CandlePool struct {
	context.Context
	*pgxpool.Pool
}

func NewCandlePool(ctx context.Context) (CandlePool, error) {
	pool, err := pgxpool.New(ctx, util.DATABASE_CONNECTION_STRING)
	if err != nil {
		return CandlePool{}, err
	}
	return CandlePool{Context: ctx, Pool: pool}, err
}

func (repo CandlePool) GetCandles(product_id string, start, end int64) (util.CandleSlice, error) {
	rows, _ := repo.Pool.Query(repo.Context, "")
	candles, err := pgx.CollectRows(rows, pgx.RowToStructByName[util.Candle])
	if err != nil {
		log.Panicf("Error: Repository-Candle-GetCandles: %v", err)
	}
	return candles, nil
}

type CandleConn struct {
	context.Context
	*pgx.Conn
}

func NewCandleConn(ctx context.Context) (CandleConn, error) {
	conn, err := pgx.Connect(ctx, util.DATABASE_CONNECTION_STRING)
	if err != nil {
		return CandleConn{}, err
	}
	pgxdecimal.Register(conn.TypeMap())
	return CandleConn{Context: ctx, Conn: conn}, nil
}

func (repo CandleConn) CopyCandles(product_id string, candles util.CandleSlice) error {
	_, err := repo.Conn.CopyFrom(
		repo.Context,
		pgx.Identifier{"candle_one_minute"},
		[]string{"ticker", "start", "open", "high", "low", "close", "volume"},
		&util.CandleSliceWithTicker{Ticker: product_id, CandleSlice: candles},
	)
	return err
}

// CreateUser creates a new user in the db..
// func (repo *CandlePool) CreateUser(candle util.Candle) (util.Candle, error) {
// 	test := [...]string{"h", "o", "h"}
// 	create_candle_query := `
// 		INSERT INTO candle_one_minute (ticker, "start", "open", high, low, "close", volume)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7)
// 		ON CONFLICT ON CONSTRAINT candle_one_minute_pk DO UPDATE SET
// 			high = EXCLUDED.high,
// 			low = EXCLUDED.low,
// 			"close" = EXCLUDED.close,
// 			volume = EXCLUDED.volume
// 		RETURNING ticker, "start", "open", high, low, "close", volume
// 	`
// 	_, err := repo.Pool.Exec(repo.Context, create_candle_query, test)

// 	result := repo.db.Create(&candle)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &domain.User{
// 		ID:        dbUser.ID,
// 		Username:  dbUser.Username,
// 		Password:  dbUser.Password,
// 		CreatedAt: dbUser.CreatedAt,
// 		UpdatedAt: dbUser.UpdatedAt,
// 	}, nil
// }
