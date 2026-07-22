package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	simpledb "github.com/auho/go-simple-db/v3"
	"github.com/auho/go-toolkit-mysql/internal/testutil"
	testmysql "github.com/auho/go-toolkit-testutil/mysql"
	gosqlmysql "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDatabase is the unified database name shared by all mysql tests.
const TestDatabase = "_test_mysql"

// LoadDSN returns a MySQL DSN pointing at TestDatabase from the TEST_MYSQL_DSN
// environment variable. The caller is responsible for loading .env.test
// (via testutil.LoadEnv) before calling this function.
func LoadDSN() (string, error) {
	dsn, err := testmysql.LoadDSN(TestDatabase)
	if err != nil {
		return "", fmt.Errorf("load dsn: %w", err)
	}

	if dsn == "" {
		return "", errors.New("TEST_MYSQL_DSN not set; create .env.test with TEST_MYSQL_DSN")
	}

	return dsn, nil
}

// EnsureDatabase creates the test database if it does not exist.
// It parses the given DSN to extract the database name, then connects
// using a DSN without the database name to avoid connection failures
// when the database is not yet created.
func EnsureDatabase(dsn string) error {
	cfg, err := gosqlmysql.ParseDSN(dsn)
	if err != nil {
		return fmt.Errorf("parse dsn: %w", err)
	}

	dbName := cfg.DBName
	cfg.DBName = ""
	rawDSN := cfg.FormatDSN()

	db, err := sql.Open("mysql", rawDSN)
	if err != nil {
		return fmt.Errorf("sql open: %w", err)
	}
	defer func() { _ = db.Close() }()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	if err != nil {
		return fmt.Errorf("create database %s: %w", dbName, err)
	}

	return nil
}

// Setup loads .env.test, builds the DSN, and ensures the test database exists.
// It returns the DSN ready for connecting.
func Setup() (string, error) {
	err := testutil.LoadEnv()
	if err != nil {
		return "", fmt.Errorf("load env: %w", err)
	}

	dsn, err := LoadDSN()
	if err != nil {
		return "", fmt.Errorf("load dsn: %w", err)
	}

	if err = EnsureDatabase(dsn); err != nil {
		return "", fmt.Errorf("ensure database: %w", err)
	}

	return dsn, nil
}

// NewDB creates a *simpledb.SimpleDB and *gorm.DB from the TEST_MYSQL_DSN
// environment variable. It ensures the test database exists, configures
// the gorm logger, and tunes the connection pool.
func NewDB() (*simpledb.SimpleDB, *gorm.DB, error) {
	dsn, err := Setup()
	if err != nil {
		return nil, nil, fmt.Errorf("setup: %w", err)
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Error,
			IgnoreRecordNotFoundError: true,
		},
	)

	simpleDB, gormDB, err := simpledb.NewMySQLGorm(dsn, &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, nil, fmt.Errorf("new mysql gorm: %w", err)
	}

	if sqlDB := simpleDB.SqlDB(); sqlDB != nil {
		conns := runtime.NumCPU() * 2
		sqlDB.SetMaxOpenConns(conns)
		sqlDB.SetMaxIdleConns(conns)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
	}

	return simpleDB, gormDB, nil
}

// RecreateTable drops and recreates the given table using the provided DDL.
func RecreateTable(db *gorm.DB, table string, ddl string) error {
	if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)).Error; err != nil {
		return fmt.Errorf("drop table %s: %w", table, err)
	}

	if err := db.Exec(ddl).Error; err != nil {
		return fmt.Errorf("create table %s: %w", table, err)
	}

	return nil
}
