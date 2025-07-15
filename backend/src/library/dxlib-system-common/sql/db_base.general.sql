CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create schema general;

create table general.property
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    nameid                       varchar(255) unique      not null,
    type                         varchar(255)             not null,
    value                        jsonb                    not null,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

INSERT INTO general.property (nameid, type, value)
VALUES ('SYSTEM-NAME', 'STRING', '{
  "value": "PGN Partner Task Dispatcher"
}'::JSONB),
       ('SYSTEM-VERSION', 'STRING', '{
         "value": "1.0.0"
       }'::JSONB),
       ('PREKEY_TTL_SECOND', 'INT', '{
         "value": 300
       }'::JSONB),
       ('SESSION_TTL_SECOND', 'INT', '{
         "value": 3000
       }'::JSONB),
       ('RELYON_INBOUND_SESSION_TTL_SECOND', 'INT', '{
         "value": 86400
       }'::JSONB);


create table general.announcement
(
    id                           bigserial                not null primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    timestamp                    timestamp with time zone not null        default now(),
    title                        varchar(255)             not null,
    content                      varchar(2048)            not null,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);
