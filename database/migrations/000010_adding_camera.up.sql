create table if not exists iot.camera (
    id serial primary key,
    sensor_id integer,
    title text,
    description text,
    description_long text,
    hostname text,
    project_name text,
    storage_location text,
    interval integer,
    next_capture_time timestamp,
    status text not null default 'new',
    alert boolean not null default false,
    active boolean not null default false,
    created_at text not null default now(),
    updated_at text
);