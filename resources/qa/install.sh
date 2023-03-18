#!/bin/bash

mkdir -p /iot
mkdir -p /iot/volume
mkdir -p /iot/dist
mkdir -p /iot/database

cd ../../
make build

FILE=/iot/exampleconfig.json
if test -f "$FILE"; then
	echo "$FILE exists."
else
	echo "$FILE does not exist"
	cp -v config/exampleconfig.json /iot
fi
ls -la /iot
cp -v api /iot
cp -v -r database /iot/
cp -v -r client/dist/iot-ng/* /iot/dist
cp -v -r /iot/dist/* /var/www/html
cd -

# installing systemd service
STATUS="$(systemctl is-active iot_server.service)"
if [ "${STATUS}" = "active" ]; then
	echo " iot systemd service is active, stopping and removing before installing"
	systemctl stop iot_server.service
	systemctl disable iot_server.service
	rm /etc/systemd/system/iot_server.service
fi

echo " installing system service iot_server.service"
cp -v systemd/iot_server.service /etc/systemd/system
systemctl daemon-reload
systemctl start iot_server.service
systemctl status iot_server.service

echo "iot qa backend installed"
