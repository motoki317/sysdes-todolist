package db

// conn.go provides helper functions for connection to DB
import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // initialize mysql driver
	"github.com/jmoiron/sqlx"
)

// DefaultDSN creates default DSN string
func DefaultDSN(host, port, user, password, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo", user, password, host, port, dbname)
}

// Connect opens connection to DB
func Connect(dsn string) (*sqlx.DB, error) {
	// Establish connection
	conn, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Check connection
	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
