package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ph20Eoow/auth-svc/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const port = "80"

var dbRetryCount int

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	conn := dbConnect()
	if conn == nil {
		log.Panic("Can't connect to Postgres")
	}
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func dbConnect() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			dbRetryCount++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if dbRetryCount > 20 {
			log.Println(err)
			return nil
		}

		log.Println("dbConnect() retry in 2 second...")
		time.Sleep(2 * time.Second)
		continue
	}
}
