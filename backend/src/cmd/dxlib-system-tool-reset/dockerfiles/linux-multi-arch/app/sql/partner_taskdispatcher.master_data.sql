create schema master_data;

-- Area is SOR Area
create table master_data.area
(
    id                           bigserial primary key,                              -- id auto numbering
    type                         varchar(50),                                        -- tipe master data
    code                         varchar(50) unique,                                 -- code master data (jika ada)
    name                         varchar(100),                                       -- nama master data yang disimpan
    value                        varchar(255),                                       -- value master data yang ditampilkan
    description                  varchar(500),                                       -- deskripsi master data
    cost_center                  varchar(50),                                        -- cost center master data (jika ada) [original varchar(50)]
    parent_group                 varchar(100),                                       -- group parent master data
    parent_value                 varchar(255),                                       -- parent master data merefer ke field mana

    status                       varchar(50),                                        -- status master data
    attribute1                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute2                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute3                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute4                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute5                   varchar(255),                                       -- kolom tambahan jika dibutuhkan

    create_method                varchar(255)                      default 'MANUAL', -- MANUAL, ETL

    is_deleted                   boolean                  not null default false,
    created_at                   timestamp with time zone not null default now(),
    created_by_user_id           varchar(255)             not null default '',
    created_by_user_nameid       varchar(255)             not null default '',
    last_modified_at             timestamp with time zone not null default now(),
    last_modified_by_user_id     varchar(255)             not null default '',
    last_modified_by_user_nameid varchar(255)             not null default ''
);

-- Location is Country Administrative Area
create table master_data.location
(
    id                           bigserial primary key,                              -- id auto numbering
    type                         varchar(50),                                        -- tipe master data
    code                         varchar(50) unique,                                 -- code master data (jika ada)
    name                         varchar(255),                                       -- nama master data yang disimpan
    value                        varchar(255),                                       -- value master data yang ditampilkan
    description                  varchar(500),                                       -- deskripsi master data
    parent_group                 varchar(100),                                       -- parent type
    parent_value                 varchar(255),                                       -- parent code
    status                       varchar(50),                                        -- status master data
    attribute1                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute2                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute3                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute4                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute5                   varchar(255),                                       -- kolom tambahan jika dibutuhkan

    create_method                varchar(255)                      default 'MANUAL', -- MANUAL, ETL

    is_deleted                   boolean                  not null default false,
    created_at                   timestamp with time zone not null default now(),
    created_by_user_id           varchar(255)             not null default '',
    created_by_user_nameid       varchar(255)             not null default '',
    last_modified_at             timestamp with time zone not null default now(),
    last_modified_by_user_id     varchar(255)             not null default '',
    last_modified_by_user_nameid varchar(255)             not null default ''
);

create MATERIALIZED view master_data.mv_location_province as
select l.*
from master_data.location l
where l.type = 'PROVINCE';
CREATE UNIQUE INDEX ON master_data.mv_location_province (code);



create table master_data.customer_ref
(
    id                           bigserial primary key,                              -- id auto numbering
    type                         varchar(50),                                        -- tipe master data
    code                         varchar(50),                                        -- code master data (jika ada)
    name                         varchar(100),                                       -- nama master data yang disimpan
    value                        varchar(255),                                       -- value master data yang ditampilkan
    description                  varchar(500),                                       -- deskripsi master data
    parent_group                 varchar(100),                                       -- group parent master data
    parent_value                 varchar(255),                                       -- parent master data merefer ke field mana
    status                       varchar(50),                                        -- status master data
    attribute1                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute2                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute3                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute4                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute5                   varchar(255),                                       -- kolom tambahan jika dibutuhkan

    create_method                varchar(255)                      default 'MANUAL', -- MANUAL, ETL

    is_deleted                   boolean                  not null default false,
    created_at                   timestamp with time zone not null default now(),
    created_by_user_id           varchar(255)             not null default '',
    created_by_user_nameid       varchar(255)             not null default '',
    last_modified_at             timestamp with time zone not null default now(),
    last_modified_by_user_id     varchar(255)             not null default '',
    last_modified_by_user_nameid varchar(255)             not null default '',

    unique (type, code)
);



create table master_data.global_lookup
(
    id                           bigserial primary key,                              -- id auto numbering
    type                         varchar(50),                                        -- tipe master data
    code                         varchar(50) unique,                                 -- code master data (jika ada)
    name                         varchar(100),                                       -- nama master data yang disimpan
    value                        varchar(255),                                       -- value master data yang ditampilkan
    description                  varchar(500),                                       -- deskripsi master data
    parent_group                 varchar(100),                                       -- group parent master data
    parent_value                 varchar(255),                                       -- parent master data merefer ke field mana
    status                       varchar(50),                                        -- status master data
    attribute1                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute2                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute3                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute4                   varchar(255),                                       -- kolom tambahan jika dibutuhkan
    attribute5                   varchar(255),                                       -- kolom tambahan jika dibutuhkan

    create_method                varchar(255)                      default 'MANUAL', -- MANUAL, ETL

    is_deleted                   boolean                  not null default false,
    created_at                   timestamp with time zone not null default now(),
    created_by_user_id           varchar(255)             not null default '',
    created_by_user_nameid       varchar(255)             not null default '',
    last_modified_at             timestamp with time zone not null default now(),
    last_modified_by_user_id     varchar(255)             not null default '',
    last_modified_by_user_nameid varchar(255)             not null default '',

    unique (type, code)
);

create table master_data.rs_customer_sector
(
    id                           bigserial primary key,
    code                         varchar(255) unique,
    name                         varchar(255) unique,

    create_method                varchar(255)                      default 'MANUAL', -- MANUAL, ETL

    is_deleted                   boolean                  not null default false,
    created_at                   timestamp with time zone not null default now(),
    created_by_user_id           varchar(255)             not null default '',
    created_by_user_nameid       varchar(255)             not null default '',
    last_modified_at             timestamp with time zone not null default now(),
    last_modified_by_user_id     varchar(255)             not null default '',
    last_modified_by_user_nameid varchar(255)             not null default ''
);

create view master_data.vw_customer_segment as
select *
from master_data.customer_ref
where type = 'CUSTOMER_SEGMENT';

create view master_data.vw_customer_type as
select *
from master_data.customer_ref
where type = 'CUSTOMER_TYPE';

create view master_data.vw_payment_scheme as
select *
from master_data.global_lookup
where type = 'SKEMA_BAYAR';
