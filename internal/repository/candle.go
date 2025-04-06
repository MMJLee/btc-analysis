package repository

import (
	"log"

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
		LIMIT $4 OFFSET $5 
	`
	rows, _ := repo.Pool.Query(repo.Context, query, ticker, start, end, limit, offset)
	candles, err := pgx.CollectRows(rows, pgx.RowToStructByName[util.Candle])
	if err != nil {
		log.Panicf("Error: Repository-Candle-GetCandles: %v", err)
	}
	return candles, nil
}

func (repo CandleConn) CopyCandles(ticker string, candles util.CandleSlice) error {
	_, err := repo.Conn.CopyFrom(
		repo.Context,
		pgx.Identifier{"candle_one_minute"},
		[]string{"ticker", "start", "open", "high", "low", "close", "volume"},
		&util.CandleSliceWithTicker{Ticker: ticker, CandleSlice: candles},
	)
	return err
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
		log.Fatal(err)
	}

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
