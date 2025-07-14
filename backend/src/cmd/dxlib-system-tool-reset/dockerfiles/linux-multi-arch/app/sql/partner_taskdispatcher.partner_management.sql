create schema partner_management;

create table partner_management.role_area
(
    id        bigserial primary key,
    uid       varchar(1024) not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint),
                                                           gen_random_uuid()::text),
    role_id   bigint        not null references user_management.role (id),
    area_code varchar(255)  not null references master_data.area (code),
    unique (role_id, area_code)
);

create view partner_management.v_role_area as
select ra.*,
       a.code  as sales_area_code,
       a.name  as sales_area_name,
       a.value as sales_area_value
from partner_management.role_area ra
         left join master_data.area a on a.code = ra.area_code;

create table partner_management.role_task_type
(
    id           bigserial primary key,
    uid          varchar(1024) not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint),
                                                              gen_random_uuid()::text),
    role_id      bigint        not null references user_management.role (id),
    task_type_id bigint        not null references task_management.task_type (id),
    unique (role_id, task_type_id)
);

create view partner_management.v_role_task_type as
select rtt.*,
       tt.code as task_type_code,
       tt.name as task_type_name
from partner_management.role_task_type rtt
         left join task_management.task_type tt on tt.id = rtt.task_type_id;

create view partner_management.v_role as
select r.*,
       vra.sales_area_code,
       vra.sales_area_name,
       vrtt.task_type_id,
       vrtt.task_type_code,
       vrtt.task_type_name
from user_management.role r
         left join partner_management.v_role_area vra on vra.role_id = r.id
         left join partner_management.v_role_task_type vrtt on vrtt.role_id = r.id;

create view partner_management.v_user_role_membership as
select urm.*,
       vr.uid         as role_uid,
       vr.nameid      as role_nameid,
       vr.name        as role_name,
       vr.description as role_description,
       vr.utag        as role_utag,
       vr.sales_area_code,
       vr.sales_area_name,
       vr.task_type_id,
       vr.task_type_code,
       vr.task_type_name
from user_management.user_role_membership urm
         left join partner_management.v_role vr on vr.id = urm.role_id;

insert into user_management.role (nameid, name, organization_types, description, utag)
values ('FIELD_EXECUTOR', 'Pelaksana Lapangan', '[
  "CONTRACTOR",
  "SUBCONTRACTOR"
]', 'Field Executor', 'FIELD_EXECUTOR');
insert into user_management.role (nameid, name, organization_types, description, utag)
values ('FIELD_SUPERVISOR', 'Supervisor Lapangan', '[
  "CONTRACTOR"
]', 'Field Supervisor', 'FIELD_SUPERVISOR');

create view partner_management.v_organization_executor as
select tor.*,
       o.uid  as organization_uid,
       o.code as organization_code,
       o.name as organization_name,
       o.type as organization_type
from user_management.organization_role tor
         left join user_management.organization o on o.id = tor.organization_id
where tor.role_id = (select id from user_management.role where utag = 'FIELD_EXECUTOR');

create table partner_management.organization_executor_expertise
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    organization_role_id         bigint                   not null references user_management.organization_role (id),
    sub_task_type_id             bigint                   not null references task_management.sub_task_type (id),
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (organization_role_id, sub_task_type_id)
);

create view partner_management.v_organization_executor_expertise as
select oee.*,
       vstt.task_type_id   as task_type_id,
       vstt.task_type_code as task_type_code,
       vstt.task_type_name as task_type_name,
       vstt.code           as sub_task_type_code,
       vstt.name           as sub_task_type_name,
       vstt.full_name      as sub_task_type_full_name,
       vstt.full_code      as sub_task_type_full_code,
       voe.organization_id,
       voe.organization_uid,
       voe.organization_code,
       voe.organization_name,
       voe.organization_type
from partner_management.organization_executor_expertise oee
         left join partner_management.v_organization_executor voe on oee.organization_role_id = voe.id
         left join task_management.v_sub_task_type vstt on vstt.id = oee.sub_task_type_id;


create table partner_management.organization_executor_area
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    organization_role_id         bigint                   not null references user_management.organization_role (id),
    area_code                    varchar(50)              not null references master_data.area (code),
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (organization_role_id, area_code)
);

create view partner_management.v_organization_executor_area as
select oea.*,
       a.type         as area_type,
       a.name         as area_name,
       a.value        as area_value,
       a.description  as area_description,
       a.cost_center  as area_cost_center,
       a.parent_group as area_parent_group,
       a.parent_value as area_parent_value,
       a.is_deleted   as area_is_deleted,
       a.status       as area_status,
       a.attribute1   as area_attribute1,
       a.attribute2   as area_attribute2,
       a.attribute3   as area_attribute3,
       a.attribute4   as area_attribute4,
       a.attribute5   as area_attribute5,
       voe.organization_id,
       voe.organization_uid,
       voe.organization_code,
       voe.organization_name,
       voe.organization_type
from partner_management.organization_executor_area oea
         left join partner_management.v_organization_executor voe on oea.organization_role_id = voe.id
         left join master_data.area a on oea.area_code = a.code;

create table partner_management.organization_executor_location
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    organization_role_id         bigint                   not null references user_management.organization_role (id),
    location_code                varchar(50)              not null references master_data.location (code),
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (organization_role_id, location_code)
);

create view partner_management.v_organization_executor_location as
select oel.*,
       l.type         as location_type,
       l.name         as location_name,
       l.value        as location_value,
       l.description  as location_description,
       l.parent_group as location_parent_group,
       l.parent_value as location_parent_value,
       l.status       as location_status,
       l.attribute1   as location_attribute1,
       l.attribute2   as location_attribute2,
       l.attribute3   as location_attribute3,
       l.attribute4   as location_attribute4,
       l.attribute5   as location_attribute5,
       voe.organization_id,
       voe.organization_uid,
       voe.organization_code,
       voe.organization_name,
       voe.organization_type
