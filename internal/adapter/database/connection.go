package database

import (
	"fmt"
	"time"

	"github.com/justjack1521/mevconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrFailedConnectToPostgres = func(err error) error {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
)

const (
	postgresMaxOpenConns    = 10
	postgresMaxIdleConns    = 2
	postgresConnMaxLifetime = 30 * time.Minute
	postgresConnMaxIdleTime = time.Minute
)

func NewPostgresConnection() (*gorm.DB, error) {
	config, err := mevconn.NewPostgresConfig()
	if err != nil {
		return nil, ErrFailedConnectToPostgres(err)
	}

	// pgx caches named prepared statements per connection by default. Behind the DigitalOcean
	// pool (PgBouncer, transaction mode) consecutive queries can run on different backends, so
	// a cached statement is missing ("prepared statement ... does not exist") or already there.
	// The simple protocol sends each query whole and prepares nothing.
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  config.Source(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, ErrFailedConnectToPostgres(err)
	}

	sql, err := db.DB()
	if err != nil {
		return nil, ErrFailedConnectToPostgres(err)
	}

	// Without a cap, a burst of requests opens one connection each until Postgres refuses
	// them, and those reads fail. Capped, a burst waits for a free connection instead.
	sql.SetMaxOpenConns(postgresMaxOpenConns)
	sql.SetMaxIdleConns(postgresMaxIdleConns)
	sql.SetConnMaxLifetime(postgresConnMaxLifetime)

	// The database's connection slots are shared by every service. An idle connection is
	// closed after a minute and few are kept, so a quiet service hands its slots back.
	sql.SetConnMaxIdleTime(postgresConnMaxIdleTime)

	return db, nil
}
