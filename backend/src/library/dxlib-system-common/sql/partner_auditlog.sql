CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

/*create table log.event
(
    id                           serial primary key,
      uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    type                         varchar(255)             not null, -- API, LOG
    at                           timestamp with time zone not null,
    request_remote_address       varchar(255)             not null,
    request_method               varchar(255)             not null,
    request_url                  varchar(255)             not null,
    request_input_data           jsonb,
    user_id                      bigint,
    user_uid                     varchar(1024),
    user_loginid                 varchar(255),
    user_fullname                varchar(255),
    session_object               jsonb,
    response_status              int                      not null,
    response_message             varchar(255)             not null,
    response_output_data         jsonb,
    is_deleted                   boolean                  not null default false,
    created_at                   timestamp with time zone not null default now(),
    created_by_user_id           varchar(255)             not null default '',
    created_by_user_nameid       varchar(255)             not null default '',
    last_modified_at             timestamp with time zone not null default now(),
    last_modified_by_user_id     varchar(255)             not null default '',
    last_modified_by_user_nameid varchar(255)             not null default ''
);*/


create schema audit_log;

create table audit_log.error_log
(
    id        serial primary key,
    uid       varchar(1024) not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    at        timestamp with time zone,
    log_level varchar(255),
    prefix    varchar(1024),
    location  varchar(1024),
    message   varchar(32768),
    stack     varchar(32768)
);

create table audit_log.user_activity_log
(
    id                        serial primary key,
    uid                       varchar(1024) not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    api_title                 varchar(1024),
    method                    varchar(255),
    api_url                   varchar(1024),
    start_time                timestamp with time zone,
    end_time                  timestamp with time zone,
    status_code               int,
    ip_address                varchar(255),
    error_message             varchar(1024),
    user_id                   bigint,
    user_uid                  varchar(1024),
    user_loginid              varchar(255),
    user_fullname             varchar(255),
    user_roles                jsonb,
    user_effective_privileges jsonb,
    activity_name             varchar(255),
    activity_result_status    varchar(255),
    activity_result_message   varchar(255),
    activity_input            jsonb,
    activity_output           jsonb
);
