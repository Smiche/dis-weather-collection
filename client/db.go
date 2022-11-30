package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn

func init_db_conn(conf Config) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)
	var err error
	conn, err = pgx.Connect(context.Background(), url)
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

	fmt.Println(greeting)
}

func close_db_conn() {
	conn.Close(context.Background())
}
