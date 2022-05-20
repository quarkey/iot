create table if not exists controllers (
    id serial primary key,
    sensor_id integer references sensors (id),
    category text not null,
    /*
        switch
        timeswitch
        thresholdswitch
    */
    title text,
    description text,
    items jsonb,
    /*
        switch:
            description
            on=true
        timeswitch:
            description
            time on
            time off
            repeat
        thresholdswitch:
            description
            datasource=d0c0
            operations grather than, less than, equal, not equal
            threshold_limit
            on=true
        
    */
    alert boolean,
    active boolean,
    created_at timestamp NOT NULL DEFAULT now()::timestamp(0)
);
/*
switch on/off
    turn on or turn off

timeswitch on/off at time
    time on     12:30
    time off    13:30
        dur: 3600s
    repeat every 24 hrs
    can work as: 
        - timeswitch
        - duration switch
        - on switch
        - on/off switch
    
thresholdswitch based on values from a particular dataset
    temperature > 25 turn on
    hydro > 30 turn on humidifier
    water level < 10l turn on water pump 
*/