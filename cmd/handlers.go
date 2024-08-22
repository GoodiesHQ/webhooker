package main

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func handleGet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	targets := webhookTable.Get(name)

	// log this information for testing/debug purposes.
	var evt *zerolog.Event

	if targets != nil {
		evt = log.Info().Int("count", len(targets))
	} else {
		evt = log.Warn()
	}

	evt.Str("name", name).Bool("found", targets != nil).Msg("GET")

	// for security purposes, just treat every GET as a 404, but log it
	respond(w, http.StatusNotFound, nil)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	targets := webhookTable.Get(name)

	if targets == nil {
		// name is not found in the configuration
		log.Warn().Str("name", name).Msg("webhook name not found")

		respond(w, http.StatusNotFound, nil)
		return
	}

	// read the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// failed to read the body
		log.Error().Str("name", name).Err(err).Msg("failed to read incoming request body")

		respond(w, http.StatusInternalServerError, nil)
		return
	}

	log.Debug().Str("body", string(body)).Msg("received body")

	// perform the webhooks to the backend targets in the background
	for _, target := range targets {
		go sendWebhook(name, target, r.Header, body)
	}

	// Ok!
	respond(w, http.StatusOK, nil)
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusMethodNotAllowed, nil)
}

func sendWebhook(name, targetUrl string, headers http.Header, body []byte) {
	var log = log.With().Str("webhook", name).Logger()

	target, err := url.Parse(targetUrl)
	if err != nil {
		log.Error().Err(err).Msg("invalid URL: " + targetUrl)
		return
	}

	// include the backend webhook hostname
	if log.Debug().Enabled() {
		log = log.With().Str("target", target.String()).Logger()
	} else {
		log = log.With().Str("target", target.Hostname()).Logger()
	}

	req, err := http.NewRequest(http.MethodPost, targetUrl, bytes.NewReader(body))
	if err != nil {
		log.Error().Err(err).Msg("failed to create new request")
		return
	}

	req.Header = headers.Clone()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to perform client request")
		return
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response")
		return
	}

	log.Info().Str("status", res.Status).Str("body", string(body)).Int("response_size", len(body)).Msg("target responded")
}

func respond(w http.ResponseWriter, statusCode int, data []byte) {
	w.WriteHeader(statusCode)
	if data == nil {
		w.Write([]byte(http.StatusText(statusCode)))
	} else {
		w.Write(data)
	}
}
