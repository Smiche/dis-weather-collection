package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func Init_db_conn(conf Config) *pgx.Conn {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return conn
}

func Close(conn *pgx.Conn) {
	conn.Close(context.Background())
}
