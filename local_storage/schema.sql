create table cluster (
    ID      integer primary key asc,
    cluster text not null
);

create table configuration_profile (
    ID            integer primary key asc,
    configuration varchar not null,
    changed_at    datetime,
    changed_by    varchar,
    description   varchar
);

create table operator_configuration (
    ID            integer primary key asc,
    cluster       integer not null,
    configuration integer not null,
    changed_at    datetime,
    changed_by    varchar,
    active        integer,
    reason        varchar,
    foreign key(cluster) references cluster(ID)
    foreign key(configuration) references configuration_profile(ID)
);

