-- Copyright 2020 Red Hat, Inc
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--      http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

--PRAGMA foreign_keys = ON;
--CREATE SCHEMA public;

create table cluster (
    ID      serial primary key,
    name    text not null
);

ALTER SEQUENCE cluster_id_seq MINVALUE 0 RESTART WITH 0;

create table configuration_profile (
    ID            serial primary key,
    configuration varchar not null,
    changed_at    timestamp,
    changed_by    varchar,
    description   varchar
);

ALTER SEQUENCE configuration_profile_id_seq MINVALUE 0 RESTART WITH 0;

create table operator_configuration (
    ID            serial primary key,
    cluster       integer not null,
    configuration integer not null,
    changed_at    timestamp,
    changed_by    varchar,
    active        integer,
    reason        varchar,
    CONSTRAINT fk_cluster
        foreign key(cluster)
        references cluster(ID)
        on delete cascade,
    CONSTRAINT fk_configuration
        foreign key (configuration)
        references configuration_profile(ID)
        on delete cascade
);

create table trigger_type (
    ID            serial primary key ,
    type          varchar not null,
    description   varchar
);

create table trigger (
    ID            serial primary key,
    type          integer not null,
    cluster       integer not null,
    reason        varchar,
    link          varchar,
    triggered_at  timestamp,
    triggered_by  varchar,
    acked_at      timestamp,
    parameters    varchar,
    active        integer,
    CONSTRAINT fk_type
        foreign key (type)
        references trigger_type(ID),
    CONSTRAINT fk_cluster
        foreign key(cluster)
        references cluster(ID)
);
