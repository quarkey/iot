#!/bin/bash
docker stop local_iot_backend
docker stop local_iot_frontend
docker stop local_iot_pg

docker rm local_iot_backend
docker rm local_iot_frontend
docker rm local_iot_pg

docker image rm local_iot_backend
docker image rm local_iot_frontend
docker image rm local_iot_pg

docker network rm local-network
