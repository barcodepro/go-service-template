package store

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// TestStoreConfig returns config for test purposes
func TestStoreConfig(_ *testing.T) *Config {
	return &Config{
		PostgresURL: "host=postgres dbname=billing_db_test user=postgres sslmode=disable",
	}
}

// TestPostgresStore provides store for testing with teardown function which allow to cleanup store after tests
func TestStore(t *testing.T, config *Config) (*Store, func(...string)) {
	t.Helper()
	var s = &Store{}

	assert.NotEmpty(t, config.PostgresURL)
	var sc = &Config{
		PostgresURL: config.PostgresURL,
	}
	s.Config = sc

	pool, err := NewPostgresStore(s.Config.PostgresURL)
	assert.NoError(t, err)
	assert.NotNil(t, pool)
	s.PgDB = pool

	//s.AccountRepository = NewAccountRepository(s)

	return s, func(tables ...string) {
		if len(tables) > 0 {
			if _, err := s.PgDB.Exec(context.Background(), fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ","))); err != nil {
				t.Fatal(err)
			}
		}
		s.Close()
	}
}
