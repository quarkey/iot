## create user

Before you run the init script, make sure you got a working Postgresql user in your database. If not, the following SQL command can be used:

```sql
CREATE USER iot WITH PASSWORD 'iot';
CREATE DATABASE iot WITH OWNER iot;
```

## building

build the app and init script to get up and running.

1. $ go build ./cmd/api/main.go
2. $ go build ./cmd/init/main.go
3. run init to add testdata $ ./cmd/init/init
4. run $ ./cmd/api/api

## arduino

Create a simple circuit that capture points you want to store.
This can be achieved by sending the data point payload to the iot backend.

```
{"sensor_id": 1, "dataset_id": "1,", "data": [123.00, 12.00]}
```

## docker setup

You'll need to have docker and golang installed on your machine.

### build image and run app

1. $ docker build -t iot .
2. $ docker run -p 6001:6001 --name iotsrv -d iot

#### tag image

1.  $ docker tag iot:latest iot:v1.0 // create tag
2.  $ docker rmi iot:v1.0 // remove tag

#### gomigrate

export POSTGRESQL_URL="postgresql://iot:iot@localhost:5432/iot"
migrate -database ${POSTGRESQL_URL} -path database/migrations down
