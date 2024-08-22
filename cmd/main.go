package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	DEFAULT_LISTEN_HOST     = ""
	DEFAULT_LISTEN_PORT     = 80
	DEFAULT_LISTEN_PORT_SSL = 443
	DEFAULT_CONFIG_PATH     = "./webhooker.yml"
)

func main() {
	defer func() { recover() }()

	// set human readable log output
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

	// initialize the global webhook table
	webhookTable.Init()

	// get the configuration file path or use the default
	cfgPath := os.Getenv("WEBHOOKER_CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = DEFAULT_CONFIG_PATH
		log.Warn().Str("WEBHOOKER_CONFIG_PATH", cfgPath).Msg("using default config path")
	} else {
		log.Info().Str("WEBHOOKER_CONFIG_PATH", cfgPath).Msg("using provided config path")
	}

	// load the configuration file
	cfg := LoadConfig(cfgPath)

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// simple mux router
	mux := http.NewServeMux()

	// set the handlers
	mux.HandleFunc("GET /", handleGet)
	mux.HandleFunc("POST /", handlePost)
	mux.HandleFunc("/", handleDefault)

	var err error

	listen := fmt.Sprintf("%s:%d", cfg.Listen.Host, cfg.Listen.Port)
	log.Info().Str("service", listen).Msg("listening...")

	if cfg.SSL() {
		err = http.ListenAndServeTLS(listen, cfg.Listen.SSL.CrtPath, cfg.Listen.SSL.KeyPath, mux)
	} else {
		err = http.ListenAndServe(listen, mux)
	}

	log.Error().Err(err).Msg("finished")
}
