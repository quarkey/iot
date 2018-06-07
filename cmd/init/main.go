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
		`create table if not exists sensordata(
			id serial primary key,
			descr text not null,
			serial text not null,
			data jsonb not null
		  );`,
		`insert into sensordata(descr, serial, data) values('descr','serialno','{ "customer": "John Doe", "items": {"product": "Beer","qty": 6}}'::jsonb);`,
		`insert into sensordata(descr, serial, data) values('descr 2','serialno 2','{"abc": "val"}'::jsonb);`,
		`insert into sensordata(descr, serial, data) values('temperatur sensor', 'serialyo','{  "sensor description": "Temperature readings",  "hardware description": "arduino uno with temp sensor",  "serial": "a8f5f167f44f4964e6c998dee827110c",  "ip address": "192.168.10.100",  "network mask": "255.255.255.0",  "server": "192.168.10.1",  "encryption key": "8ed358a7da3cc760364909d4aaf7321e",  "record interval": "1800",  "data": {"serial": "a8f5f167f44f4964e6c998dee827110c","temp c": ["33.1","22.1"],"record time": ["113030","1200"]}}'::jsonb);`,
	}
	for _, q := range query {
		_, err := db.Exec(q)
		if err != nil {
			fmt.Println(err)
		}
	}
}
