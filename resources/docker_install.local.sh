#!/bin/bash

echo "building dist file"
cd ../client
npm install
ng build

cp -r dist/* ../resources/qa/ng/dist
cd ../resources

echo "creating QA docker network"
docker network create qa-network

cd ..
echo "building QA iot backend ..."
docker build -t qa_iot_backend .

cd resources/qa/ng
echo "building QA iot frontend"
docker build -t qa_iot_frontend .

cd ../pg
echo "building QA iot pg"
docker build -t qa_iot_pg .

# back to resources again
cd ../../

echo "starting containers ..."

docker run --restart always --name qa_iot_pg --net qa-network -p 15432:5432 -d -e POSTGRES_USER=iot -e POSTGRES_DB=iot -e POSTGRES_PASSWORD=iot qa_iot_pg
docker run --restart always --name qa_iot_frontend --net qa-network -p 8081:80 -d qa_iot_frontend
docker run --restart always --name qa_iot_backend --net qa-network -p 6002:6002 -d qa_iot_backend
