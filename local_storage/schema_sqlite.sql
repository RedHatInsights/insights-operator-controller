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

PRAGMA foreign_keys = ON;

create table cluster (
    ID      integer primary key asc,
    name    text not null
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
    CONSTRAINT fk_cluster
        foreign key(cluster)
        references cluster(ID)
        on delete cascade
    CONSTRAINT fk_configuration
        foreign key (configuration)
        references configuration_profile(ID)
        on delete cascade
);

create table trigger_type (
    ID            integer primary key asc,
    type          varchar not null,
    description   varchar
);

create table trigger (
    ID            integer primary key asc,
    type          integer not null,
    cluster       integer not null,
    reason        varchar,
    link          varchar,
    triggered_at  datetime,
    triggered_by  varchar,
    acked_at      datetime,
    parameters    varchar,
    active        integer,
    CONSTRAINT fk_type
        foreign key (type)
        references trigger_type(ID)
        on delete cascade
    CONSTRAINT fk_cluster
        foreign key(cluster)
        references cluster(ID)
        on delete cascade
);
