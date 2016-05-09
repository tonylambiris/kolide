package config

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/mephux/common"
)

// Duration type
type Duration struct {
	time.Duration
}

// UnmarshalText custom duration type
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// Config for database, session and other information
type Config struct {
	Session  *session
	Database *database
	Server   *server
}

type session struct {
	Key           string
	EncryptionKey string `toml:"encryption_key"`
	Size          int
	Network       string
	Address       string
	Password      string
}

type database struct {
	Address  string
	Username string
	Password string
	Database string
	SSL      string
	Crt      string
	Key      string
}

type server struct {
	QueryTimeout Duration `toml:"query_timeout"`
	Production   bool
	Debug        bool
	Address      string
	Crt          string
	Key          string
	EnrollSecret string `toml:"enroll_secret"`
}

// Default will return a default configuration
func Default(debug, production bool) *Config {
	return &Config{
		Session: &session{
			Size:     10,
			Network:  "tcp",
			Address:  ":6379",
			Password: "",
			Key:      string(common.RandomCreateBytes(50)),
		},
		Database: &database{},
		Server: &server{
			// QueryTimeout: duration("30s").Duration,
			Production:   production,
			Debug:        debug,
			Address:      "127.0.0.1:8000",
			Crt:          "",
			Key:          "",
			EnrollSecret: "kolidedev",
		},
	}
}

// Load a given configuration by path
func Load(configPath string) (*Config, error) {
	var config Config

	_, err := toml.DecodeFile(configPath, &config)

	if err != nil {
		return nil, fmt.Errorf("error loading configuration file - %s", err.Error())
	}

	return &config, nil
}
