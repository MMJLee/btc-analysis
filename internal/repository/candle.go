package repository

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mmjlee/btc-analysis/internal/util"
)

func (repo CandlePool) GetCandles(ticker, start, end, limit, offset string) (util.CandleSlice, error) {
	query := `
		SELECT ticker, "start", "open", high, low, "close", volume
		FROM candle_one_minute 
		WHERE ticker = $1 
		AND "start" BETWEEN $2 AND $3
		ORDER BY ticker, "start"
		LIMIT $4 OFFSET $5;
	`
	rows, _ := repo.Pool.Query(repo.Context, query, ticker, start, end, limit, offset)
	candles, err := pgx.CollectRows(rows, pgx.RowToStructByName[util.Candle])
	if err != nil {
		return candles, util.WrappedError{Err: err, Message: "Repository-GetCandles-CollectRows"}
	}
	return candles, nil
}

func (repo CandlePool) GetMissingCandles(ticker, start, end, limit, offset string) (util.CandleSlice, error) {
	query := `
		WITH tmp AS (
			SELECT generate_series($2, $3, 60) AS "start"
			LIMIT $4 OFFSET $5
		) SELECT t.* FROM tmp t 
		LEFT JOIN candle_one_minute com 
		ON com.ticker = $1
		AND t."start" = com."start"
		WHERE com."start" IS NULL
	`
	rows, _ := repo.Pool.Query(repo.Context, query, ticker, start, end, limit, offset)
	candles, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[util.Candle])
	if err != nil {
		return candles, util.WrappedError{Err: err, Message: "Repository-GetMissingCandles-CollectRows"}
	}
	return candles, nil
}

func (repo CandleConn) CopyCandles(table_name string, ticker string, candles util.CandleSlice) error {
	_, err := repo.Conn.CopyFrom(
		repo.Context,
		pgx.Identifier{table_name},
		[]string{"ticker", "start", "open", "high", "low", "close", "volume"},
		&util.CandleSliceWithTicker{Ticker: ticker, CandleSlice: candles},
	)
	return err
}

func (repo CandleConn) CreateTable(table_name string) error {
	query := fmt.Sprintf(`CREATE TABLE %s (LIKE candle_one_minute INCLUDING DEFAULTS);`, table_name)
	_, err := repo.Conn.Exec(repo.Context, query)
	return err
}

func (repo CandleConn) DropTable(table_name string) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, table_name)
	_, err := repo.Conn.Exec(repo.Context, query)
	return err
}

func (repo CandleConn) InsertFromStaging(table_name string) error {
	query := fmt.Sprintf(` 
		INSERT INTO candle_one_minute (ticker, "start", "open", high, low, "close", volume) 
		SELECT ticker, "start", "open", high, low, "close", volume
		FROM %s
		ON CONFLICT ON CONSTRAINT candle_one_minute_pk DO UPDATE SET
			high = EXCLUDED.high,
			low = EXCLUDED.low,
			"close" = EXCLUDED.close,
			volume = EXCLUDED.volume
	`, table_name)
	_, err := repo.Conn.Exec(repo.Context, query)
	return err
}

func (repo CandleConn) BulkLogCandles(ticker string, candles util.CandleSlice) error {
	table_name := "staging_candle_one_minute"
	if err := repo.DropTable(table_name); err != nil {
		return err
	}
	if err := repo.CreateTable(table_name); err != nil {
		return err
	}
	if err := repo.CopyCandles(table_name, ticker, candles); err != nil {
		return err
	}
	if err := repo.InsertFromStaging(table_name); err != nil {
		return err
	}
	if err := repo.DropTable(table_name); err != nil {
		return err
	}
	return nil
}

func (repo CandleConn) InsertCandles(ticker string, candles util.CandleSlice) error {
	query := `
		INSERT INTO candle_one_minute (ticker, "start", "open", high, low, "close", volume) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT ON CONSTRAINT candle_one_minute_pk DO UPDATE SET
			high = EXCLUDED.high,
			low = EXCLUDED.low,
			"close" = EXCLUDED.close,
			volume = EXCLUDED.volume
		RETURNING ticker, "start", "open", high, low, "close", volume
	`
	batch := &pgx.Batch{}
	for _, candle := range candles {
		batch.Queue(query, ticker, candle.Start, candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
	}
	err := repo.Conn.SendBatch(repo.Context, batch).Close()
	if err != nil {
		return util.WrappedError{Err: err, Message: "Repository-InsertCandles-SendBatch"}
	}

	return err
}
