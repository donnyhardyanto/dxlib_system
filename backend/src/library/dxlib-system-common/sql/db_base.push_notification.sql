create schema push_notification;

create table push_notification.fcm_application
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    nameid                       text                     not null unique,
    service_account_data         jsonb                    not null,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table push_notification.fcm_user_token
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    user_id                      bigint                   not null references user_management.user (id),
    fcm_application_id           bigint                   not null references push_notification.fcm_application (id),
    fcm_token                    text                     not null,
    device_type                  text                     not null, -- ANDROID or IOS
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table push_notification.fcm_message
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    fcm_user_token_id            bigint                   not null references push_notification.fcm_user_token (id),
    status                       varchar(255)             not null,
    title                        text                     not null,
    body                         text                     not null,
    data                         jsonb                    not null,
    next_retry_time              timestamp with time zone,
    is_read                      boolean                  not null        default false,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create view push_notification.v_fcm_message as
select m.*,
       ut.user_id,
       ut.fcm_application_id,
       ut.fcm_token,
       ut.device_type
from push_notification.fcm_message m
         left join push_notification.fcm_user_token ut on m.fcm_user_token_id = ut.id

