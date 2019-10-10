insert into cluster (id, cluster) values (0, 'cluster0');
insert into cluster (id, cluster) values (1, 'cluster1');

insert into configuration_profile (id, configuration, changed_at, changed_by, description) values (0, 'configuration1', '2019-01-01', 'tester', 'cfg1');
insert into configuration_profile (id, configuration, changed_at, changed_by, description) values (1, 'configuration2', '2019-01-01', 'tester', 'cfg1');

insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (0, 0, 0, '2019-01-01', 'tester', 1, 'no reason');
