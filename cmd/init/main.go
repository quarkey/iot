package main

import (
	"flag"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/quarkey/iot/models"
)

func main() {
	confPath := flag.String("conf", "", "path to your config")
	testdata := flag.Bool("testdata", false, "some test data will be included when creating")
	drop := flag.Bool("drop", false, "dropping tables before creating schemas")

	flag.Parse()

	if *confPath == "" {
		log.Fatalf("ERROR: missing configuration jsonfile")
	}
	srv := models.NewDB(*confPath)
	if *drop {
		drop := []query{
			{"dropping iot schema with cascade", "drop schema if exists iot cascade;"},
		}
		runCommandsDescr(drop, srv.DB)
	}
	querys := []query{
		{"creating schema iot", "create schema iot;"},
		{"creating sensor table",
			`create table sensor (
			id serial primary key,
			title text not null,
			description text not null, -- aurdino sensor description
			arduino_key text not null, -- unique identifyer key
			created_at timestamp NOT NULL DEFAULT now()::timestamp(0)
		  );`},
		{"creating dataset table",
			`create table dataset (
			id serial primary key,
			sensor_id integer references sensor (id),
			title text not null,
			description text not null,
			reference text not null,
			intervalsec int not null,
			fields jsonb,
			created_at timestamp NOT NULL DEFAULT now()::timestamp(0)
		  );`},
		{"creating sensordata table",
			`create table sensordata (
			id serial primary key,
			sensor_id integer references sensor (id),
			dataset_id integer references dataset (id),
			data jsonb,
			time timestamp NOT NULL DEFAULT now()::timestamp(0)
		  );`},
	}
	runCommandsDescr(querys, srv.DB)

	if *testdata {
		log.Println("Inserting testdata ...")
		// TODO: use struct instead of array
		testdata := []string{
			`insert into sensor(title, description, arduino_key) values('Arduino + Ethernet shield','Arduino UNO with Ethernet shield. LM35 temperatur sensor and hydrosensor. Used for project X', '8a1bbddba98a8d8512787d311352d951');`,
			`insert into dataset(sensor_id, title, description, reference, intervalsec, fields) values(1,'temp&hydro','Temperatur/hydro measurement, growhouse 1','8a1bbddba98a8d8512787d311352d951',1800,'["temp", "hydro"]');`,
			`insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('dataset', 'id')),'["23.13","59.32"]');`,
			`insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('dataset', 'id')),'["23.13","59.32"]');`,
			`insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('dataset', 'id')),'["23.13","59.32"]');`,
			`insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('dataset', 'id')),'["23.13","59.32"]');`,
			`insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('dataset', 'id')),'["23.13","59.32"]');`,

			`insert into sensor(title, description, arduino_key) values('Arduino + GPS','Arduino UNO with GPS tracking', '4987fb174ae91dc702394024378fc1cd');`,
			`insert into dataset(sensor_id, title, description, reference, intervalsec, fields) values(2,'Bicycle to work','Battery-driven lat/long tracker','4987fb174ae91dc702394024378fc1cd',1800,'["lat (n)", "long (e)", "direction"]');`,
			`insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('dataset', 'id')),'["58.8533","5.7329","e"]');`,
			`insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('dataset', 'id')),'["58.8533","5.7329","n/e"]');`,
			`insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('dataset', 'id')),'["58.8532","5.7329","n"]');`,
			`insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('dataset', 'id')),'["58.8531","5.7329","e"]');`,

			`insert into sensor(title, description, arduino_key) values('SuperduperRecordingbox','this device is awesome', 'dummy');`,
			// `insert into sensor(title, description, arduino_key) values('temp og hydro', 'a long description', dummy');`,
			// `insert into sensor(title, description, arduino_key) values('temp og hydro', 'a long description', dummy');`,
			// `insert into sensor(title, description, arduino_key) values('temp og hydro', 'a long description', dummy');`,
			// `insert into sensor(title, description, arduino_key) values('temp og hydro', 'a long description', dummy');`,
			// `insert into sensor(title, description, arduino_key) values('temp og hydro', 'a long description', dummy');`,
		}
		runCommands(testdata, srv.DB)
	}
}

type query struct {
	descr string
	query string
}

func runCommandsDescr(q []query, db *sqlx.DB) {
	for _, q := range q {
		log.Printf("%s\n", q.descr)
		_, err := db.Exec(q.query)
		if err != nil {
			log.Printf("DB ERROR: %v", err)
		}
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
