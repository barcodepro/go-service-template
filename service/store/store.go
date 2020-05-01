package store

import (
	"context"
	"fmt"
	"github.com/barcodepro/go-service-template/service/internal/log"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

// Config defines configuration for Store
type Config struct {
	PostgresURL string // URL for connecting to Postgres service
}

// Store is the persistent storage for application data
type Store struct {
	Config *Config
	PgDB   *pgxpool.Pool
	//AccountRepository *AccountRepository
}

// NewStore creates new store
func NewStore(ctx context.Context, c *Config) (*Store, error) {
	var s = new(Store)
	var t = time.NewTicker(3 * time.Second)
	var giveup = time.NewTicker(600 * time.Second)

	for {
		pgdbStore, err := NewPostgresStore(c.PostgresURL)
		if err == nil {
			log.Debug("connection to postgres successful")
			s.PgDB = pgdbStore
			break
		}
		log.Warnln("failed connect to postgres: ", err)
		select {
		case <-t.C:
			continue
		case <-giveup.C:
			return nil, fmt.Errorf("give up connecting to postgres")
		case <-ctx.Done():
			log.Info("context interrupt received")
			return nil, fmt.Errorf("context interrupt received")
		}
	}

	s.Config = c
	//s.AccountRepository = NewAccountRepository(s)
	return s, nil
}

// NewPostgresStore creates new Postgres connection pool using specified DSN
func NewPostgresStore(postgresURL string) (*pgxpool.Pool, error) {
	pgConfig, err := pgxpool.ParseConfig(postgresURL)
	if err != nil {
		return nil, err
	}

	pgConfig.MaxConns = 10
	pgConfig.MaxConnIdleTime = 2 * time.Minute
	pgConfig.MaxConnLifetime = 1 * time.Hour
	pgConfig.ConnConfig.PreferSimpleProtocol = true

	return pgxpool.ConnectConfig(context.Background(), pgConfig)
}

// Close gracefully closes connections to store services
func (s *Store) Close() {
	// close postgres connection
	s.PgDB.Close()
}
