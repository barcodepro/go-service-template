package app

import (
	"fmt"
	"net"
)

type Config struct {
	ListenAddress  net.TCPAddr // Network address and port where the application should listen on
	AllowedOrigins string      // CORS policy allowed origins
	PostgresURL    string      // URL for connecting to Postgres service
}

// Validate checks configuration for stupid values
func (c *Config) Validate() error {
	if c.PostgresURL == "" {
		return fmt.Errorf("POSTGRES_URL should not be empty")
	}

	return nil
}
