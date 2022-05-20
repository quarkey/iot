package models

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type query struct {
	descr string
	query string
}

// InsertTestdata insert data needed for tests to run
func (s *Server) InsertTestdata() error {
	log.Println("Inserting testdata ...")
	testdata := []query{
		{"adding arduino sensor 1", `insert into sensors(title, description, arduino_key) values('Arduino + Ethernet shield','Arduino UNO with Ethernet shield. LM35 temperatur sensor and hydrosensor. Used for project X', '8a1bbddba98a8d8512787d311352d951');`},
		{"adding dataset 1", `insert into datasets(sensor_id, title, description, reference, intervalsec, fields, types, showcharts) 
		values(1,'temp&hydro','Temperatur/hydro measurement, growhouse 1','8a1bbddba98a8d8512787d311352d951',4,'["temp", "hydro"]' ,'["float", "float"]', '["true","true"]');`},
		{"adding data point 1", `insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('datasets', 'id')),'["23.13","59.32"]');`},
		{"adding data point 2", `insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('datasets', 'id')),'["23.13","59.32"]');`},
		{"adding data point 3", `insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('datasets', 'id')),'["23.13","59.32"]');`},
		{"adding data point 4", `insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('datasets', 'id')),'["23.13","59.32"]');`},
		{"adding data point 5", `insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('datasets', 'id')),'["23.13","59.32"]');`},

		{"adding arduino sensor 2", `insert into sensors(title, description, arduino_key) values('Arduino + GPS','Arduino UNO with GPS tracking', '4987fb174ae91dc702394024378fc1cd');`},
		{"adding dataset 2", `insert into datasets(sensor_id, title, description, reference, intervalsec, fields, types, showcharts) 
		values(2,'Bicycle to work 1','Battery-driven lat/long tracker','4987fb174ae91dc702394024378fc1cd',4,'["lat (n)", "long (e)", "direction"]', '["float", "float", "string"]', '["true","true","false"]');`},
		{"data point 1", `insert into sensordata(sensor_id, dataset_id, data) values(2,currval(pg_get_serial_sequence('datasets', 'id')),'["58.8533","5.7329","e"]');`},
		{"data point 2", `insert into sensordata(sensor_id, dataset_id, data) values(2,currval(pg_get_serial_sequence('datasets', 'id')),'["58.8533","5.7329","n/e"]');`},
		{"data point 3", `insert into sensordata(sensor_id, dataset_id, data) values(2,currval(pg_get_serial_sequence('datasets', 'id')),'["58.8532","5.7329","n"]');`},
		{"data point 4", `insert into sensordata(sensor_id, dataset_id, data) values(2,currval(pg_get_serial_sequence('datasets', 'id')),'["58.8531","5.7329","e"]');`},

		{"adding arduino sensor 3", `insert into sensors(title, description, arduino_key) values('SuperduperRecordingbox','this device is awesome', 'dummykey');`},
		{"adding dataset 3", `insert into datasets(sensor_id, title, description, reference, intervalsec, fields, types) 
		values(3,'Bicycle to work 2','Battery-driven lat/long tracker','dummykey',4,'["lat (n)", "long (e)", "direction"]','["float","float","string"]');`},
		{"data point 1", `insert into sensordata(sensor_id, dataset_id, data) values(3,currval(pg_get_serial_sequence('datasets', 'id')),'["58.8533","5.7329","e"]');`},

		// controllers
		{"sensor with on off relay", `insert into sensors(title, description, arduino_key)
		values ('Arduino on off relay', 'roof fan', 'arduino_key123');
		`},
		{"controller threshold switch 1", `insert into controllers(sensor_id, category, title, description, items, alert, active) values
		(
			4,
			'thresholdswitch',
			'turn on fan above 25c',
			'when temperature is above 25c turn on roof fan',
			'{
				"description": "dataset col 0 > 25c", 
				"datasource": "d1c0", 
				"operation": 
				"greather than", 
				"threshold_limit": 25, 
				"on": true 
			}',
			't',
			't'
		);
		`},
	}
	runCommandsDescr(testdata, s.DB)

	return nil
}

func runCommandsDescr(q []query, db *sqlx.DB) {
	for _, q := range q {
		log.Printf("%s\n", q.descr)
		_, err := db.Exec(q.query)
		if err != nil {
			log.Printf("DB ERROR: %v (%s)", err, q.descr)
		}
	}
}
