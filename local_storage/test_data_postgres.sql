insert into cluster (id, name) values (default, '00000000-0000-0000-0000-000000000000');
insert into cluster (id, name) values (default, '00000000-0000-0000-0000-000000000001');
insert into cluster (id, name) values (default, '00000000-0000-0000-0000-000000000002');
insert into cluster (id, name) values (default, '00000000-0000-0000-0000-000000000003');
insert into cluster (id, name) values (default, '00000000-0000-0000-0000-000000000004');

insert into configuration_profile (id, configuration, changed_at, changed_by, description) values (default, '{"no_op":"X", "watch":["a","b","c"]}', '2019-01-01', 'tester', 'cfg1');
insert into configuration_profile (id, configuration, changed_at, changed_by, description) values (default, '{"no_op":"X", "watch":["a","b","c"]}', '2019-01-01', 'tester', 'cfg2');
insert into configuration_profile (id, configuration, changed_at, changed_by, description) values (default, '{"no_op":"X", "watch":["a","b","c"]}', '2019-10-11', 'tester', 'cfg3');
insert into configuration_profile (id, configuration, changed_at, changed_by, description) values (default, '{"no_op":"Y", "watch":["d","e"]}', '2019-10-11', 'tester', 'cfg3');

insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (0, 0, 0, '2019-01-01', 'tester', 0, 'no reason');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (1, 0, 1, '2019-01-01', 'tester', 0, 'no reason');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (2, 0, 2, '2019-01-01', 'tester', 1, 'no reason');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (3, 1, 1, '2019-01-01', 'tester', 1, 'no reason');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (4, 2, 2, '2019-10-11', 'tester', 1, 'no reason so far');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (5, 3, 0, '2019-10-11', 'tester', 0, 'disabled one');
insert into operator_configuration (id, cluster, configuration, changed_at, changed_by, active, reason) values (6, 4, 1, '2019-10-11', 'tester', 0, 'disabled one');

insert into trigger_type (type, description) values ('must-gather', 'Triggers must-gather operation on selected cluster');

insert into trigger(type, cluster, reason, link, triggered_at, triggered_by, parameters, active, acked_at) values (1, 0, 'reason', 'link', '2019-11-01', 'tester', '{}', 0, '1970-01-01');
insert into trigger(type, cluster, reason, link, triggered_at, triggered_by, parameters, active, acked_at) values (1, 0, 'reason', 'link', '2019-11-01', 'tester', '{}', 1, '1970-01-01');
insert into trigger(type, cluster, reason, link, triggered_at, triggered_by, parameters, active, acked_at) values (1, 1, 'reason', 'link', '2019-11-01', 'tester', '{}', 0, '1970-01-01');
insert into trigger(type, cluster, reason, link, triggered_at, triggered_by, parameters, active, acked_at) values (1, 1, 'reason', 'link', '2019-11-01', 'tester', '{}', 1, '1970-01-01');
