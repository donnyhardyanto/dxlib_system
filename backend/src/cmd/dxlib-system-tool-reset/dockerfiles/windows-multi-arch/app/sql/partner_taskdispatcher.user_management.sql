CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

create schema user_management;

create table user_management.user
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    loginid                      varchar(255)             not null unique,
    email                        varchar(255)             not null        default '',
    fullname                     varchar(255)             not null        default '',
    phonenumber                  varchar(255)             not null        default '',
    status                       varchar(255)             not null        default 'ACTIVE', -- ACTIVE, SUSPENDED, DELETED
    attribute                    varchar(1024)            not null        default '',
    identity_number              varchar(255) unique,
    identity_type                varchar(255)             not null        default '',
    gender                       varchar(1),                                                -- M, F
    address_on_identity_card     varchar(1024),
    must_change_password         boolean                  not null        default false,
    is_avatar_exist              bool                     not null        default false,
    utag                         varchar(255) unique,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table user_management.user_password
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    user_id                      bigint                   not null references user_management.user (id),
    value                        varchar(4096)            not null,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table user_management.user_message
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    fcm_application_id           bigint,
    fcm_message_id               bigint,
    sent_at                      timestamp with time zone,
    arrive_at                    timestamp with time zone,
    read_at                      timestamp with time zone,
    user_id                      bigint                   not null references user_management.user (id),
    title                        text                     not null,
    body                         text                     not null,
    data                         jsonb                    not null,
    is_read                      boolean                  not null        default false,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table user_management.organization
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    code                         varchar(255)             not null unique,
    name                         varchar(1024)            not null unique,
    parent_id                    bigint references user_management.organization (id),
    type                         varchar(1024)            not null,                         -- OWNER, CONTRACTOR, SUBCONTRACTOR
    address                      varchar(1024)            not null,
    npwp                         varchar(255),
    email                        varchar(255),
    phonenumber                  varchar(255),
    status                       varchar(255)             not null        default 'ACTIVE', --ACTIVE, SUSPENDED, DELETED
    auth_source1                 varchar(255),                                              -- LDAP1, LDAP2, INTERNAL, (EMPTY/NULL)
    auth_source2                 varchar(255),
    attribute1                   varchar(1024),
    attribute2                   varchar(1024),
    utag                         varchar(255) unique,
    tags                         varchar(1024),
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table user_management.user_organization_membership
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    user_id                      bigint                   not null references user_management.user (id),
    organization_id              bigint                   not null references user_management.organization (id),
    membership_number            varchar(255)             not null        default '',
    order_index                  integer                  not null        default 0,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    -- unique (user_id, organization_id), -- UserOrganizationMembershipTypeMultipleOrganizationPerUser: one user can only have one organization
    unique (user_id) -- UserOrganizationMembershipTypeOneOrganizationPerUser: one user can only have one organization
);

create view user_management.v_user_organization_membership as
select a.*,
       b.uid          as organization_uid,
       b.name         as organization_name,
       b.type         as organization_type,
       b.address      as organization_address,
       b.status       as organization_state,
       b.auth_source1 as organization_auth_source1,
       b.auth_source2 as organization_auth_source2,
       b.attribute1   as organization_attribute1,
       b.attribute2   as organization_attribute2
from user_management.user_organization_membership a
         left join user_management.organization b on a.organization_id = b.id;

create view user_management.v_user as
select a.*,
       uom.membership_number,
       uom.organization_id,
       uom.organization_uid,
       uom.organization_name,
       uom.organization_type,
       uom.organization_address,
       uom.organization_state,
       uom.organization_auth_source1,
       uom.organization_auth_source2,
       uom.organization_attribute1,
       uom.organization_attribute2
from user_management.user a
         left join user_management.v_user_organization_membership uom on a.id = uom.user_id;

create table user_management.role
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    organization_types           JSON, -- - NULL: for universal role that can be applied to any organization type, OWNER, PARTNER, VENDOR, or array of string like ["PARTNER","VENDOR"]
    nameid                       varchar(255)             not null unique,
    name                         varchar(255)             not null,
    description                  varchar(255)             not null,
    utag                         varchar(255) unique,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table user_management.organization_role
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    organization_id              bigint                   not null references user_management.organization (id),
    role_id                      bigint                   not null references user_management.role (id),
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (organization_id, role_id)
);

create view user_management.v_organization_role as
select a.*,
       r.uid                as role_uid,
       r.organization_types as role_organization_types,
       r.nameid             as role_nameid,
       r.name               as role_name,
       r.description        as role_description,
       r.utag               as role_utag,
       o.uid                as organization_uid,
       o.name               as organization_name,
       o.type               as organization_type,
       o.address            as organization_address,
       o.status             as organization_state,
       o.auth_source1       as organization_auth_source1,
       o.auth_source2       as organization_auth_source2,
       o.attribute1         as organization_attribute1,
       o.attribute2         as organization_attribute2
