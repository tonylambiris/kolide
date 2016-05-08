package config

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/mephux/common"
)

type Duration struct {
	time.Duration
}

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
	Name     string
	Type     string
	Key      string
	Size     int
	Network  string
	Address  string
	Password string
}

type database struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
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

var (
	// DefaultSessionName to use for all session keys
	DefaultSessionName = "kolide_session"

	// DefaultCookieName to use for all cookies
	DefaultCookieName = "kolide_session"
)

// Default will return a default configuration
func Default(debug, production bool) *Config {
	return &Config{
		Session: &session{
			Name: DefaultCookieName,
			Type: "cookie",
			Key:  string(common.RandomCreateBytes(50)),
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
