create table data (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    sensor_description text not null,
    serial text not null,
    temp text not null,
    time text not null,
);

-- id  sensor_description  serial  temp
-- 1   arduino network xxx 15

 insert into data('sensor_description','serial','temp','time') values('sensor descr','serial no','23.12','2018-01-01 123301');

{"sensor_description": "lala description","serial": "lala serial","temp": "23.11","time": "2018-0101 123301"}