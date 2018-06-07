-- sqlite:
create table data (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    sensor_description text not null,
    serial text not null,
    temp text not null,
    time text not null,
);
insert into data('sensor_description','serial','temp','time') values('sensor descr','serial no','23.12','2018-01-01 123301');

-- id  sensor_description  serial  temp
-- 1   arduino network xxx 15
-- {"sensor_description": "lala description","serial": "lala serial","temp": "23.11","time": "2018-0101 123301"}
 

-- postgres:
create table sensordata(
  id serial primary key,
  descr text not null,
  serial text not null,
  data jsonb not null
);

-- test data
insert into sensordata(descr, serial, data) values('descr','serialno','{ "customer": "John Doe", "items": {"product": "Beer","qty": 6}}'::jsonb);
insert into sensordata(descr, serial, data) values('descr 2','serialno 2','{"abc": "val"}'::jsonb);
insert into sensordata(descr, serial, data) values('temperatur sensor', 'serialyo','{  "sensor description": "Temperature readings",  "hardware description": "arduino uno with temp sensor",  "serial": "a8f5f167f44f4964e6c998dee827110c",  "ip address": "192.168.10.100",  "network mask": "255.255.255.0",  "server": "192.168.10.1",  "encryption key": "8ed358a7da3cc760364909d4aaf7321e",  "record interval": "1800",  "data": {"serial": "a8f5f167f44f4964e6c998dee827110c","temp c": ["33.1","22.1"],"record time": ["113030","1200"]}}'::jsonb);
  