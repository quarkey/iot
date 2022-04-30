create table if not exists events (
    id serial primary key,
    category text not null,
    message text not null,
    event_time timestamp NOT NULL DEFAULT now()::timestamp(0)
);
