package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB() *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("Gagal Bikin Koneksi Ke database", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Gagal Ping DB", err)
	}
	return pool
}
