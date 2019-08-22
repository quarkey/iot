# iot

iot stuff

requirements:
a running postgres database 9.4+
golang

# create user

CREATE USER iot WITH PASSWORD 'iot';
CREATE DATABASE iot WITH OWNER iot;

# building init + api

1. go build main.go
2. go build ./cmd/init/main.go
3. run init script ./cmd/init/init.exe to setup the tables required (double check flags)
4. run ./iot.exe to start api

# arduino sketch

1. write sketch to arduino
   post: {"sensor_id": 1, "dataset_id": "1,", "data": [123.00, 12.00]}