from partner_management.organization_executor_location oel
         left join partner_management.v_organization_executor voe on oel.organization_role_id = voe.id
         left join master_data.location l on oel.location_code = l.code;

create view partner_management.v_field_executor
as
select urm.*,
       u.uid                      as user_uid,
       u.loginid                  as user_loginid,
       u.email                    as user_email,
       u.fullname                 as user_fullname,
       u.phonenumber              as user_phonenumber,
       u.status                   as user_status,
       u.identity_number          as user_identity_number,
       u.identity_type            as user_identity_type,
       u.gender                   as user_gender,
       u.address_on_identity_card as user_address_on_identity_card,
       u.is_deleted               as user_is_deleted,
       o.uid                      as organization_uid,
       o.name                     as organization_name,
       o.type                     as organization_type,
       uom.membership_number      as organization_user_membership_number
from user_management.user_role_membership urm
         left join user_management.user u on urm.user_id = u.id
         left join user_management.user_organization_membership uom on u.id = uom.user_id
         left join user_management.organization o on urm.organization_id = o.id
where urm.role_id = (select id from user_management.role where utag = 'FIELD_EXECUTOR');

create table partner_management.field_executor_expertise
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    user_role_membership_id      bigint                   not null references user_management.user_role_membership (id),
    sub_task_type_id             bigint                   not null references task_management.sub_task_type (id),
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (user_role_membership_id, sub_task_type_id)
);

create view partner_management.v_field_executor_expertise as
select fee.*,
       vstt.task_type_id   as task_type_id,
       vstt.task_type_code as task_type_code,
       vstt.task_type_name as task_type_name,
       vstt.code           as sub_task_type_code,
       vstt.name           as sub_task_type_name,
       vstt.full_name      as sub_task_type_full_name,
       vstt.full_code      as sub_task_type_full_code
from partner_management.field_executor_expertise fee
         left join partner_management.v_field_executor vfe on fee.user_role_membership_id = vfe.id
         left join task_management.v_sub_task_type vstt on vstt.id = fee.sub_task_type_id;

create table partner_management.field_executor_location
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    user_role_membership_id      bigint                   not null references user_management.user_role_membership (id),
    location_code                varchar(50)              not null references master_data.location (code),
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (user_role_membership_id, location_code)
);

create view partner_management.v_field_executor_location as
select fel.*,
       l.type         as location_type,
       l.name         as location_name,
       l.value        as location_value,
       l.description  as location_description,
       l.parent_group as location_parent_group,
       l.parent_value as location_parent_value,
       l.status       as location_status,
       l.attribute1   as location_attribute1,
       l.attribute2   as location_attribute2,
       l.attribute3   as location_attribute3,
       l.attribute4   as location_attribute4,
       l.attribute5   as location_attribute5
from partner_management.field_executor_location fel
         left join partner_management.v_field_executor vfe on fel.user_role_membership_id = vfe.id
         left join master_data.location l on fel.location_code = l.code;

CREATE VIEW partner_management.v_field_executor_effective_location AS
SELECT vfe.*,
       vfe.id         AS user_role_membership_id,
       l.type         as location_type,
       l.name         as location_name,
       l.value        as location_value,
       l.description  as location_description,
       l.parent_group as location_parent_group,
       l.parent_value as location_parent_value,
       l.status       as location_status
FROM partner_management.v_field_executor vfe
         INNER JOIN partner_management.field_executor_location fel
                    ON vfe.id = fel.user_role_membership_id
         INNER JOIN partner_management.v_organization_executor_location voel
                    ON vfe.organization_id = voel.organization_id
                        AND voel.location_code = fel.location_code
         INNER JOIN master_data.location l
                    ON l.code = fel.location_code;

create table partner_management.field_executor_area
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    user_role_membership_id      bigint                   not null references user_management.user_role_membership (id),
    area_code                    varchar(50)              not null references master_data.area (code),
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (user_role_membership_id, area_code)
);

create view partner_management.v_field_executor_area as
select fea.*,
       a.type         as area_type,
       a.name         as area_name,
       a.value        as area_value,
       a.description  as area_description,
       a.parent_group as area_parent_group,
       a.parent_value as area_parent_value,
       a.status       as area_status,
       a.cost_center  as area_cost_center,
       a.is_deleted   as area_is_deleted,
       a.attribute1   as area_attribute1,
       a.attribute2   as area_attribute2,
       a.attribute3   as area_attribute3,
       a.attribute4   as area_attribute4,
       a.attribute5   as area_attribute5
from partner_management.field_executor_area fea
         left join partner_management.v_field_executor vfe on fea.user_role_membership_id = vfe.id
         left join master_data.area a on fea.area_code = a.code;

CREATE VIEW partner_management.v_field_executor_effective_area AS
SELECT vfe.*,
       vfe.id         AS user_role_membership_id,
       fea.area_code,
       a.type         as area_type,
       a.name         as area_name,
       a.value        as area_value,
       a.description  as area_description,
       a.parent_group as area_parent_group,
       a.parent_value as area_parent_value,
       a.status       as area_status
FROM partner_management.v_field_executor vfe
         INNER JOIN partner_management.field_executor_area fea
                    ON vfe.id = fea.user_role_membership_id
         INNER JOIN partner_management.v_organization_executor_area voea
                    ON vfe.organization_id = voea.organization_id
                        AND voea.area_code = fea.area_code
         INNER JOIN master_data.area a
                    ON a.code = fea.area_code;

CREATE VIEW partner_management.v_field_executor_effective_expertise AS
SELECT vfe.*,
       vfe.id         AS user_role_membership_id,
       fee.sub_task_type_id,
       vstt.task_type_id,
       vstt.task_type_code,
       vstt.task_type_name,
       vstt.code      AS sub_task_type_code,
       vstt.name      AS sub_task_type_name,
       vstt.full_name AS sub_task_type_full_name,
       vstt.full_code AS sub_task_type_full_code
