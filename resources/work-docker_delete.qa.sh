#!/bin/bash
docker stop qa_iot_backend
docker stop qa_iot_frontend
docker stop qa_iot_pg

docker rm qa_iot_backend
docker rm qa_iot_frontend
docker rm qa_iot_pg

docker image rm qa_iot_backend
docker image rm qa_iot_frontend
docker image rm qa_iot_pg

docker network rm qa-network
