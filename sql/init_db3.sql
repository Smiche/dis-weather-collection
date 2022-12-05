
create extension if not exists postgres_fdw;

create server meas_db2_fdw foreign data wrapper postgres_fdw options (
    dbname 'weatherdb',
    host 'db2',
    port '5432'
);

create server meas_db1_fdw foreign data wrapper postgres_fdw options (
    dbname 'weatherdb',
    host 'db1',
    port '5432'
);

create user mapping for "dis-db-user" server meas_db2_fdw options (
   "user" 'dis-db-user'
);

create user mapping for "dis-db-user" server meas_db1_fdw options (
   "user" 'dis-db-user'
);

create table phenomenon_type (
    id bigint generated by default as identity, 
    name text not null, 
    primary key(id)
);

create subscription phen_subscription_2 
connection 'host=db1 port=5432 dbname=weatherdb user=dis-db-user' 
publication phen_publication;

create table unit (
    id bigint generated by default as identity, 
    name text not null, 
    abbreviation text not null, 
    primary key(id)
);

create subscription unit_subscription_2 
connection 'host=db1 port=5432 dbname=weatherdb user=dis-db-user' 
publication unit_publication;

create table organization (
    id bigint generated by default as identity, 
    name text not null, 
    country text not null
    -- primary key(id, country)
)
partition by list (country);

create table organization_local 
partition of organization for values in ('Norway');

create foreign table organization_db2 
partition of organization for values in ('Sweden')
server meas_db2_fdw options (table_name 'organization_local');

create foreign table organization_db1 
partition of organization for values in ('Finland')
server meas_db1_fdw options (table_name 'organization_local');

create table users_local (
    id bigint generated by default as identity, 
    organization bigint not null, -- references organization(id) not null, 
    name text not null, 
    email text, 
    password text not null, 
    primary key (id)
);

create table station_local (
    id bigint generated by default as identity, 
    name text not null, 
    number bigint not null, 
    organization bigint not null, -- references organization(id) not null, 
    type text not null, 
    latitude double precision not null, 
    longitude double precision not null, 
    altitude double precision not null, 
    city text not null, 
    primary key (id)
);

create foreign table station_db2 (
    id bigint, 
    name text, 
    number bigint, 
    organization bigint,
    type text, 
    latitude double precision, 
    longitude double precision, 
    altitude double precision, 
    city text
)
server meas_db2_fdw options (table_name 'station_local');

create foreign table station_db1 (
    id bigint, 
    name text, 
    number bigint, 
    organization bigint,
    type text, 
    latitude double precision, 
    longitude double precision, 
    altitude double precision, 
    city text
)
server meas_db1_fdw options (table_name 'station_local');

create view station_all as 
    select *, 'Norway' as country from station_local
    union all select *, 'Sweden' as country from station_db2
    union all select *, 'Finland' as country from station_db1;

create table device_local (
    id bigint generated by default as identity, 
    station bigint references station_local(id) not null, 
    name text not null, 
    primary key (id)
);

create foreign table device_db2 (
    id bigint, 
    station bigint, 
    name text
)
server meas_db2_fdw options (table_name 'device_local');

create foreign table device_db1 (
    id bigint, 
    station bigint, 
    name text
)
server meas_db1_fdw options (table_name 'device_local');

create view device_all as
    select *, 'Norway' as country from device_local
    union all select *, 'Sweden' as country from device_db2
    union all select *, 'Finland' as country from device_db1;

create table measurement_local (
    id bigint generated by default as identity, 
    device bigint references device_local(id) not null, 
    value double precision not null, 
    time timestamp without time zone not null, 
    type bigint references phenomenon_type(id) not null, 
    unit bigint references unit(id) not null, 
    primary key (id)
);

create index meas_time_idx on measurement_local (time);

create foreign table measurement_db2 (
    id bigint, 
    device bigint, 
    value double precision, 
    time timestamp without time zone, 
    type bigint, 
    unit bigint
) 
server meas_db2_fdw options (table_name 'measurement_local');

create foreign table measurement_db1 (
    id bigint, 
    device bigint, 
    value double precision, 
    time timestamp without time zone, 
    type bigint, 
    unit bigint
) 
server meas_db1_fdw options (table_name 'measurement_local');

create view measurement_all as
    select *, 'Norway' as country from measurement_local
    union all select *, 'Sweden' as country from measurement_db2
    union all select *, 'Finland' as country from measurement_db1;

create materialized view station_info_local as 
    select station_local.id as id, 
        station_local.name as station, 
        type as station_type, 
        organization_local.name as organization, 
        latitude, 
        longitude, 
        altitude, 
        city, 
        organization_local.country as country 
    from station_local 
    join organization_local on station_local.organization = organization_local.id;

create materialized view meas_min_max_day_local as
    select distinct on (station_info_local, min, max, avg, time, phenomenon_type, unit)
        meas_day.id as id, 
        station_info_local.id as station_info, 
        meas_day.min as min, 
        meas_day.max as max, 
        meas_day.avg as avg, 
        meas_day.time as time, 
        phenomenon_type.name as phenomenon_type, 
        unit.abbreviation as unit 
    from (
        select id,
            date_trunc('day', time) as time,
            min(value) over (partition by type, unit, date_trunc('day', time)) as min,
            max(value) over (partition by type, unit, date_trunc('day', time)) as max,
            avg(value) over (partition by type, unit, date_trunc('day', time)) as avg,
            device,
            type,
            unit
        from measurement_local
    ) meas_day
    join device_local on meas_day.device = device_local.id
    join station_info_local on device_local.station = station_info_local.id
    join phenomenon_type on meas_day.type = phenomenon_type.id
    join unit on meas_day.unit = unit.id;

create index meas_view_time_idx on meas_min_max_day_local (time);

create foreign table meas_min_max_day_db1 (
    id bigint, 
    station_info bigint, 
    min double precision, 
    max double precision, 
    avg double precision, 
    time timestamp without time zone, 
    phenomenon_type text, 
    unit text
) 
server meas_db1_fdw options (table_name 'meas_min_max_day_local');

create foreign table meas_min_max_day_db2 (
    id bigint, 
    station_info bigint, 
    min double precision, 
    max double precision, 
    avg double precision, 
    time timestamp without time zone, 
    phenomenon_type text, 
    unit text
) 
server meas_db2_fdw options (table_name 'meas_min_max_day_local');

create materialized view station_info_all as 
    select station_all.id as id, 
        station_all.name as station, 
        type as station_type, 
        organization.name as organization, 
        latitude, 
        longitude, 
        altitude, 
        city, 
        organization.country as country 
    from station_all 
    join organization on station_all.organization = organization.id
        and station_all.country = organization.country;

create materialized view meas_min_max_day_all as
    select *, 'Norway' as country from meas_min_max_day_local
    union all select *, 'Sweden' as country from meas_min_max_day_db2
    union all select *, 'Finland' as country from meas_min_max_day_db1;

create index meas_all_view_time_idx on meas_min_max_day_all (time);

create extension pg_cron;