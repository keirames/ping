package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var L zerolog.Logger

func New() {
	L =
		zerolog.
			New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			With().
			Timestamp().
			Caller().
			Logger()
}

func ExecQueryError(err error) {
	L.Err(err).Msg("Exec query error")
}
