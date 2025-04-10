package repository

import (
	"context"
	"log"
	"os"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmjlee/btc-analysis/internal/util"
)

type DBPool struct {
	*pgxpool.Pool
}

func NewPool() DBPool {
	dbConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_CONNECTION_STRING"))
	if err != nil {
		log.Panic(util.WrappedError{Err: err, Message: "Repository-NewCandlePool-ParseConfig"}.Error())
	}
	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Panic(util.WrappedError{Err: err, Message: "Repository-NewCandlePool-NewWithConfig"}.Error())
	}
	return DBPool{pool}
}

type DBConn struct {
	*pgx.Conn
}

func NewConn() DBConn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_CONNECTION_STRING"))
	if err != nil {
		log.Panic(util.WrappedError{Err: err, Message: "Repository-NewCandleConn-Connect"}.Error())
	}
	pgxdecimal.Register(conn.TypeMap())
	return DBConn{conn}
}
