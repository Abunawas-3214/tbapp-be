package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	// Format waktu standar ISO8601 (nyaman untuk SIEM)
	zerolog.TimeFieldFormat = time.RFC3339

	// Log level default
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Mode Development: Output ke console dengan warna dan format cantik
	// Di Production: Sebaiknya hapus ConsoleWriter agar output JSON murni ke stdout
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05",
	})
}

// GetLogger returns a pointer to the global zerolog logger
func GetLogger() *zerolog.Logger {
	return &log.Logger
}
