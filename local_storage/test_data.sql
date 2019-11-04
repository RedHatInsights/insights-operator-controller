insert into cluster (id, name) values (0, 'cluster0');
insert into cluster (id, name) values (1, 'cluster1');
insert into cluster (id, name) values (2, 'cluster2');
insert into cluster (id, name) values (3, 'cluster3');
insert into cluster (id, name) values (4, 'cluster4');

insert into configuration_profile (id, configuration, changed_at, changed_by, description) values (0, '{"no_op":"X", "watch":["a","b","c"]}', '2019-01-01', 'tester', 'cfg1');
insert into configuration_profile (id, configuration, changed_at, changed_by, description) values (1, '{"no_op":"X", "watch":["a","b","c"]}', '2019-01-01', 'tester', 'cfg2');
insert into configuration_profile (id, configuration, changed_at, changed_by, description) values (2, '{"no_op":"X", "watch":["a","b","c"]}', '2019-10-11', 'tester', 'cfg3');
insert into configuration_profile (id, configuration, changed_at, changed_by, description) values (3, '{"no_op":"Y", "watch":["d","e"]}', '2019-10-11', 'tester', 'cfg3');

insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (0, 0, 0, '2019-01-01', 'tester', 0, 'no reason');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (1, 0, 1, '2019-01-01', 'tester', 0, 'no reason');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (2, 0, 2, '2019-01-01', 'tester', 1, 'no reason');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (3, 1, 1, '2019-01-01', 'tester', 1, 'no reason');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (4, 2, 2, '2019-10-11', 'tester', 1, 'no reason so far');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (5, 3, 0, '2019-10-11', 'tester', 0, 'disabled one');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (6, 4, 1, '2019-10-11', 'tester', 0, 'disabled one');

insert into trigger_type (type, description) values ('must-gather', 'Triggers must-gather operation on selected cluster');

insert into trigger(type, cluster, reason, link, triggered_at, triggered_by, parameters, active) values (1, 0, 'reason', 'link', '2019-11-01', 'tester', '{}', 0);
insert into trigger(type, cluster, reason, link, triggered_at, triggered_by, parameters, active) values (1, 0, 'reason', 'link', '2019-11-01', 'tester', '{}', 1);
insert into trigger(type, cluster, reason, link, triggered_at, triggered_by, parameters, active) values (1, 1, 'reason', 'link', '2019-11-01', 'tester', '{}', 0);
insert into trigger(type, cluster, reason, link, triggered_at, triggered_by, parameters, active) values (1, 1, 'reason', 'link', '2019-11-01', 'tester', '{}', 1);