FROM partner_management.v_field_executor vfe
         INNER JOIN partner_management.field_executor_expertise fee
                    ON vfe.id = fee.user_role_membership_id
         INNER JOIN partner_management.v_organization_executor_expertise voee
                    ON vfe.organization_id = voee.organization_id
                        AND voee.sub_task_type_id = fee.sub_task_type_id
         INNER JOIN task_management.v_sub_task_type vstt
                    ON vstt.id = fee.sub_task_type_id;

create table partner_management.special_flag
(
    id     bigserial primary key,
    uid    varchar(1024) not null unique default CONCAT(to_hex((extract(epoch from now()) * 1000000)::bigint),
                                                        gen_random_uuid()::text),
    nameid varchar       not null
);

insert into partner_management.special_flag (nameid)
values ('SK_PREMIER');

create table partner_management.field_executor_special_flag
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    user_role_membership_id      bigint                   not null references user_management.user_role_membership (id),
    special_flag_id              bigint                   not null references partner_management.special_flag (id),
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (user_role_membership_id, special_flag_id)
);

create view partner_management.v_organization_supervisor as
select tor.*,
       o.uid  as organization_uid,
       o.code as organization_code,
       o.name as organization_name,
       o.type as organization_type
from user_management.organization_role tor
         left join user_management.organization o on o.id = tor.organization_id
where tor.role_id = (select id from user_management.role where utag = 'FIELD_SUPERVISOR');

create table partner_management.organization_supervisor_expertise
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    organization_role_id         bigint                   not null references user_management.organization_role (id),
    sub_task_type_id             bigint                   not null references task_management.sub_task_type (id),
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (organization_role_id, sub_task_type_id)
);

create view partner_management.v_organization_supervisor_expertise as
select ose.*,
       vstt.task_type_id   as task_type_id,
       vstt.task_type_code as task_type_code,
       vstt.task_type_name as task_type_name,
       vstt.code           as sub_task_type_code,
       vstt.name           as sub_task_type_name,
       vstt.full_name      as sub_task_type_full_name,
       vstt.full_code      as sub_task_type_full_code,
       vos.organization_id,
       vos.organization_uid,
       vos.organization_code,
       vos.organization_name,
       vos.organization_type
from partner_management.organization_supervisor_expertise ose
         left join partner_management.v_organization_supervisor vos on ose.organization_role_id = vos.id
         left join task_management.v_sub_task_type vstt on vstt.id = ose.sub_task_type_id;


create table partner_management.organization_supervisor_area
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    organization_role_id         bigint                   not null references user_management.organization_role (id),
    area_code                    varchar(50)              not null references master_data.area (code),
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (organization_role_id, area_code)
);

create view partner_management.v_organization_supervisor_area as
select osa.*,
       a.type         as area_type,
       a.name         as area_name,
       a.value        as area_value,
       a.description  as area_description,
       a.cost_center  as area_cost_center,
       a.parent_group as area_parent_group,
       a.parent_value as area_parent_value,
       a.is_deleted   as area_is_deleted,
       a.status       as area_status,
       a.attribute1   as area_attribute1,
       a.attribute2   as area_attribute2,
       a.attribute3   as area_attribute3,
       a.attribute4   as area_attribute4,
       a.attribute5   as area_attribute5,
       vos.organization_id,
       vos.organization_uid,
       vos.organization_code,
       vos.organization_name,
       vos.organization_type
from partner_management.organization_supervisor_area osa
         left join partner_management.v_organization_supervisor vos on osa.organization_role_id = vos.id
         left join master_data.area a on osa.area_code = a.code;

create table partner_management.organization_supervisor_location
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    organization_role_id         bigint                   not null references user_management.organization_role (id),
    location_code                varchar(50)              not null references master_data.location (code),
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (organization_role_id, location_code)
);

create view partner_management.v_organization_supervisor_location as
select osl.*,
       l.type         as location_type,
       l.name         as location_name,
       l.value        as location_value,
       l.description  as location_description,
       l.parent_group as location_parent_group,
       l.parent_value as location_parent_value,
       l.status       as location_status,
       l.attribute1   as location_attribute1,
       l.attribute2   as location_attribute2,
       l.attribute3   as location_attribute3,
       l.attribute4   as location_attribute4,
       l.attribute5   as location_attribute5,
       vos.organization_id,
       vos.organization_uid,
       vos.organization_code,
       vos.organization_name,
       vos.organization_type
from partner_management.organization_supervisor_location osl
         left join partner_management.v_organization_supervisor vos on osl.organization_role_id = vos.id
         left join master_data.location l on osl.location_code = l.code;

create view partner_management.v_field_supervisor
as
select urm.*,
       u.uid                      as user_uid,
       u.loginid                  as user_loginid,
       u.email                    as user_email,
       u.fullname                 as user_fullname,
       u.phonenumber              as user_phonenumber,
       u.status                   as user_status,
       u.is_deleted               as user_is_deleted,
       u.identity_number          as user_identity_number,
       u.identity_type            as user_identity_type,
       u.gender                   as user_gender,
       u.address_on_identity_card as user_address_on_identity_card,
       o.uid                      as organization_uid,
       o.name                     as organization_name,
       o.type                     as organization_type,
       uom.membership_number      as organization_user_membership_number
from user_management.user_role_membership urm
         left join user_management.user u on urm.user_id = u.id
         left join user_management.user_organization_membership uom on u.id = uom.user_id
         left join user_management.organization o on urm.organization_id = o.id
where urm.role_id = (select id from user_management.role where utag = 'FIELD_SUPERVISOR');

create table partner_management.field_supervisor_area
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    user_role_membership_id      bigint                   not null references user_management.user_role_membership (id),
    area_code                    varchar(50)              not null references master_data.area (code),
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (user_role_membership_id, area_code)
);

