package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/Pak3n/catalog-service/internal/app/config/section"
	"github.com/Pak3n/catalog-service/migration"
)

type Client struct {
	_bunDB   bun.IDB
	rawBunDB *bun.DB
	cfg      section.RepositoryPostgres
}

func (c *Client) GetRawBunDB() *bun.DB {
	return c.rawBunDB
}

func NewClient(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Host:   cfg.Address,
		Path:   "/" + cfg.Name,
	}
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("read_timeout", cfg.ReadTimeout.String()) // ← добавляем read_timeout
	q.Set("write_timeout", cfg.WriteTimeout.String())
	u.RawQuery = q.Encode()

	dsn := u.String()
	log.Printf("ReadTimeout: %s, WriteTimeout: %s", cfg.ReadTimeout, cfg.WriteTimeout)

	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	sqlDB := sql.OpenDB(connector)
	sqlDB.SetMaxOpenConns(10)

	bunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	ctxPing, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := bunDB.PingContext(ctxPing); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{
		_bunDB:   bunDB,
		rawBunDB: bunDB,
		cfg:      cfg,
	}, nil
}

func getMaxVersion(ctx context.Context, migrator *migrate.Migrator) (int64, error) {
	applied, err := migrator.AppliedMigrations(ctx)
	if err != nil {
		return 0, err
	}
	if len(applied) == 0 {
		return 0, nil
	}

	var maxVer int64
	for _, m := range applied {
		v, err := strconv.ParseInt(m.Name, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse migration version %s: %w", m.Name, err)
		}
		if v > maxVer {
			maxVer = v
		}
	}
	return maxVer, nil
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {
	migrations := migrate.NewMigrations()

	if err := migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}

	opts := []migrate.MigratorOption{
		migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable + "_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	}

	migrator := migrate.NewMigrator(c.rawBunDB, migrations, opts...)

	if err := migrator.Init(ctx); err != nil {
		return 0, 0, fmt.Errorf("failed to init migrator: %w", err)
	}

	oldVer, err = getMaxVersion(ctx, migrator)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get old version: %w", err)
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to apply migrations: %w", err)
	}

	newVer = oldVer
	for _, m := range group.Migrations {
		version, err := strconv.ParseInt(m.Name, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse migration version %s: %w", m.Name, err)
		}
		if version > newVer {
			newVer = version
		}
	}

	return oldVer, newVer, nil
}
