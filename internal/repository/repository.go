package repository

import (
	"context"
	"log"
	"os"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Config() *pgxpool.Config {
	conn_string := os.Getenv("DATABASE_CONNECTION_STRING")
	dbConfig, err := pgxpool.ParseConfig(conn_string)
	if err != nil {
		log.Panicf("Error: Repository-Config: %v", err)
	}
	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}
	return dbConfig
}

type CandlePool struct {
	context.Context
	*pgxpool.Pool
}

func NewCandlePool(ctx context.Context) CandlePool {
	pool, err := pgxpool.NewWithConfig(ctx, Config())
	if err != nil {
		log.Panicf("Error: Repository-NewCandlePool: %v", err)
	}
	return CandlePool{Context: ctx, Pool: pool}
}

type CandleConn struct {
	context.Context
	*pgx.Conn
}

func NewCandleConn(ctx context.Context) CandleConn {
	conn_string := os.Getenv("DATABASE_CONNECTION_STRING")
	conn, err := pgx.Connect(ctx, conn_string)
	if err != nil {
		log.Panicf("Error: Repository-NewCandleConn: %v", err)
	}
	pgxdecimal.Register(conn.TypeMap())
	return CandleConn{Context: ctx, Conn: conn}
}
