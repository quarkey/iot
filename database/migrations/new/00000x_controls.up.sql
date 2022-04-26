create table if not exists  controls (
    id serial primary key,
    sensor_id integer references sensors (id),
    dataset_id integer references datasets (id),
    category text not null, -- is enum a better option?
    title text not null,
    description text not null,

    active boolean not null
    created_at timestamp NOT NULL DEFAULT now()::timestamp(0),
    
);

/*

scenarip 1: light switch - on/off switch

insert into controls (sensor_id, category, title, description)
values(1, "switch", "light switch rom 1", "just a normal on/off switch")

scenario 2: light switch - on/off switch controlled by time or threshold by a light sensor (dataset)


*/