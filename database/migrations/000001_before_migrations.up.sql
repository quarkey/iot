create schema if not exists iot;

create table if not exists  sensors (
    id serial primary key,
    title text not null,
    description text not null, -- aurdino sensor description
    arduino_key text not null, -- unique identifyer key
    created_at timestamp NOT NULL DEFAULT now()::timestamp(0)
);

create table if not exists datasets (
    id serial primary key,
    sensor_id integer references sensors (id),
    title text not null,
    description text not null,
    reference text not null,
    intervalsec int not null,
    fields jsonb,
    types jsonb,
    created_at timestamp NOT NULL DEFAULT now()::timestamp(0)
);

create table if not exists sensordata (
    id serial primary key,
    sensor_id integer references sensors (id),
    dataset_id integer references datasets (id),
    data jsonb,
    time timestamp NOT NULL DEFAULT now()::timestamp(0)
);