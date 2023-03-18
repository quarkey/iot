# iot sensorboard

IoT Sensor Board is a versatile backend yet not finished, written in Golang that enables you to effortlessly store data
points collected by Arduino devices. It provides a wide range of features, including the ability to add new devices,
configure datasets, and more. The Angular frontend application offers a user-friendly interface that simplifies the
process of setting up sensors and datasets.

Aside from data storage, you can also monitor signals (via socket) written to a dataset using live plots. This allows
you to keep track of your data in real-time and make informed decisions accordingly.

To further customize your data collection, you can create datasets with specified columns, such as temperature (float),
based on your Arduino sketch capabilities. This dataset can serve as a reference for controller operations to turn on
fans when the temperature reaches a certain threshold. Additionally, the system enables you to create timed switches and
on/off switches.

IoT Sensor Board is ideally suited for monitoring greenhouses that can control a range of temperature sensors and
automatic watering systems. In the future, we plan to incorporate solar panels into the system for added sustainability.

Please note that the current version of the system is not intended for critical environments and is recommended solely
for private use.

<insert drawing here>

## building backend and frontend for dev

The build process of the system is controlled by make commands and can be easily initiated by following the steps below.
Additional make features are available but not documented.

```
golang api

1. $ make build     // building go bins
2. $ make testdata  // adding testdata to db
3. $ make run       // start up go backend

or

1. go build
2. ./api --conf ./config/exampleconfig.json

angular frontend:
1. npm install
2. ng serve
```

Running the API locally on your machine is useful for testing and development. However, for a more permanent solution,
we recommend hosting everything via Docker containers. Please note that the API depends on a running Postgres server. To
learn more about how to set up Postgres, please refer to the instructions below.

### exampleconfig.json

The api require the following settings to be able to run.

```
{
    "connectString": "host=localhost port=15432 user=iot password=iot dbname=iot sslmode=disable",
    "driver": "postgres",
    "api_addr": ":6001",

    // not implemented
    "encryptionkey": "enter encryption key",

    // full path to up/down sql-files
    "migration": "database/migrations",

    // the numbers of seconds between telemetry checks
    "checkTelemetryTimer": 10,

    // enables simulation mode for development, this should be deactivated when running application
    "allowSim": true
}
```

## Docker container setup

To run iot sensorboard, you can use the "build" script ./docker_install.local.sh. located in ./resources/local. This
script will start three containers that are required for the application to run.

### nginx

The Angular frontend is hosted using Nginx. It is important to note that additional configuration is required to host a
Single Page Application (SPA). However, this configuration is already taken care of in the Dockerfile. For more
information, please refer to the documentation available at
https://medium.com/@technicadil_001/deploy-an-angular-app-with-nginx-a79cc1a44b49.

```
1. $ docker build -t ngimg .
2. $ docker run --name frontend -p 8080:80 -d ngimg
```

### postgres database

You can either start up a docker container directly or via the traditional Dockerfile.

```
$ docker run --restart always --platform linux/amd64 -v /Users/slundin/iotsensorboard/pg_data:/var/lib/postgresql/data --name pg -p 5432:5432 -d -e POSTGRES_PASSWORD=password postgres:latest

$ docker run --restart always --platform linux/amd64 -v /Users/slundin/iotsensorboard/pg_data:/var/lib/postgresql/data --name pg -p 15432:5432 -d -e POSTGRES_PASSWORD=password pgtest

// build from Dockerfile
$ docker build -t pgserverimg
$ docker run --restart always --name pgtwo -p 15432:5432 -d -e POSTGRES_USER=iot -e POSTGRES_DB=iot -e POSTGRES_PASSWORD=iot pgserver

```

## Arduino setup to run the application

To capture data points, for example, temperature, you can create a simple circuit that utilizes a temperature sensor and
an Arduino board. The circuit can be set up to collect a datapoint every n seconds and transmit it to the server with a
JSON payload. The JSON payload can include information such as the value of the datapoint and a timestamp. This can be
achieved by using an HTTP POST request to send the JSON payload to the server.

See example payload below

```
{"sensor_id": 1, "dataset_id": "1,", "data": ["123.00", "12.00"]}
```

## Setup commands QA environment.

To configure the QA environment for the iot sensorboard, you need to follow the steps outlined in the
docker_install.qa.sh script. This script contains detailed instructions for configuring the Go backend, Angular
frontend, and postgres. Please refer to the script for more information.

```
part 1
$ git clone git@github.com:quarkey/iot.git
$ cd iot
$ cd client
$ ng build --configuration qa-m1mini
$ cp dist ../../resources/qa/ng/dist
$ make build

part 2 can be done by running docker_install.qa.sh
$ docker network create qa-network
$ docker build -t qa_iot_backend .
$ cd resources/ng
$ docker build -t qa_iot_frontend .
$ cd ../pg
$ docker build -t qa_iot_pg .
$ cd ..
$ docker run --restart always --name qa_iot_pg --net qa-network -p 15432:5432 -d -e POSTGRES_USER=iot -e POSTGRES_DB=iot -e POSTGRES_PASSWORD=iot qa_iot_pg
$ docker run --restart always --name qa_iot_frontend --net qa-network -p 8081:80 -d qa_iot_frontend
$ docker run --name qa_iot_backend --net qa-network -p 6001:6001 -d qa_iot_backend

done!

```

#### migrate postgres database from server a to server b

```
0. docker stop qa_iot_backend
1. dump rpi database and copy over to new server
   pg_dump -Fc -f iot.dump.db -h localhost iot
   scp iot.dump.db slundin@192.168.10.159:/Users/slundin/devel/iot/resources
   docker cp iot.dump.db qa_iot_pg:/tmp

2. restore database
   docker stop qa_iot_backend
   docker exec -it qa_iot_pg sh
   su - postgres
   cd /tmp
   dropdb -h localhost iot
   createdb -h localhost -T template0 iot
   pg_restore -d iot -h localhost iot.dump.db
   psql -U iot
   update iot.schema_migrations set dirty='f';
   exit
   exit

3. start up again
    docker start qa_iot_backend
```

#### sqls for troubleshooting

Useful sql-commands

```sql
CREATE USER iot WITH PASSWORD 'iot';
CREATE DATABASE iot WITH OWNER iot;
GRANT ALL PRIVILEGES ON DATABASE iot TO iot;

alter user iot with password 'iot';
```

## database auto migrate

If you start the backend with the "--automigrate" switch, it will automatically migrate the database. However, you need
to ensure that you have installed "go-migrate" before running the migration command from the command line.

```
$ brew install golang-migrate

$ export POSTGRESQL_URL="postgresql://iot:iot@localhost:15432/iot"
$ migrate -database ${POSTGRESQL_URL} -path database/migrations down
$ migrate -database ${POSTGRESQL_URL} -path database/migrations up
```

#### raspberry pi server setup

Commands used to set up raspberry pi

```bash
    $ systemctrl enable ssh
    $ systemctrl start ssh
    $ apt-get update
    $ apt-get upgrade
    $ curl -fsSL https://get.docker.com -o get-docker.sh
    $ chmod +x ./get-docker.sh
    $ ./get-docker.sh
    $ usermod -aG docker slundin
    $ usermod -aG docker slundin

    # install go bin files
```

## Future ideas

Below are some features that I have in mind for future development:

### Backend

- Enable communication between Arduino devices and the backend to switch relay states.
- Add Auth0 authentication and JWT tokens to verify Arduino signals.

### arduino

- Provide Arduino sketches and examples on how to set up a "test" rig that covers all functionality in the IoT
  Sensorboard application.

### Frontend

- Make website more usable on mobile
