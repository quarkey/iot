# iot sensorboard

Iot sensorboard is a simple golang server that can store datapoints captured by arduino devices. The system enables you to add new devices, configure datasets and much more.
<insert drawing here>

## building backend

Build process governed by make commands and can be intialized by the steps below.

```
1. $ make build
2. $ make testdata
3. $ make run
```

Additional make features are available, but not documented.

## arduino setup

Create a simple circuit that can capture datapoints, e.g temperature sensor. Datapoint will be collected every n seconds and sent to the server with a json payload such as:

```
{"sensor_id": 1, "dataset_id": "1,", "data": [123.00, 12.00]}
```

## docker setup

### nginx

Angular frontend is hosted by nginx, but please note that additional configuration is required to host SPA. Read https://medium.com/@technicadil_001/deploy-an-angular-app-with-nginx-a79cc1a44b49 for more information.

```
1. $ docker build -t iotng .
2. $ docker run --name iotbackend -p 8080:80 -d n
```

### postgres database

Not currently working, create a working docker file and testing is needed.

#### sqls

Useful sql-commands

```sql
CREATE USER iot WITH PASSWORD 'iot';
CREATE DATABASE iot WITH OWNER iot;
alter user iot with password 'iot';
GRANT ALL PRIVILEGES ON DATABASE iot TO iot;
```

## database automigrate

```
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