from user_management.organization_role a
         join user_management.role r on a.role_id = r.id
         join user_management.organization o on a.organization_id = o.id;

create table user_management.privilege
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    nameid                       varchar(255)             not null unique,
    name                         varchar(255)             not null,
    description                  varchar(255)             not null,
    utag                         varchar(255) unique,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table user_management.role_privilege
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    role_id                      bigint                   not null references user_management.role (id),
    privilege_id                 bigint                   not null references user_management.privilege (id),
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (role_id, privilege_id)
);

create view user_management.v_role_privilege
as
select a.*,
       b.nameid      as privilege_nameid,
       b.name        as privilege_name,
       b.description as privilege_description,
       b.is_deleted  as privilege_is_deleted
from user_management.role_privilege a
         left join user_management.privilege b on a.privilege_id = b.id;

/* Note: Even the structure allowed to a user to have same role in multiple organizations, but the system will only allow one organization per one user and role combination. */
create table user_management.user_role_membership
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    user_id                      bigint                   not null references user_management.user (id),
    organization_id              bigint                   not null references user_management.organization (id),
    role_id                      bigint                   not null references user_management.role (id),
    data                         jsonb,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (user_id, role_id)
);

create view user_management.v_user_role_membership as
select a.*,
       r.nameid      as role_nameid,
       r.name        as role_name,
       r.description as role_description
from user_management.user_role_membership a
         join user_management.role r on a.role_id = r.id;


create table user_management.menu_item
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    parent_id                    bigint,
    nameid                       varchar(255)             not null,
    name                         varchar(255)             not null,
    composite_nameid             varchar(4096)            not null unique,
    item_index                   integer                  not null        default 0,
    privilege_id                 bigint references user_management.privilege (id),
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create view user_management.v_menu_item as
select a.*,
       b.nameid      as privilege_nameid,
       b.name        as privilege_name,
       b.description as privilege_description,
       b.is_deleted  as privilege_is_deleted
from user_management.menu_item a
         left join user_management.privilege b on a.privilege_id = b.id;


CREATE OR REPLACE FUNCTION user_management.update_composite_nameid() RETURNS TRIGGER AS
'
    DECLARE
        parent_composite_nameid varchar(4096);
    BEGIN
        IF NEW.parent_id IS NOT NULL THEN
            SELECT composite_nameid
            INTO parent_composite_nameid
            FROM user_management.menu_item
            WHERE id = NEW.parent_id;
            NEW.composite_nameid := parent_composite_nameid || ''.'' || NEW.nameid;
        ELSE
            NEW.composite_nameid := NEW.nameid;
        END IF;
        RETURN NEW;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER trg_update_composite_nameid
    BEFORE INSERT OR UPDATE
    ON user_management.menu_item
    FOR EACH ROW
EXECUTE FUNCTION user_management.update_composite_nameid();

/* standard role, privilege, organization_roles and role_privilege for organization type OWNER */

insert into user_management.role (organization_types, nameid, name, description, utag)
values ('[
  "OWNER"
]', 'SUPER-ADMINISTRATOR', 'Super Administrator', 'Super administrator role', 'SUPER-ADMINISTRATOR');

insert into user_management.privilege (nameid, name, description, utag)
values ('EVERYTHING', 'Every thing', 'Everything allowed. This is only for Super administrator role.', 'EVERYTHING');

insert into user_management.role_privilege(role_id, privilege_id)
values ((select id from user_management.role where utag = 'SUPER-ADMINISTRATOR'),
        (select id from user_management.privilege where utag = 'EVERYTHING'));

/* initial user and organization for OWNER */

insert into user_management.organization (code, name, parent_id, type, address, status, auth_source1, auth_source2,
                                          utag)
values ('0001', 'Owner', null, 'OWNER', '', 'ACTIVE', '', '', 'OWNER');

insert into user_management.organization_role (organization_id, role_id)
values ((select id from user_management.organization where utag = 'OWNER'),
        (select id from user_management.role where utag = 'SUPER-ADMINISTRATOR'));

insert into user_management.user (loginid, fullname, status, created_by_user_id, created_by_user_nameid,
                                  last_modified_by_user_id, last_modified_by_user_nameid, utag)
values ('superadmin', 'Super Administrator', 'ACTIVE', 0, 'SYSTEM', 0, 'SYSTEM', 'SUPERADMIN');

insert into user_management.user_organization_membership (user_id, organization_id)
values ((select id from user_management.user where utag = 'SUPERADMIN'),
        (select id from user_management.organization where utag = 'OWNER'));

insert into user_management.user_role_membership(user_id, organization_id, role_id)
values ((select id from user_management.user where utag = 'SUPERADMIN'),
        (select id from user_management.organization where utag = 'OWNER'),
        (select id from user_management.role where utag = 'SUPER-ADMINISTRATOR'));