create view partner_management.v_field_supervisor_area as
select fsa.*,
       a.type         as area_type,
       a.name         as area_name,
       a.value        as area_value,
       a.description  as area_description,
       a.parent_group as area_parent_group,
       a.parent_value as area_parent_value,
       a.status       as area_status,
       a.cost_center  as area_cost_center,
       a.is_deleted   as area_is_deleted,
       a.attribute1   as area_attribute1,
       a.attribute2   as area_attribute2,
       a.attribute3   as area_attribute3,
       a.attribute4   as area_attribute4,
       a.attribute5   as area_attribute5
from partner_management.field_supervisor_area fsa
         left join partner_management.v_field_supervisor vfs on fsa.user_role_membership_id = vfs.id
         left join master_data.area a on fsa.area_code = a.code;

CREATE VIEW partner_management.v_field_supervisor_effective_area AS
SELECT vfs.*,
       vfs.id         AS user_role_membership_id,
       fsa.area_code,
       a.type         as area_type,
       a.name         as area_name,
       a.value        as area_value,
       a.description  as area_description,
       a.parent_group as area_parent_group,
       a.parent_value as area_parent_value,
       a.status       as area_status
FROM partner_management.v_field_supervisor vfs
         INNER JOIN partner_management.field_supervisor_area fsa
                    ON vfs.id = fsa.user_role_membership_id
         INNER JOIN partner_management.v_organization_supervisor_area vosa
                    ON vfs.organization_id = vosa.organization_id
                        AND vosa.area_code = fsa.area_code
         INNER JOIN master_data.area a
                    ON a.code = fsa.area_code;

create table partner_management.field_supervisor_location
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    user_role_membership_id      bigint                   not null references user_management.user_role_membership (id),
    location_code                varchar(50)              not null references master_data.location (code),
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (user_role_membership_id, location_code)
);

create view partner_management.v_field_supervisor_location as
select fsl.*,
       l.type         as location_type,
       l.name         as location_name,
       l.value        as location_value,
       l.description  as location_description,
       l.parent_group as location_parent_group,
       l.parent_value as location_parent_value,
       l.status       as location_status,
       l.attribute1   as location_attribute1,
       l.attribute2   as location_attribute2,
       l.attribute3   as location_attribute3,
       l.attribute4   as location_attribute4,
       l.attribute5   as location_attribute5
from partner_management.field_supervisor_location fsl
         left join partner_management.v_field_supervisor vfs on fsl.user_role_membership_id = vfs.id
         left join master_data.location l on fsl.location_code = l.code;

CREATE VIEW partner_management.v_field_supervisor_effective_location AS
SELECT vfs.*,
       vfs.id         AS user_role_membership_id,
       fsl.location_code,
       l.type         as location_type,
       l.name         as location_name,
       l.value        as location_value,
       l.description  as location_description,
       l.parent_group as location_parent_group,
       l.parent_value as location_parent_value,
       l.status       as location_status
FROM partner_management.v_field_supervisor vfs
         INNER JOIN partner_management.field_supervisor_location fsl
                    ON vfs.id = fsl.user_role_membership_id
         INNER JOIN partner_management.v_organization_supervisor_location vosl
                    ON vfs.organization_id = vosl.organization_id
                        AND vosl.location_code = fsl.location_code
         INNER JOIN master_data.location l
                    ON l.code = fsl.location_code;

create table partner_management.field_supervisor_expertise
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    user_role_membership_id      bigint                   not null references user_management.user_role_membership (id),
    sub_task_type_id             bigint                   not null references task_management.sub_task_type (id),
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (user_role_membership_id, sub_task_type_id)
);

create view partner_management.v_field_supervisor_expertise as
select fse.*,
       vstt.task_type_id   as task_type_id,
       vstt.task_type_code as task_type_code,
       vstt.task_type_name as task_type_name,
       vstt.code           as sub_task_type_code,
       vstt.name           as sub_task_type_name,
       vstt.full_name      as sub_task_type_full_name,
       vstt.full_code      as sub_task_type_full_code
from partner_management.field_supervisor_expertise fse
         left join partner_management.v_field_supervisor vfs on fse.user_role_membership_id = vfs.id
         left join task_management.v_sub_task_type vstt on vstt.id = fse.sub_task_type_id;


CREATE VIEW partner_management.v_field_supervisor_effective_expertise AS
SELECT vfs.*,
       vfs.id         AS user_role_membership_id,
       fse.sub_task_type_id,
       vstt.task_type_id,
       vstt.task_type_code,
       vstt.task_type_name,
       vstt.code      AS sub_task_type_code,
       vstt.name      AS sub_task_type_name,
       vstt.full_name AS sub_task_type_full_name,
       vstt.full_code AS sub_task_type_full_code
FROM partner_management.v_field_supervisor vfs
         INNER JOIN partner_management.field_supervisor_expertise fse
                    ON vfs.id = fse.user_role_membership_id
         INNER JOIN partner_management.v_organization_supervisor_expertise vose
                    ON vfs.organization_id = vose.organization_id
                        AND vose.sub_task_type_id = fse.sub_task_type_id
         INNER JOIN task_management.v_sub_task_type vstt
                    ON vstt.id = fse.sub_task_type_id;

-- Executor Expertise
CREATE MATERIALIZED VIEW partner_management.mv_field_executor_effective_expertise AS
SELECT vfe.*,
       vfe.id         AS user_role_membership_id,
       fee.sub_task_type_id,
       vstt.task_type_id,
       vstt.task_type_code,
       vstt.task_type_name,
       vstt.code      AS sub_task_type_code,
       vstt.name      AS sub_task_type_name,
       vstt.full_name AS sub_task_type_full_name,
       vstt.full_code AS sub_task_type_full_code
