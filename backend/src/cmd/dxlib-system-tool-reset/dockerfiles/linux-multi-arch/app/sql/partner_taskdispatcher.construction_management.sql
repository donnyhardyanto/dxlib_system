create schema construction_management;

create table construction_management.gas_appliance
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    code                         varchar(255) unique,
    name                         varchar(255) unique,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

insert into construction_management.gas_appliance (code, name)
values ('01', 'Kompor Portable 1 Tungku'),
       ('02', 'Kompor Portable 2 Tungku'),
       ('03', 'Kompor Portable 3 Tungku'),
       ('04', 'Kompor Tanam 1 Tungku'),
       ('05', 'Kompor Tanam 2 Tungku'),
       ('06', 'Kompor Tanam 3 Tungku'),
       ('07', 'Kompor Tanam 4 Tungku'),
       ('08', 'Kompor Low/High Pressure'),
       ('09', 'Oven'),
       ('10', 'Water Header'),
       ('11', 'Gas Dryer Pakaian'),
       ('12', 'Gas Rice Cooker');

create table construction_management.tapping_saddle_appliance
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    code                         varchar(255) unique,
    name                         varchar(255) unique,
    is_deleted                   boolean                  not null default false,
    created_at                   timestamp with time zone not null default now(),
    created_by_user_id           varchar(255)             not null default '',
    created_by_user_nameid       varchar(255)             not null default '',
    last_modified_at             timestamp with time zone not null default now(),
    last_modified_by_user_id     varchar(255)             not null default '',
    last_modified_by_user_nameid varchar(255)             not null default ''
);

insert into construction_management.tapping_saddle_appliance (code, name)
values ('01', '63 mm x 20 mm'),
       ('02', '90 mm x 20 mm'),
       ('03', '125 mm x 20 mm'),
       ('04', '180 mm x 20 mm'),
       ('05', '63 mm x 32 mm'),
       ('06', '90 mm x 32 mm'),
       ('07', '125 mm x 32 mm'),
       ('08', '180 mm x 32 mm');

create table construction_management.meter_appliance_type
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    code                         varchar(255) unique,
    name                         varchar(255) unique,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

insert into construction_management.meter_appliance_type (code, name)
values ('01', 'Meter Konvensional'),
       ('02', 'Smart Meter');

create table construction_management.regulator_appliance
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    code                         varchar(255) unique,
    name                         varchar(255) unique,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

insert into construction_management.regulator_appliance (code, name)
values ('01', 'RV12'),
       ('02', 'RV20'),
       ('03', 'RV47'),
       ('04', 'RV48');

create table construction_management.g_size
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    code                         varchar(255) unique,
    name                         varchar(255) unique,
    qmin                         float,
    qmax                         float,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

insert into construction_management.g_size (code, name, qmin, qmax)
values ('01', 'G1.6', 0.02, 2.5),
       ('02', 'G2.5', 0.03, 4),
       ('03', 'G4', 0.04, 6),
       ('04', 'G6', 0.06, 10),
       ('05', 'G10', 0.1, 16),
       ('06', 'G16', 0.16, 25),
       ('07', 'G25', 0.25, 40);
