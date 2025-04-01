package repository

import (
	"context"

	"mjlee.dev/btc-analysis/internal/util"
)

func (q *Queries) Create(ctx context.Context, ticker string, arg util.Candle) (util.Candle, error) {
	create_candle_query := `
		INSERT INTO candle_one_minute (ticker, "start", "open", high, low, "close", volume)
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		ON CONFLICT ON CONSTRAINT candle_one_minute_pk DO UPDATE SET
			high = EXCLUDED.high,
			low = EXCLUDED.low,
			"close" = EXCLUDED.close,
			volume = EXCLUDED.volume
		RETURNING ticker, "start", "open", high, low, "close", volume
	`
	row := q.db.QueryRow(ctx, create_candle_query, ticker, arg.Start, arg.Open, arg.High, arg.Low, arg.Close, arg.Volume)
	var c util.Candle
	err := row.Scan(
		&c.Ticker,
		&c.Start,
		&c.Open,
		&c.High,
		&c.Low,
		&c.Close,
		&c.Volume,
	)
	return c, err
}

func (q *Queries) Read(ctx context.Context, ticker string, arg util.Candle) (util.Candle, error) {
	create_candle_query := `
		SELECT ticker, "start", "open", high, low, "close", volume
		FROM candle_one_minute 
		WHERE ticker = $1 AND "start" = $2
	`
	row := q.db.QueryRow(ctx, create_candle_query, ticker, arg.Start)
	var c util.Candle
	err := row.Scan(
		&c.Ticker,
		&c.Start,
		&c.Open,
		&c.High,
		&c.Low,
		&c.Close,
		&c.Volume,
	)
	return c, err
}

func (q *Queries) ReadList(ctx context.Context, ticker string, arg util.Candle) (util.Candle, error) {
	create_candle_query := `
		SELECT ticker, "start", "open", high, low, "close", volume
		FROM candle_one_minute 
		WHERE ticker = $1 AND "start" = $2
	`
	row := q.db.QueryRow(ctx, create_candle_query, ticker, arg.Start)
	var c util.Candle
	err := row.Scan(
		&c.Ticker,
		&c.Start,
		&c.Open,
		&c.High,
		&c.Low,
		&c.Close,
		&c.Volume,
	)
	return c, err
}
