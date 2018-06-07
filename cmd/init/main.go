package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

func main() {
	db, err := sqlx.Open("postgres", "host=localhost port=25432 user=iot password=iot dbname=iot sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	query := []string{
		`create table sensor (
			id serial primary key,
			description text not null, -- aurdino sensor description
			arduino_key text not null, -- unique identifyer key
			created_at timestamp NOT NULL DEFAULT now()::timestamp(0)
		  );`,
		`create table dataset(
			id serial primary key,
			sensor_id integer references sensor (id),
			description text not null,
			reference text not null,
			intervalsec int not null,
			fields jsonb,
			created_at timestamp NOT NULL DEFAULT now()::timestamp(0)
		  );`,
		`create table sensordata (
			id serial primary key,
			sensor_id integer references sensor (id),
			dataset_id integer references dataset (id),
			data jsonb,
			time timestamp NOT NULL DEFAULT now()::timestamp(0)
		  );`,
	}
	for _, q := range query {
		_, err := db.Exec(q)
		if err != nil {
			fmt.Println(err)
		}
	}
}
