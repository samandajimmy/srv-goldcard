package logger

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v9"
)

// DbLogger struct for db logging
type DbLogger struct{}

// BeforeQuery hook function before querying
func (d DbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

// AfterQuery hook function after querying
func (d DbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	sql, err := q.FormattedQuery()

	if err != nil {
		fmt.Println()
		fmt.Println("[*] SQL query format err:", err)
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("[*] SQL query:", sql)
	fmt.Println()

	return nil
}
