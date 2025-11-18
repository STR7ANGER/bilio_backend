package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger

func New(level string) Logger {
	zl := zerolog.New(os.Stdout).With().Timestamp().Logger()
	lvl, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.TimeFieldFormat = time.RFC3339
	return zl.Level(lvl)
}
