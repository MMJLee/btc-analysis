package database

import (
	"context"
	"log"
	"os"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type DBPool struct {
	*pgxpool.Pool
}

func NewPool() DBPool {
	dbConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_CONNECTION_STRING"))
	if err != nil {
		log.Panic(err)
	}
	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Panic(err)
	}
	return DBPool{pool}
}

type DBConn struct {
	*pgx.Conn
}

func NewConn() DBConn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_CONNECTION_STRING"))
	if err != nil {
		log.Panic(err)
	}
	pgxdecimal.Register(conn.TypeMap())
	return DBConn{conn}
}

type RedisClient struct {
	*redis.Client
}

func NewRedis() RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "",
		DB:       0,
	})
	return RedisClient{rdb}
}
