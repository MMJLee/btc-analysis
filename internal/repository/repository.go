package repository

import (
	"context"
	"os"

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
	conn_string := os.Getenv("DATABASE_CONNECTION_STRING")
	dbConfig, err := pgxpool.ParseConfig(conn_string)
	if err != nil {
		return CandlePool{}, util.WrappedError{Err: err, Message: "Repository-NewCandlePool-ParseConfig"}
	}
	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}
	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return CandlePool{}, util.WrappedError{Err: err, Message: "Repository-NewCandlePool-NewWithConfig"}
	}
	return CandlePool{Context: ctx, Pool: pool}, nil
}

type CandleConn struct {
	context.Context
	*pgx.Conn
}

func NewCandleConn(ctx context.Context) (CandleConn, error) {
	conn_string := os.Getenv("DATABASE_CONNECTION_STRING")
	conn, err := pgx.Connect(ctx, conn_string)
	if err != nil {
		return CandleConn{}, util.WrappedError{Err: err, Message: "Repository-NewCandleConn-Connect"}
	}
	pgxdecimal.Register(conn.TypeMap())
	return CandleConn{Context: ctx, Conn: conn}, nil
}
