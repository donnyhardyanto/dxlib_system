CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create schema configuration;

create table configuration.external_system
(
    id                           bigserial                not null primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    nameid                       varchar(255)             not null unique,
    type                         varchar(255)             not null, -- LDAP, SMTP, RelyOn
    configuration                jsonb,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);


create schema settings;

create table settings.general_template
(
    id                           bigserial                not null primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    nameid                       varchar(255)             not null unique,
    content_type                 varchar(255)             not null, -- text/plain, text/html
    subject                      varchar(255)             not null,
    body                         varchar(8096)            not null,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table settings.email_template
(
    id                           bigserial                not null primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    nameid                       varchar(255)             not null unique,
    content_type                 varchar(255)             not null, -- text/plain, text/html
    subject                      varchar(255)             not null,
    body                         varchar(8096)            not null,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table settings.sms_template
(
    id                           bigserial                not null primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    nameid                       varchar(255)             not null unique,
    content_type                 varchar(255)             not null, -- text/plain, text/html
    subject                      varchar(255)             not null,
    body                         varchar(8096)            not null,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);
