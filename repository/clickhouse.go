package repository

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type clickhouseRepository struct {
	conn    driver.Conn
	timeout time.Duration
}

func (r clickhouseRepository) Write(entrytime time.Time,
	level string,
	service string,
	userID string,
	message string) error {
	queryCtx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	query := `
    INSERT INTO app_logs (
        timestamp, level, service, user_id, message
    ) VALUES (?, ?, ?, ?, ?)
    `
	return r.conn.Exec(queryCtx, query, entrytime, level, service, userID, message)
}

func (r clickhouseRepository) Close() error {
	return r.conn.Close()
}
