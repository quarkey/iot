-- postgres:

drop table sensordata;
drop table dataset;
drop table sensor;

create table if not exists sensor (
  id serial primary key,
  description text not null, -- aurdino sensor description
  arduino_key text not null, -- unique identifyer key
  created_at timestamp NOT NULL DEFAULT now()::timestamp(0)
);
create table if not exists dataset(
  id serial primary key,
  sensor_id integer references sensor (id),
  description text not null,
  reference text not null,
  intervalsec int not null,
  fields jsonb,
  created_at timestamp NOT NULL DEFAULT now()::timestamp(0)
);
create table if not exists sensordata (
  id serial primary key,
  sensor_id integer references sensor (id),
  dataset_id integer references dataset (id),
  data text not null,
  time timestamp NOT NULL DEFAULT now()::timestamp(0)
);

insert into sensor(description, arduino_key) values('temp og hydro', 'arduino serial');
insert into dataset(sensor_id, description, reference, intervalsec, fields) values(1,'temperatur measurement, growhouse','reference x',1800,'["temp", "hydro"]');
insert into sensordata(sensor_id, dataset_id, data) values(1,currval(pg_get_serial_sequence('dataset', 'id')),'["23.13","59.32","ubro"]');

select * from sensor;
select * from dataset;
select * from sensordata;



--sqlite3
create table sensor (
  id serial primary key,
  descr text not null,
  serial text not null,
  data jsonb not null
);
-- test data
insert into sensordata(descr, serial, data) values('descr','serialno','{ "customer": "John Doe", "items": {"product": "Beer","qty": 6}}'::jsonb);
insert into sensordata(descr, serial, data) values('descr 2','serialno 2','{"abc": "val"}'::jsonb);
insert into sensordata(descr, serial, data) values('temperatur sensor', 'serialyo','{  "sensor description": "Temperature readings",  "hardware description": "arduino uno with temp sensor",  "serial": "a8f5f167f44f4964e6c998dee827110c",  "ip address": "192.168.10.100",  "network mask": "255.255.255.0",  "server": "192.168.10.1",  "encryption key": "8ed358a7da3cc760364909d4aaf7321e",  "record interval": "1800",  "data": {"serial": "a8f5f167f44f4964e6c998dee827110c","temp c": ["33.1","22.1"],"record time": ["113030","1200"]}}'::jsonb);


  select 
    a.id,
    a.sensor_id,
    a.title,
    a.description,
    a.reference,
    a.intervalsec,
    a.fields,
    a.created_at,
    b.title as sensor_title
    from dataset a, sensors b
    where a.sensor_id = b.id;


select
  a.id,
  a.data,
  a.time
from 
  sensordata a,
  dataset b
where
  b.reference='8a1bbddba98a8d8512787d311352d951'
  and b.id = a.dataset_id