package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mmjlee/btc-analysis/internal/util"
)

func (repo DBPool) GetCandles(c context.Context, ticker, start, end, limit, offset string, missing bool) (util.CandleSlice, error) {
	query := `
	SELECT ticker, "start", "open", high, low, "close", volume
	FROM candle_one_minute 
	WHERE ticker = $1 
	AND "start" BETWEEN $2 AND $3
	ORDER BY ticker, "start"
	LIMIT $4 OFFSET $5;
	`
	if missing {
		query = `
			WITH tmp AS (
				SELECT generate_series($2, $3, 60) AS "start"
			) SELECT t.* FROM tmp t 
			LEFT JOIN candle_one_minute com 
			ON com.ticker = $1
			AND t."start" = com."start"
			WHERE com."start" IS NULL
			LIMIT $4 OFFSET $5;
		`
	}
	rows, _ := repo.Pool.Query(c, query, ticker, start, end, limit, offset)
	candles, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[util.Candle])
	if err != nil {
		return candles, util.WrappedError{Err: err, Message: "Repository-GetCandles-CollectRows"}
	}
	return candles, nil
}

func (repo DBConn) CopyCandles(c context.Context, tableName string, ticker string, candles util.CandleSlice) error {
	_, err := repo.Conn.CopyFrom(
		c,
		pgx.Identifier{tableName},
		[]string{"ticker", "start", "open", "high", "low", "close", "volume"},
		&util.CandleSliceWithTicker{Ticker: ticker, CandleSlice: candles},
	)
	return util.WrappedError{Err: err, Message: "Repository-CopyCandles-CopyFrom"}
}

func (repo DBConn) CreateTable(c context.Context, tableName string) error {
	query := fmt.Sprintf(`CREATE TABLE %s (LIKE candle_one_minute INCLUDING DEFAULTS);`, tableName)
	_, err := repo.Conn.Exec(c, query)
	return util.WrappedError{Err: err, Message: "Repository-CreateTable-Exec"}
}

func (repo DBConn) DropTable(c context.Context, tableName string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, tableName)
	_, err := repo.Conn.Exec(c, query)
	return util.WrappedError{Err: err, Message: "Repository-DropTable-Exec"}
}

func (repo DBConn) InsertFromStaging(c context.Context, tableName string) error {
	query := fmt.Sprintf(` 
		INSERT INTO candle_one_minute (ticker, "start", "open", high, low, "close", volume) 
		SELECT ticker, "start", "open", high, low, "close", volume
		FROM %s
		ON CONFLICT ON CONSTRAINT candle_one_minute_pk DO UPDATE SET
			high = EXCLUDED.high,
			low = EXCLUDED.low,
			"close" = EXCLUDED.close,
			volume = EXCLUDED.volume;
	`, tableName)
	_, err := repo.Conn.Exec(c, query)
	return util.WrappedError{Err: err, Message: "Repository-InsertFromStaging-Exec"}
}

func (repo DBConn) BulkLogCandles(c context.Context, ticker string, candles util.CandleSlice) error {
	tableName := "staging_candle_one_minute"
	if err := repo.DropTable(c, tableName); err != nil {
		return util.WrappedError{Err: err, Message: "Repository-BulkLogCandles-DropTable"}
	}
	if err := repo.CreateTable(c, tableName); err != nil {
		return util.WrappedError{Err: err, Message: "Repository-BulkLogCandles-CreateTable"}
	}
	if err := repo.CopyCandles(c, tableName, ticker, candles); err != nil {
		return util.WrappedError{Err: err, Message: "Repository-BulkLogCandles-CopyCandles"}
	}
	if err := repo.InsertFromStaging(c, tableName); err != nil {
		return util.WrappedError{Err: err, Message: "Repository-BulkLogCandles-InsertFromStaging"}
	}
	if err := repo.DropTable(c, tableName); err != nil {
		return util.WrappedError{Err: err, Message: "Repository-InsertFromStaging-DropTable"}
	}
	return nil
}

func (repo DBConn) InsertCandles(c context.Context, ticker string, candles util.CandleSlice) error {
	query := `
		INSERT INTO candle_one_minute (ticker, "start", "open", high, low, "close", volume) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT ON CONSTRAINT candle_one_minute_pk DO UPDATE SET
			high = EXCLUDED.high,
			low = EXCLUDED.low,
			"close" = EXCLUDED.close,
			volume = EXCLUDED.volume
		RETURNING ticker, "start", "open", high, low, "close", volume;
	`
	batch := &pgx.Batch{}
	for _, candle := range candles {
		batch.Queue(query, ticker, candle.Start, candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
	}
	err := repo.Conn.SendBatch(c, batch).Close()
	if err != nil {
		return util.WrappedError{Err: err, Message: "Repository-InsertCandles-SendBatch"}
	}
	return nil
}