FROM partner_management.v_field_executor vfe
         INNER JOIN partner_management.field_executor_expertise fee
                    ON vfe.id = fee.user_role_membership_id
         INNER JOIN partner_management.v_organization_executor_expertise voee
                    ON vfe.organization_id = voee.organization_id
                        AND voee.sub_task_type_id = fee.sub_task_type_id
         INNER JOIN task_management.v_sub_task_type vstt
                    ON vstt.id = fee.sub_task_type_id
WITH DATA;

-- Primary unique index for expertise
CREATE UNIQUE INDEX idx_mv_field_executor_effective_expertise_pk
    ON partner_management.mv_field_executor_effective_expertise (user_role_membership_id, sub_task_type_id);

-- Secondary indexes for common queries
CREATE INDEX idx_mv_field_executor_effective_expertise_uid
    ON partner_management.mv_field_executor_effective_expertise (user_uid);
CREATE INDEX idx_mv_field_executor_effective_expertise_org
    ON partner_management.mv_field_executor_effective_expertise (organization_id);

-- Executor Area
CREATE MATERIALIZED VIEW partner_management.mv_field_executor_effective_area AS
SELECT vfe.*,
       vfe.id         AS user_role_membership_id,
       fea.area_code,
       a.type         as area_type,
       a.name         as area_name,
       a.value        as area_value,
       a.description  as area_description,
       a.parent_group as area_parent_group,
       a.parent_value as area_parent_value,
       a.status       as area_status
FROM partner_management.v_field_executor vfe
         INNER JOIN partner_management.field_executor_area fea
                    ON vfe.id = fea.user_role_membership_id
         INNER JOIN partner_management.v_organization_executor_area voea
                    ON vfe.organization_id = voea.organization_id
                        AND voea.area_code = fea.area_code
         INNER JOIN master_data.area a
                    ON a.code = fea.area_code
WITH DATA;

-- Primary unique index for area
CREATE UNIQUE INDEX idx_mv_field_executor_effective_area_pk
    ON partner_management.mv_field_executor_effective_area (user_role_membership_id, area_code);

-- Secondary indexes for common queries
CREATE INDEX idx_mv_field_executor_effective_area_uid
    ON partner_management.mv_field_executor_effective_area (user_uid);
CREATE INDEX idx_mv_field_executor_effective_area_org
    ON partner_management.mv_field_executor_effective_area (organization_id);

-- Executor Location
CREATE MATERIALIZED VIEW partner_management.mv_field_executor_effective_location AS
SELECT vfe.*,
       vfe.id         AS user_role_membership_id,
       fel.location_code,
       l.type         as location_type,
       l.name         as location_name,
       l.value        as location_value,
       l.description  as location_description,
       l.parent_group as location_parent_group,
       l.parent_value as location_parent_value,
       l.status       as location_status
FROM partner_management.v_field_executor vfe
         INNER JOIN partner_management.field_executor_location fel
                    ON vfe.id = fel.user_role_membership_id
         INNER JOIN partner_management.v_organization_executor_location voel
                    ON vfe.organization_id = voel.organization_id
                        AND voel.location_code = fel.location_code
         INNER JOIN master_data.location l
                    ON l.code = fel.location_code
WITH DATA;

-- Primary unique index for location
CREATE UNIQUE INDEX idx_mv_field_executor_effective_location_pk
    ON partner_management.mv_field_executor_effective_location (user_role_membership_id, location_code);

-- Secondary indexes for common queries
CREATE INDEX idx_mv_field_executor_effective_location_uid
    ON partner_management.mv_field_executor_effective_location (user_uid);
CREATE INDEX idx_mv_field_executor_effective_location_org
    ON partner_management.mv_field_executor_effective_location (organization_id);

-- Supervisor Expertise
CREATE MATERIALIZED VIEW partner_management.mv_field_supervisor_effective_expertise AS
SELECT vfs.*,
       vfs.id         AS user_role_membership_id,
       fse.sub_task_type_id,
       vstt.task_type_id,
       vstt.task_type_code,
       vstt.task_type_name,
       vstt.code      AS sub_task_type_code,
       vstt.name      AS sub_task_type_name,
       vstt.full_name AS sub_task_type_full_name,
       vstt.full_code AS sub_task_type_full_code
FROM partner_management.v_field_supervisor vfs
         INNER JOIN partner_management.field_supervisor_expertise fse
                    ON vfs.id = fse.user_role_membership_id
         INNER JOIN partner_management.v_organization_supervisor_expertise vose
                    ON vfs.organization_id = vose.organization_id
                        AND vose.sub_task_type_id = fse.sub_task_type_id
         INNER JOIN task_management.v_sub_task_type vstt
                    ON vstt.id = fse.sub_task_type_id
WITH DATA;

-- Primary unique index for expertise
CREATE UNIQUE INDEX idx_mv_field_supervisor_effective_expertise_pk
    ON partner_management.mv_field_supervisor_effective_expertise (user_role_membership_id, sub_task_type_id);

-- Secondary indexes for common queries
CREATE INDEX idx_mv_field_supervisor_effective_expertise_uid
    ON partner_management.mv_field_supervisor_effective_expertise (user_uid);
CREATE INDEX idx_mv_field_supervisor_effective_expertise_org
    ON partner_management.mv_field_supervisor_effective_expertise (organization_id);

-- Create materialized view for supervisor area
CREATE MATERIALIZED VIEW partner_management.mv_field_supervisor_effective_area AS
SELECT vfs.*,
       vfs.id         AS user_role_membership_id,
       fsa.area_code,
       a.type         as area_type,
       a.name         as area_name,
       a.value        as area_value,
       a.description  as area_description,
       a.parent_group as area_parent_group,
       a.parent_value as area_parent_value,
       a.status       as area_status
