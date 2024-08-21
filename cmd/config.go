package main

import (
	"net/url"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// global webhooks table
var webhookTable WebhookTable

type Config struct {
	Debug  bool `yaml:"debug,omitempty"`
	Listen struct {
		Host string `yaml:"host"`
		Port uint16 `yaml:"port"`
		SSL  struct {
			KeyPath string `yaml:"key"`
			CrtPath string `yaml:"crt"`
		} `yaml:"ssl"`
	}
	Webhooks []struct {
		Name    string   `yaml:"name"`
		Targets []string `yaml:"targets"`
	} `yaml:"webhooks"`
}

func (config *Config) SSL() bool {
	return config.Listen.SSL.KeyPath != "" && config.Listen.SSL.CrtPath != ""
}

// load a configuration from a YAML file, this function will panic on failure
func LoadConfig(filename string) *Config {
	var config Config
	log.Info().Msg("loading config file")

	// read the target file
	yamlData, err := os.ReadFile(filename)
	if err != nil {
		log.Panic().Err(err).Msg("failed to read the configuration file")
	}

	// parse the yaml configuration
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		log.Panic().Err(err).Msg("failed to parse the configuration file")
	}

	for _, webhook := range config.Webhooks {
		// create a list of valid targets
		validTargets := []string{}

		for _, target := range webhook.Targets {
			// parse each target to ensure it is a valid URL
			_, err := url.Parse(target)
			if err != nil {
				log.Warn().Str("url", target).Msg("URL is invalid")
				continue
			}

			// only append valid URLs to the slice
			validTargets = append(validTargets, target)
		}

		// set the targets for this name in the webhook table
		if err := webhookTable.Set(webhook.Name, validTargets, false); err != nil {
			log.Panic().Err(err).Send()
		}
	}

	/* handle defaults and zero values */

	// if the listen host is empty, set default
	if config.Listen.Host == "" {
		config.Listen.Host = DEFAULT_LISTEN_HOST
	}

	// if the listen port is empty, set default based on ssl
	if config.Listen.Port == 0 {
		if config.SSL() {
			config.Listen.Port = DEFAULT_LISTEN_PORT_SSL
		} else {
			config.Listen.Port = DEFAULT_LISTEN_PORT
		}
	}

	log.Info().Msgf("Loaded %d webhook handlers", len(config.Webhooks))

	return &config
}
