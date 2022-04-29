#!/bin/bash

echo "building dist file"
cd ../client
ng build --configuration production-m1mini

cp -r dist/* ../resources/prod/ng/dist
cd ../resources

echo "creating prod docker network"
docker network create prod-network

cd ..
echo "building prod iot backend ..."
docker build -t prod_iot_backend .

cd resources/prod/ng
echo "building prod iot frontend"
docker build -t prod_iot_frontend .

cd ../pg
echo "building prod iot pg"
docker build -t prod_iot_pg .

# back to resources again
cd ..

echo "starting containers ..."

docker run --restart always --name prod_iot_pg --net prod-network -v /Users/slundin/devel/iot/resources/pg_data/prod:/var/lib/postgresql/data  -p 5432:5432 -d -e POSTGRES_USER=iot -e POSTGRES_DB=iot -e POSTGRES_PASSWORD=iot prod_iot_pg
docker run --restart always --name prod_iot_frontend --net prod-network -p 8080:80 -d prod_iot_frontend
docker run --restart always --name prod_iot_backend --net prod-network -p 6001:6001 -d prod_iot_backend