FROM partner_management.v_field_supervisor vfs
         INNER JOIN partner_management.field_supervisor_area fsa
                    ON vfs.id = fsa.user_role_membership_id
         INNER JOIN partner_management.v_organization_supervisor_area vosa
                    ON vfs.organization_id = vosa.organization_id
                        AND vosa.area_code = fsa.area_code
         INNER JOIN master_data.area a
                    ON a.code = fsa.area_code
WITH DATA;

-- Create indexes for supervisor area
CREATE UNIQUE INDEX idx_mv_field_supervisor_effective_area_pk
    ON partner_management.mv_field_supervisor_effective_area (user_role_membership_id, area_code);
CREATE INDEX idx_mv_field_supervisor_effective_area_uid
    ON partner_management.mv_field_supervisor_effective_area (user_uid);
CREATE INDEX idx_mv_field_supervisor_effective_area_org
    ON partner_management.mv_field_supervisor_effective_area (organization_id);

-- Create materialized view for supervisor location
CREATE MATERIALIZED VIEW partner_management.mv_field_supervisor_effective_location AS
SELECT vfs.*,
       vfs.id         AS user_role_membership_id,
       fsl.location_code,
       l.type         as location_type,
       l.name         as location_name,
       l.value        as location_value,
       l.description  as location_description,
       l.parent_group as location_parent_group,
       l.parent_value as location_parent_value,
       l.status       as location_status
FROM partner_management.v_field_supervisor vfs
         INNER JOIN partner_management.field_supervisor_location fsl
                    ON vfs.id = fsl.user_role_membership_id
         INNER JOIN partner_management.v_organization_supervisor_location vosl
                    ON vfs.organization_id = vosl.organization_id
                        AND vosl.location_code = fsl.location_code
         INNER JOIN master_data.location l
                    ON l.code = fsl.location_code
WITH DATA;

-- Create indexes for supervisor location
CREATE UNIQUE INDEX idx_mv_field_supervisor_effective_location_pk
    ON partner_management.mv_field_supervisor_effective_location (user_role_membership_id, location_code);
CREATE INDEX idx_mv_field_supervisor_effective_location_uid
    ON partner_management.mv_field_supervisor_effective_location (user_uid);
CREATE INDEX idx_mv_field_supervisor_effective_location_org
    ON partner_management.mv_field_supervisor_effective_location (organization_id);

-- Create trigger functions for materialized view refresh
CREATE OR REPLACE FUNCTION partner_management.refresh_supervisor_area_mv()
    RETURNS TRIGGER AS
$$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY partner_management.mv_field_supervisor_effective_area;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION partner_management.refresh_supervisor_location_mv()
    RETURNS TRIGGER AS
$$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY partner_management.mv_field_supervisor_effective_location;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Add this trigger function
CREATE OR REPLACE FUNCTION partner_management.refresh_supervisor_expertise_mv()
    RETURNS TRIGGER AS
$$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY partner_management.mv_field_supervisor_effective_expertise;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Add these triggers
CREATE TRIGGER refresh_supervisor_expertise_mv_on_fse
    AFTER INSERT OR UPDATE OR DELETE
    ON partner_management.field_supervisor_expertise
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_supervisor_expertise_mv();

CREATE TRIGGER refresh_supervisor_expertise_mv_on_ose
    AFTER INSERT OR UPDATE OR DELETE
    ON partner_management.organization_supervisor_expertise
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_supervisor_expertise_mv();
-- Create triggers for Supervisor Area
CREATE TRIGGER refresh_supervisor_area_mv_on_fsa
    AFTER INSERT OR UPDATE OR DELETE
    ON partner_management.field_supervisor_area
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_supervisor_area_mv();

CREATE TRIGGER refresh_supervisor_area_mv_on_osa
    AFTER INSERT OR UPDATE OR DELETE
    ON partner_management.organization_supervisor_area
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_supervisor_area_mv();

-- Create triggers for Supervisor Location
CREATE TRIGGER refresh_supervisor_location_mv_on_fsl
    AFTER INSERT OR UPDATE OR DELETE
    ON partner_management.field_supervisor_location
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_supervisor_location_mv();

CREATE TRIGGER refresh_supervisor_location_mv_on_osl
    AFTER INSERT OR UPDATE OR DELETE
    ON partner_management.organization_supervisor_location
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_supervisor_location_mv();

-- Create trigger functions for executor materialized views
CREATE OR REPLACE FUNCTION partner_management.refresh_executor_expertise_mv()
    RETURNS TRIGGER AS
$$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY partner_management.mv_field_executor_effective_expertise;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION partner_management.refresh_executor_area_mv()
    RETURNS TRIGGER AS
$$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY partner_management.mv_field_executor_effective_area;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION partner_management.refresh_executor_location_mv()
    RETURNS TRIGGER AS
$$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY partner_management.mv_field_executor_effective_location;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for Executor Expertise
CREATE TRIGGER refresh_executor_expertise_mv_on_fee
    AFTER INSERT OR UPDATE OR DELETE
    ON partner_management.field_executor_expertise
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_executor_expertise_mv();

CREATE TRIGGER refresh_executor_expertise_mv_on_oee
    AFTER INSERT OR UPDATE OR DELETE
    ON partner_management.organization_executor_expertise
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_executor_expertise_mv();

-- Create triggers for Executor Area
CREATE TRIGGER refresh_executor_area_mv_on_fea
    AFTER INSERT OR UPDATE OR DELETE
    ON partner_management.field_executor_area
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_executor_area_mv();

CREATE TRIGGER refresh_executor_area_mv_on_oea
    AFTER INSERT OR UPDATE OR DELETE
    ON partner_management.organization_executor_area
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_executor_area_mv();

-- Create triggers for Executor Location
CREATE TRIGGER refresh_executor_location_mv_on_fel
    AFTER INSERT OR UPDATE OR DELETE
    ON partner_management.field_executor_location
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_executor_location_mv();

