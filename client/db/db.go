package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type LocalStation struct {
	ID             int32
	Name           string
	Number         int32
	OrganizationId int32
	Type           string
	Latitude       float32
	Longitude      float32
	Altitude       float32
	City           string
}

type GlobalStation struct {
	LocalStation
	Country string
}

type MeasurementMinMax struct {
	ID        int32
	StationID int32
	Min       float64
	Max       float64
	Avg       float64
	Time      time.Time
	PhenType  string
	Unit      string
}

type MeasurementMinMaxGlobal struct {
	MeasurementMinMax
	Country string
}

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

func Get_stations(conn *pgx.Conn) ([]GlobalStation, error) {
	rows, err := conn.Query(context.Background(), "select * from station_all")
	if err != nil {
		fmt.Println(err)
	}

	stations, err := pgx.CollectRows(rows, pgx.RowToStructByPos[GlobalStation])
	return stations, err
}

func Get_local_stations(conn *pgx.Conn) ([]LocalStation, error) {
	rows, err := conn.Query(context.Background(), "select * from station_local")
	if err != nil {
		fmt.Println(err)
	}

	stations, err := pgx.CollectRows(rows, pgx.RowToStructByPos[LocalStation])
	return stations, err
}

func Query_local_data(conn *pgx.Conn, station int, startTime time.Time, endTime time.Time) []MeasurementMinMax {
	rows, err := conn.Query(context.Background(), "select * from meas_min_max_day_local where station_info=$1 and time >= $2 and time <= $3 order by time ASC", station, startTime, endTime)
	if err != nil {
		fmt.Println(err)
	}
	measurements, err := pgx.CollectRows(rows, pgx.RowToStructByPos[MeasurementMinMax])
	if err != nil {
		log.Fatal(err)
	}

	return measurements
}

func Query_global_data(conn *pgx.Conn, station int, startTime time.Time, endTime time.Time) []MeasurementMinMaxGlobal {
	rows, err := conn.Query(context.Background(), "select * from meas_min_max_day_all where station_info=$1 and time >= $2 and time <= $3 order by time ASC", station, startTime, endTime)
	if err != nil {
		fmt.Println(err)
	}

	measurements, err := pgx.CollectRows(rows, pgx.RowToStructByPos[MeasurementMinMaxGlobal])
	if err != nil {
		log.Fatal(err)
	}

	return measurements
}

func printRows(rows pgx.Rows) {
	for rows.Next() {
		columnValues, _ := rows.Values()
		for i, v := range columnValues {
			fmt.Printf("Type of value at %v=%T, value=%v | ", i, v, v)
		}
	}
}
