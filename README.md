# iot sensorboard

Iot sensorboard is a simple golang server that can store data points captured by Arduino devices. The system enables you to add new devices, configure datasets, and more. Setting up sensors and datasets can be done via an angular frontend.

<insert drawing here>

## building backend for dev

The build process is governed by make commands and can be initialized by the steps below.

```
1. $ make build     // building go bins
2. $ make testdata  // adding testdata to db
3. $ make run       // start up go backend
```

Additional make features are available but not documented.

## Arduino setup

Create a simple circuit that can capture data points, e.g., temperature sensor. Datapoint will be collected every n second and sent to the server with a JSON payload such as:

```
{"sensor_id": 1, "dataset_id": "1,", "data": ["123.00", "12.00"]}
```

## docker setup

### nginx

The angular frontend is hosted with Nginx, but please note that additional configuration is required to host SPA. However this is taken care of in the Dockerfile. Read https://medium.com/@technicadil_001/deploy-an-angular-app-with-nginx-a79cc1a44b49 for more information.

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

## Setup commands QA environment.

steps needed to configure qa environment, from go backend, ng and docker containers. Read docker_install.qa.sh for more details

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
## migrate postgres database from server a to server b
    
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

go-migrate must be installed before you can run migrate from the command line.

```
$ brew install golang-migrate

$ export POSTGRESQL_URL="postgresql://iot:iot@localhost:5432/iot"
$ migrate -database ${POSTGRESQL_URL} -path database/migrations down
$ migrate -database ${POSTGRESQL_URL} -path database/migrations up
```

## raspberry pi server setup

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
