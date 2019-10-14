insert into cluster (id, name) values (0, 'cluster0');
insert into cluster (id, name) values (1, 'cluster1');
insert into cluster (id, name) values (2, 'cluster2');
insert into cluster (id, name) values (3, 'cluster3');
insert into cluster (id, name) values (4, 'cluster4');

insert into configuration_profile (id, configuration, changed_at, changed_by, description) values (0, 'configuration1', '2019-01-01', 'tester', 'cfg1');
insert into configuration_profile (id, configuration, changed_at, changed_by, description) values (1, 'configuration2', '2019-01-01', 'tester', 'cfg2');
insert into configuration_profile (id, configuration, changed_at, changed_by, description) values (2, 'configuration3', '2019-10-11', 'tester', 'cfg3');

insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (0, 0, 0, '2019-01-01', 'tester', 1, 'no reason');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (1, 1, 1, '2019-01-01', 'tester', 1, 'no reason');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (2, 2, 2, '2019-10-11', 'tester', 1, 'no reason so far');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (3, 3, 0, '2019-10-11', 'tester', 0, 'disabled one');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (4, 4, 1, '2019-10-11', 'tester', 0, 'disabled one');
