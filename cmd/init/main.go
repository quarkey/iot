package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sqlx.Open("sqlite3", "../../database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	query := []string{
		`create table sensordata (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			sensor_description text not null,
			serial text not null,
			temp text not null,
			time text not null
		);`,
		// `insert into data('sensor_description','serial','temp','time') values('sensor descr','serial no','23.12','2018-01-01 123301');`,
	}
	for _, q := range query {
		// fmt.Println("query:", q)
		_, err := db.Exec(q)
		if err != nil {
			fmt.Println(err)
		}
	}
}
