#!/bin/bash

echo "[INFO] installing angular, and building dist for local docker environment"
cd ../../client
npm install --force
ng build

cp -r dist/* ../resources/local/ng/dist
cd ../resources/local

echo "[INFO] creating local docker network"
docker network create local-network

cd ../../
echo "[INFO] building local iot backend ..."
docker build -t local_iot_backend .

cd resources/local/ng
echo "[INFO] building nginx image local iot frontend ..."
docker build -t local_iot_frontend .

cd ../pg
echo "[INFO] building local iot postgres database ..."
docker build -t local_iot_pg .

# back to resources again
cd ../../

echo "[INFO] starting all containers ..."

docker run --restart always --name local_iot_pg --net local-network -p 25432:5432 -d -e POSTGRES_USER=iot -e POSTGRES_DB=iot -e POSTGRES_PASSWORD=iot -e POSTGRES_HOST=postgres local_iot_pg
docker run --restart always --name local_iot_frontend --net local-network -p 8081:80 -d local_iot_frontend
docker run --restart always --name local_iot_backend --net local-network -p 6002:6002 -d local_iot_backend
