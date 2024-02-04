create table if not exists iot.metrics (
    id serial primary key,
    metric_id integer not null,
    name text,
    help text,
    value numeric,
    created_at text not null default now()
);