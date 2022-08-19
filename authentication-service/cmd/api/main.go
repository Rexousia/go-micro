package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Rexousia/go-micro/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	//connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}
	//set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}
	//setting up server and handler
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	//listen and serve on port 80
	//serving app.routes
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

//opening the DB
func openDB(dsn string) (*sql.DB, error) {
	//driver name //data source name (dsn)
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

//connecting to the opened db
func connectToDB() *sql.DB {
	//geting environment variables to access the db from the dsn
	dsn := os.Getenv("DSN")

	//looping over the connection until connection is made
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}
		//if not 10 counts log err and return out
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue

	}
}
