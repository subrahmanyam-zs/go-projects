package datastore

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/zopsmart/gofr/pkg"
)

func (c *SQLClient) Query(query string, args ...interface{}) (*sql.Rows, error) {
	begin := time.Now()
	rows, err := c.DB.Query(query, args...)

	c.monitorQuery(begin, query)

	return rows, err
}

func (c *SQLClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	begin := time.Now()
	rows, err := c.DB.Exec(query, args...)

	c.monitorQuery(begin, query)

	return rows, err
}

func (c *SQLClient) QueryRow(query string, args ...interface{}) *sql.Row {
	begin := time.Now()

	row := c.DB.QueryRow(query, args...)

	c.monitorQuery(begin, query)

	return row
}

func (c *SQLClient) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	begin := time.Now()

	rows, err := c.DB.QueryContext(ctx, query, args...)

	c.monitorQuery(begin, query)

	return rows, err
}

func (c *SQLClient) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	begin := time.Now()

	rows, err := c.DB.ExecContext(ctx, query, args...)

	c.monitorQuery(begin, query)

	return rows, err
}

func (c *SQLClient) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	begin := time.Now()
	row := c.DB.QueryRowContext(ctx, query, args...)

	c.monitorQuery(begin, query)

	return row
}

func checkQueryOperation(query string) string {
	query = strings.ToLower(query)
	query = strings.TrimSpace(query)

	if strings.HasPrefix(query, "select") {
		return "SELECT"
	} else if strings.HasPrefix(query, "update") {
		return "UPDATE"
	} else if strings.HasPrefix(query, "delete") {
		return "DELETE"
	} else if strings.HasPrefix(query, "commit") {
		return "COMMIT"
	} else if strings.HasPrefix(query, "begin") {
		return "BEGIN"
	} else if strings.HasPrefix(query, "set") {
		return "SET"
	}

	return "INSERT"
}

func (c *SQLClient) monitorQuery(begin time.Time, query string) {
	var (
		hostName string
		dbName   string
	)

	dur := time.Since(begin).Seconds()

	if c.config != nil {
		hostName = c.config.HostName
		dbName = c.config.Database
	}

	// push stats to prometheus
	sqlStats.WithLabelValues(checkQueryOperation(query), hostName, dbName).Observe(dur)

	ql := QueryLogger{
		Query:     []string{query},
		DataStore: pkg.SQL,
	}

	// log the query
	if c.logger != nil {
		ql.Duration = time.Since(begin).Microseconds()
		c.logger.Debug(ql)
	}
}

func (c *SQLClient) Begin() (*SQLTx, error) {
	begin := time.Now()

	tx, err := c.DB.Begin()
	c.monitorQuery(begin, "BEGIN")

	return &SQLTx{Tx: tx, logger: c.logger, config: c.config}, err
}

func (c *SQLClient) BeginTx(ctx context.Context, opts *sql.TxOptions) (*SQLTx, error) {
	begin := time.Now()

	tx, err := c.DB.BeginTx(ctx, opts)
	c.monitorQuery(begin, "BEGIN TRANSACTION")

	return &SQLTx{Tx: tx, logger: c.logger, config: c.config}, err
}

func (c *SQLTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	begin := time.Now()

	result, err := c.Tx.Exec(query, args...)
	c.monitorQuery(begin, query)

	return result, err
}

func (c *SQLTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	begin := time.Now()

	rows, err := c.Tx.Query(query, args...)
	c.monitorQuery(begin, query)

	return rows, err
}

func (c *SQLTx) QueryRow(query string, args ...interface{}) *sql.Row {
	begin := time.Now()

	row := c.Tx.QueryRow(query, args...)
	c.monitorQuery(begin, query)

	return row
}

func (c *SQLTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	begin := time.Now()

	result, err := c.Tx.ExecContext(ctx, query, args...)
	c.monitorQuery(begin, query)

	return result, err
}

func (c *SQLTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	begin := time.Now()

	rows, err := c.Tx.QueryContext(ctx, query, args...)
	c.monitorQuery(begin, query)

	return rows, err
}

func (c *SQLTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	begin := time.Now()

	row := c.Tx.QueryRowContext(ctx, query, args...)
	c.monitorQuery(begin, query)

	return row
}

func (c *SQLTx) Commit() error {
	begin := time.Now()

	err := c.Tx.Commit()
	c.monitorQuery(begin, "COMMIT")

	return err
}

func (c *SQLTx) monitorQuery(begin time.Time, query string) {
	var (
		hostName string
		dbName   string
	)

	dur := time.Since(begin).Seconds()

	if c.config != nil {
		hostName = c.config.HostName
		dbName = c.config.Database
	}

	var ql QueryLogger

	ql.Query = append(ql.Query, query)
	ql.Duration = time.Since(begin).Microseconds()
	ql.StartTime = begin
	ql.DataStore = pkg.SQL

	// push stats to prometheus
	sqlStats.WithLabelValues(checkQueryOperation(query), hostName, dbName).Observe(dur)

	// log the query
	if c.logger != nil {
		c.logger.Debug(ql)
	}
}
