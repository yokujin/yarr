package logging

import (
	"os"

	"github.com/nkanaev/yarr/src/utils/cli"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup() {
	if cli.LogFile != "" {
		file, err := os.OpenFile(cli.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal().Err(err).Msg("setup log file")
		}
		// defer file.Close()

		log.Logger = zerolog.New(file).With().Timestamp().Logger()

	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05.999"})
	}
}
