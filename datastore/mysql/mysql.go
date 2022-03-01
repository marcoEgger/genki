package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type MySQL struct {
	db   *sqlx.DB
	dsn  string
	opts *Options
}

const DriverName = "mysql"

// New will connect to the MySQL server using the given DSN
func New(dsn string, options ...Option) (*MySQL, error) {
	opts := newOptions(options...)

	db, err := sqlx.Connect(DriverName, dsn)
	if err != nil {
		return nil, err
	}

	// configure connection mysql pool
	db.SetMaxOpenConns(opts.MaxOpenConnections)
	db.SetMaxIdleConns(opts.MaxIdleConnections)
	db.SetConnMaxLifetime(opts.MaxConnectionLifetime)

	return &MySQL{
		db:   db,
		dsn:  dsn,
		opts: opts,
	}, nil
}

// Migrate to a specific version. It's assumed t
func (m MySQL) Migrate(version uint) error {
	db, err := sql.Open(DriverName, m.dsn)
	if err != nil {
		return errors.Wrap(err, "unable to open database connection")
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	migrations, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", m.opts.MigrationPath),
		DriverName,
		driver)
	if err != nil {
		return errors.Wrap(err, "unable initialize migrations")
	}

	err = migrations.Migrate(version)
	if err != nil {
		if strings.Contains(err.Error(), "no change") {
			return nil
		}
		return errors.Wrap(err, "failed to apply migrations")
	}

	return nil
}

// Close is just a proxy for convenient access to db.Close()
func (m MySQL) Close() error {
	return m.db.Close()
}

// DB is just a proxy for convenient access to the underlying sqlx implementation
// This method is used a lot, therefore it's name is abbreviated.
func (m MySQL) DB() *sqlx.DB {
	return m.db
}

// Options returns the currently set options.
func (m MySQL) Options() *Options {
	return m.opts
}
