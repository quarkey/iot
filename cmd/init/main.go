package main

import (
	"flag"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

func main() {
	testdata := flag.Bool("testdata", false, "some test data will be included when creating")
	drop := flag.Bool("drop", false, "dropping tables before creating schemas")

	flag.Parse()

	db, err := sqlx.Open("postgres", "host=localhost port=25432 user=iot password=iot dbname=iot sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	if *drop {
		drop := []string{
			`drop table sensordata`,
			`drop table dataset`,
			`drop table sensor`,
		}
		runCommands(drop, db)
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
	runCommands(query, db)
	if *testdata {
		testdata := []string{
			`insert into sensor(description, arduino_key) values('temp og hydro', 'arduino serial');`,
			`insert into dataset(sensor_id, description, reference, intervalsec, fields) values(1,'temperatur measurement, growhouse','reference x',1800,'["temp", "hydro"]');`,
			`insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('dataset', 'id')),'["23.13","59.32","ubro"]');`,
		}
		runCommands(testdata, db)
	}
}
func runCommands(commands []string, db *sqlx.DB) {
	for _, q := range commands {
		_, err := db.Exec(q)
		if err != nil {
			log.Println(err)
		}
	}
}
