package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"auth-service/internal/config"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(cfg *config.DBConfig) *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.DBName))
	if err != nil {
		slog.Error("Проблема з підключенням до Postgres:", "error", err.Error())
		db.Close()
		return nil
	}

	pingCounts := 0
	for {
		err := db.Ping()
		if err == nil {
			slog.Info("Успішне підключення до Postgres")
			break
		}

		slog.Warn("Невдале підключення до Postgres", "count", pingCounts, "error", err.Error())
		pingCounts++

		if pingCounts > 3 {
			slog.Error("Невдале підключення до Postgres після 3 спроб:", "error", err.Error())
			os.Exit(1)
		}

		slog.Info("Очікування 5 секунд...")
		time.Sleep(5 * time.Second)
	}

	return db
}
