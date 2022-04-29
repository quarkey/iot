#!/bin/bash
docker stop prod_iot_backend
docker stop prod_iot_frontend
docker stop prod_iot_pg

docker rm prod_iot_backend
docker rm prod_iot_frontend
docker rm prod_iot_pg

docker image rm prod_iot_backend
docker image rm prod_iot_frontend
docker image rm prod_iot_pg

docker network rm prod-network