CREATE TRIGGER refresh_executor_location_mv_on_oel
    AFTER INSERT OR UPDATE OR DELETE
    ON partner_management.organization_executor_location
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_executor_location_mv();

-- Create materialized view
CREATE MATERIALIZED VIEW partner_management.mv_field_executor_sub_task_status_summary AS
WITH field_executor_user_sub_tasks AS (SELECT st.last_field_executor_user_id as user_id,
                                              st.sub_task_type_id,
                                              st.status,
                                              COUNT(*)                       as count
                                       FROM task_management.sub_task st
                                       WHERE st.is_deleted = false
                                       GROUP BY st.last_field_executor_user_id, st.sub_task_type_id, st.status)
SELECT vfe.user_id                                                                                 as user_id,
       COALESCE(SUM(CASE WHEN feust.sub_task_type_id = 1 THEN feust.count ELSE 0 END), 0)::integer as total_sub_task_construction_sk,
       COALESCE(SUM(CASE WHEN feust.sub_task_type_id = 2 THEN feust.count ELSE 0 END), 0)::integer as total_sub_task_construction_sr,
       COALESCE(SUM(CASE WHEN feust.sub_task_type_id = 3 THEN feust.count ELSE 0 END),
                0)::integer                                                                        as total_sub_task_construction_meter_installation,
       COALESCE(SUM(CASE WHEN feust.sub_task_type_id = 4 THEN feust.count ELSE 0 END),
                0)::integer                                                                        as total_sub_task_construction_gas_in,
       COALESCE(SUM(CASE
                        WHEN feust.status IN
                             ('WAITING_VERIFICATION', 'VERIFICATION_SUCCESS', 'CGP_VERIFICATION_SUCCESS')
                            THEN feust.count
                        ELSE 0 END),
                0)::integer                                                                        as total_sub_task_done,
       COALESCE(SUM(CASE WHEN feust.status IN ('ASSIGNED', 'WORKING', 'PAUSED') THEN feust.count ELSE 0 END),
                0)::integer                                                                        as total_sub_task_onprogress,
       COALESCE(SUM(CASE
                        WHEN feust.status IN ('VERIFICATION_FAIL', 'CGP_VERIFICATION_FAIL', 'FIXING') THEN feust.count
                        ELSE 0 END),
                0)::integer                                                                        as total_sub_task_revision,
       COALESCE(SUM(CASE
                        WHEN feust.status IN ('CANCELED_BY_CUSTOMER', 'CANCELED_BY_FIELD_EXECUTOR') THEN feust.count
                        ELSE 0 END),
                0)::integer                                                                        as total_sub_task_canceled,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = ANY (ARRAY[1::bigint, 2::bigint, 3::bigint, 4::bigint]) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_konstruksi,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = ANY (ARRAY[5::bigint, 6::bigint, 7::bigint, 8::bigint]) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_piutang,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = ANY (ARRAY[9::bigint, 10::bigint]) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_pengaduan_gangguan,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = ANY (ARRAY[9::bigint, 10::bigint]) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_layanan_teknis,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = 1 THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_konstruksi_sk,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = 2 THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_konstruksi_sr,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = 3 THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_konstruksi_meter_installation,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = 4 THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_konstruksi_gas_in,
       COALESCE(sum(
                        CASE
                            WHEN (feust.sub_task_type_id = ANY (ARRAY[1::bigint, 2::bigint, 3::bigint, 4::bigint])) AND (feust.status::text = ANY (ARRAY['WAITING_VERIFICATION'::character varying::text, 'VERIFICATION_SUCCESS'::character varying::text, 'CGP_VERIFICATION_SUCCESS'::character varying::text])) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_konstruksi_done,
       COALESCE(sum(
                        CASE
                            WHEN (feust.sub_task_type_id = ANY (ARRAY[1::bigint, 2::bigint, 3::bigint, 4::bigint])) AND (feust.status::text = ANY (ARRAY['ASSIGNED'::character varying::text, 'WORKING'::character varying::text, 'PAUSED'::character varying::text])) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_konstruksi_onprogress,
       COALESCE(sum(
                        CASE
                            WHEN (feust.sub_task_type_id = ANY (ARRAY[1::bigint, 2::bigint, 3::bigint, 4::bigint])) AND (feust.status::text = ANY (ARRAY['VERIFICATION_FAIL'::character varying::text, 'CGP_VERIFICATION_FAIL'::character varying::text, 'FIXING'::character varying::text])) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_konstruksi_revision,
       COALESCE(sum(
                        CASE
                            WHEN (feust.sub_task_type_id = ANY (ARRAY[1::bigint, 2::bigint, 3::bigint, 4::bigint])) AND (feust.status::text = ANY (ARRAY['CANCELED_BY_CUSTOMER'::character varying::text, 'CANCELED_BY_FIELD_EXECUTOR'::character varying::text])) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_konstruksi_canceled,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = 6 THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_piutang_cabut_meter_gas,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = 8 THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_piutang_pasang_meter_gas,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = 5 THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_piutang_pasang_segel,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = 7 THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_piutang_cabut_segel,
       COALESCE(sum(
                        CASE
                            WHEN (feust.sub_task_type_id = ANY (ARRAY[5::bigint, 6::bigint, 7::bigint, 8::bigint])) AND (feust.status::text = ANY (ARRAY['CANCELED_BY_CUSTOMER'::character varying::text, 'CANCELED_BY_EXECUTOR'::character varying::text, 'CANCELED_BY_FORCE_MAYOR'::character varying::text, 'CANCELED_BY_OTHER'::text, 'CANCELED_BY_EXPIRED'::text])) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_piutang_canceled,
       COALESCE(sum(
                        CASE
                            WHEN (feust.sub_task_type_id = ANY (ARRAY[5::bigint, 6::bigint, 7::bigint, 8::bigint])) AND (feust.status::text = ANY (ARRAY['PAID'::character varying::text, 'CANCELED_BY_CUSTOMER'::character varying::text])) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_piutang_canceled_paid,
       COALESCE(sum(
                        CASE
                            WHEN (feust.sub_task_type_id = ANY (ARRAY[5::bigint, 6::bigint, 7::bigint, 8::bigint])) AND (feust.status::text = ANY (ARRAY['DONE'::character varying::text])) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_piutang_done,
       COALESCE(sum(
                        CASE
                            WHEN (feust.sub_task_type_id = ANY (ARRAY[5::bigint, 6::bigint, 7::bigint, 8::bigint])) AND (feust.status::text = ANY (ARRAY['EXPIRED'::character varying::text])) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_piutang_expired,
       COALESCE(sum(
                        CASE
                            WHEN (feust.sub_task_type_id = ANY (ARRAY[5::bigint, 6::bigint, 7::bigint, 8::bigint])) AND (feust.status::text = ANY (ARRAY['WORKING'::character varying::text])) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_piutang_onprogress,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = 9 THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_pengaduan_gangguan_penyaluran,
       COALESCE(sum(
                        CASE
                            WHEN feust.sub_task_type_id = 10 THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_pengaduan_konstruksi_pipa,
       COALESCE(sum(
                        CASE
                            WHEN (feust.sub_task_type_id = ANY (ARRAY[9::bigint, 10::bigint])) AND (feust.status::text = ANY (ARRAY['CANCELED_BY_CUSTOMER'::character varying::text, 'CANCELED_BY_EXECUTOR'::character varying::text, 'CANCELED_BY_FORCE_MAYOR'::character varying::text, 'CANCELED_BY_OTHER'::text, 'CANCELED_BY_EXPIRED'::text])) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_pengaduan_canceled,
       COALESCE(sum(
                        CASE
                            WHEN (feust.sub_task_type_id = ANY (ARRAY[9::bigint, 10::bigint])) AND (feust.status::text = ANY (ARRAY['DONE'::character varying::text])) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_pengaduan_onprogress,
       COALESCE(sum(
                        CASE
                            WHEN (feust.sub_task_type_id = ANY (ARRAY[9::bigint, 10::bigint])) AND (feust.status::text = ANY (ARRAY['WORKING'::character varying::text])) THEN feust.count
                            ELSE 0::bigint
                            END), 0::numeric)::integer AS total_penanganan_pengaduan_done
FROM partner_management.v_field_executor vfe
         LEFT JOIN field_executor_user_sub_tasks feust ON vfe.user_id = feust.user_id
WHERE vfe.is_deleted = false
GROUP BY vfe.user_id;

CREATE UNIQUE INDEX idx_mv_field_executor_sub_task_status_summary_unique
    ON partner_management.mv_field_executor_sub_task_status_summary (user_id);

-- Create refresh function with error handling
CREATE OR REPLACE FUNCTION partner_management.refresh_field_executor_sub_task_status_summary()
    RETURNS TRIGGER AS
$$
BEGIN
    BEGIN
        REFRESH MATERIALIZED VIEW CONCURRENTLY partner_management.mv_field_executor_sub_task_status_summary;
    EXCEPTION
        WHEN OTHERS THEN
            RAISE WARNING 'Failed to refresh field_executor_sub_task_status_summary: %', SQLERRM;
    END;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
CREATE TRIGGER refresh_report_after_sub_task_change
    AFTER INSERT OR UPDATE OR DELETE
    ON task_management.sub_task
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_field_executor_sub_task_status_summary();

-- Create materialized view for supervisor verification statistics
CREATE MATERIALIZED VIEW partner_management.mv_field_supervisor_verification_stats AS
SELECT vfs.id        as user_role_membership_id,
       vfs.user_id,
       vfs.user_uid,
       vfs.user_loginid,
       vfs.user_email,
       vfs.user_fullname,
       vfs.user_phonenumber,
       vfs.organization_id,
       vfs.organization_uid,
       vfs.organization_name,
       vfs.organization_type,
       COALESCE(SUM(CASE WHEN str.sub_task_status = 'VERIFICATION_SUCCESS' THEN 1 ELSE 0 END),
                0)   as verification_success_count,
       COALESCE(SUM(CASE WHEN str.sub_task_status = 'VERIFICATION_FAIL' THEN 1 ELSE 0 END),
                0)   as verification_fail_count,
       COUNT(str.id) as total_verifications
FROM partner_management.v_field_supervisor vfs
         LEFT JOIN task_management.sub_task_report str ON str.user_id = vfs.user_id
    AND str.sub_task_status IN ('VERIFICATION_SUCCESS', 'VERIFICATION_FAIL')
    AND str.is_deleted = false
WHERE vfs.is_deleted = false
GROUP BY vfs.id,
         vfs.user_id,
         vfs.user_uid,
         vfs.user_loginid,
         vfs.user_email,
         vfs.user_fullname,
         vfs.user_phonenumber,
         vfs.organization_id,
         vfs.organization_uid,
         vfs.organization_name,
         vfs.organization_type;

-- Create index for better query performance
CREATE UNIQUE INDEX idx_supervisor_verification_stats_pk
    ON partner_management.mv_field_supervisor_verification_stats (user_role_membership_id);

CREATE INDEX idx_supervisor_verification_stats_org
    ON partner_management.mv_field_supervisor_verification_stats (organization_id);

-- Create refresh function
CREATE OR REPLACE FUNCTION partner_management.refresh_supervisor_verification_stats_mv()
    RETURNS TRIGGER AS
$$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY partner_management.mv_field_supervisor_verification_stats;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to refresh materialized view
CREATE TRIGGER refresh_supervisor_verification_stats_mv
    AFTER INSERT OR UPDATE OR DELETE
    ON task_management.sub_task_report
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_supervisor_verification_stats_mv();