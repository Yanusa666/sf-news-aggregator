package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	LogLevel      string           `json:"log_level"`
	HttpServer    HttpServerConfig `json:"http_server"`
	RSS           []string         `json:"rss"`
	RequestPeriod uint64           `json:"request_period"`
	Postgres      PostgresConfig   `json:"postgres"`
}

func NewConfig() *Config {
	conf := new(Config)

	b, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, conf)
	if err != nil {
		log.Fatal(err)
	}

	return conf
}

type HttpServerConfig struct {
	ListenAddress string `json:"listen_address"`
}

type PostgresConfig struct {
	URI string `json:"URI"`
}
