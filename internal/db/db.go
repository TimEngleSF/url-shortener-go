package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pgOnce sync.Once

type PGDataSource struct {
	Host   string
	Port   string
	DbName string
	User   string
	Pass   string
	SSL    string
}

type Postgres struct {
	DB  *pgxpool.Pool
	Dsn *PGDataSource
}

func (p *Postgres) OpenDb() error {
	var err error
	var db *pgxpool.Pool
	pgOnce.Do(func() {
		db, err = pgxpool.New(context.Background(), p.Dsn.String())
	})
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}
	p.DB = db
	return nil
}

func ConvSSL(sslVar bool) string {
	sslStr := "disable"
	if sslVar {
		sslStr = "enable"
	}
	return sslStr
}

func (dsn PGDataSource) String() string {
	var dsnStr string
	dsnStr += fmt.Sprintf("postgres://%s", dsn.User)
	if dsn.Pass != "" {
		dsnStr += fmt.Sprintf(":%s", dsn.Pass)
	}
	dsnStr += fmt.Sprintf("@%s:%s/%s?sslmode=%s", dsn.Host, dsn.Port, dsn.DbName, dsn.SSL)
	return dsnStr
}
