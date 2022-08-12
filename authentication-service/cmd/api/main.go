package main

import "database/sql"

const webPort = "70"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

}
