package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func ConnectDB() *pgxpool.Pool {
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 2. Ambil URL database dari Neon
	connStr := os.Getenv("DATABASE_URL")

	// 3. Buat Connection Pool (pgxpool jauh lebih stabil untuk SaaS)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Gagal parse config database: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Gagal koneksi ke database Neon: %v", err)
	}

	return pool
}
