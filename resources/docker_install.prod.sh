#!/bin/bash
echo "creating prod docker network"
docker network create prod-network

cd ..
echo "building prod iot backend ..."
docker build -t prod_iot_backend .

cd resources/ng
echo "building prod iot frontend"
docker build -t prod_iot_frontend .

cd ../pg
echo "building prod iot pg"
docker build -t prod_iot_pg .

# back to resources again
cd ..

echo "starting containers ..."

docker run --restart always --name prod_iot_pg --net prod-network -p 15432:5432 -d -e POSTGRES_USER=iot -e POSTGRES_DB=iot -e POSTGRES_PASSWORD=iot prod_iot_pg
docker run --restart always --name prod_iot_frontend --net prod-network -p 8081:80 -d prod_iot_frontend
docker run --name prod_iot_backend --net prod-network -p 6001:6001 -d prod_iot_backend

