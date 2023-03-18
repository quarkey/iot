#!/bin/bash

cp iot_server.service /etc/systemd/system

# restart systemd demon and start service
systemctl daemon-reload
systemctl start iot_server.service

