CREATE EXTENSION IF NOT EXISTS postgis;

create view task_management.v_customer as
select pel.*,
       l.name  as address_kelurahan_location_name,
       c.name  as address_kecamatan_location_name,
       b.name  as address_kabupaten_location_name,
       p.name  as address_province_location_name,
       concat(
               CASE
                   WHEN pel.address_name IS NOT NULL AND pel.address_name <> ''
                       THEN pel.address_name || ', '
                   ELSE ''
                   END,
               pel.address_street, ',',
               'RT', pel.address_rt, 'RW', pel.address_rw, ',',
               l.name, ',',
               c.name, ',',
               b.name, ',',
               p.name
       )       as address,
       a.name  as sales_area_name,
       rs.name as rs_customer_sector_name
from task_management.customer pel
         left join master_data.location l on pel.address_kelurahan_location_code = l.code
         left join master_data.location c on pel.address_kecamatan_location_code = c.code
         left join master_data.location b on pel.address_kabupaten_location_code = b.code
         left join master_data.location p on pel.address_province_location_code = p.code
         left join master_data.area a on pel.sales_area_code = a.code
         left join master_data.rs_customer_sector rs on pel.rs_customer_sector_code = rs.code;


create view task_management.v_task as
select t.*,
       tt.name                           as task_type_name,
       tt.code                           as task_type_code,
       c.uid                             as customer_uid,
       c.special_flag_sk_primer          as customer_special_flag_sk_primer,
       c.registration_number             as customer_registration_number,
       c.customer_number                 as customer_number,
       c.fullname                        as customer_fullname,
       c.email                           as customer_email,
       c.phonenumber                     as customer_phonenumber,
       c.korespondensi_media             as customer_korespondensi_media,
       c.identity_type                   as customer_identity_type,
       c.identity_number                 as customer_identity_number,
       c.npwp                            as customer_npwp,
       c.customer_segment_code           as customer_segment_code,
       c.customer_type_code              as customer_type_code,
       c.jenis_anggaran                  as customer_jenis_anggaran,
       c.rs_customer_sector_code         as customer_rs_customer_sector_code,
       c.sales_area_code                 as customer_sales_area_code,
       c.latitude                        as customer_latitude,
       c.longitude                       as customer_longitude,
       c.geom                            as customer_geom,
       c.address                         as customer_address,
       c.address_name                    as customer_address_name,
       c.address_street                  as customer_address_street,
       c.address_rt                      as customer_address_rt,
       c.address_rw                      as customer_address_rw,
       c.address_kelurahan_location_code as customer_address_kelurahan_location_code,
       c.address_kecamatan_location_code as customer_address_kecamatan_location_code,
       c.address_kabupaten_location_code as customer_address_kabupaten_location_code,
       c.address_province_location_code  as customer_address_province_location_code,
       c.address_postal_code             as customer_address_postal_code,
       c.address_kelurahan_location_name as customer_address_kelurahan_location_name,
       c.address_kecamatan_location_name as customer_address_kecamatan_location_name,
       c.address_kabupaten_location_name as customer_address_kabupaten_location_name,
       c.address_province_location_name  as customer_address_province_location_name,
       c.register_at                     as customer_register_at,
       c.jenis_bangunan                  as customer_jenis_bangunan,
       c.program_pelanggan               as customer_program_pelanggan,
       c.payment_scheme_code             as customer_payment_scheme_code,
       c.kategory_wilayah                as customer_kategory_wilayah,
       c.is_deleted                      as customer_is_deleted
from task_management.task t
         left join task_management.v_customer c on t.customer_id = c.id
         left join task_management.task_type tt on t.task_type_id = tt.id;

create view task_management.v_sub_task as
select st.*,
       tt.uid                                  as task_type_uid,
       tt.id                                   as task_type_id,
       tt.code                                 as task_type_code,
       tt.name                                 as task_type_name,
       vstt.name                               as sub_task_type_name,
       vstt.full_code                          as sub_task_type_full_code,
       vstt.full_name                          as sub_task_type_full_name,
       t.uid                                   as task_uid,
       t.code                                  as task_code,
       concat(
               t.code, '-', st.code
       )                                       as full_code,
       c.id                                    as customer_id,
       c.uid                                   as customer_uid,
       c.registration_number                   as customer_registration_number,
       c.customer_number                       as customer_number,
       c.fullname                              as customer_fullname,
       c.special_flag_sk_primer                as customer_special_flag_sk_primer,
       c.email                                 as customer_email,
       c.phonenumber                           as customer_phonenumber,
       c.korespondensi_media                   as customer_korespondensi_media,
       c.identity_type                         as customer_identity_type,
       c.identity_number                       as customer_identity_number,
       c.npwp                                  as customer_npwp,
       c.customer_segment_code                 as customer_segment_code,
       c.customer_type_code                    as customer_type_code,
       c.jenis_anggaran                        as customer_jenis_anggaran,
       c.rs_customer_sector_code               as customer_rs_customer_sector_code,
       c.rs_customer_sector_name               as customer_rs_customer_sector_name,
       c.sales_area_code                       as customer_sales_area_code,
       c.sales_area_name                       as customer_sales_area_name,
       c.latitude                              as customer_latitude,
       c.longitude                             as customer_longitude,
       c.geom                                  as customer_geom,
       c.address_name                          as customer_address_name,
       c.address_street                        as customer_address_street,
       c.address_rt                            as customer_address_rt,
       c.address_rw                            as customer_address_rw,
       c.address_kelurahan_location_code       as customer_address_kelurahan_location_code,
       c.address_kecamatan_location_code       as customer_address_kecamatan_location_code,
       c.address_kabupaten_location_code       as customer_address_kabupaten_location_code,
       c.address_province_location_code        as customer_address_province_location_code,
       c.address_postal_code                   as customer_address_postal_code,
       c.register_at                           as customer_register_at,
       c.jenis_bangunan                        as customer_jenis_bangunan,
       c.program_pelanggan                     as customer_program_pelanggan,
       c.payment_scheme_code                   as customer_payment_scheme_code,
       c.kategory_wilayah                      as customer_kategory_wilayah,
       c.is_deleted                            as customer_is_deleted,
       vfe.user_phonenumber                    as last_field_executor_user_phonenumber,
       vfe.user_identity_number                as last_field_executor_user_identity_number,
       vfe.user_identity_type                  as last_field_executor_user_identity_type,
       vfe.organization_name                   as last_field_executor_user_organization_name,
       vfe.organization_type                   as last_field_executor_user_organization_type,
       vfe.organization_user_membership_number as last_field_executor_user_organization_user_membership_number,
       vfs.organization_type                   as last_field_supervisor_user_organization_type
from task_management.sub_task st
         left join task_management.task t on st.task_id = t.id
         left join task_management.task_type tt on t.task_type_id = tt.id
         left join task_management.v_sub_task_type vstt on st.sub_task_type_id = vstt.id
         left join task_management.v_customer c on t.customer_id = c.id
         left join partner_management.v_field_executor vfe on st.last_field_executor_user_id = vfe.user_id
         left join partner_management.v_field_supervisor vfs on st.last_field_supervisor_user_id = vfs.user_id;


create view task_management.v_sub_task_report as
select str.*,
       st.code                           as sub_task_code,
       st.last_field_executor_user_id    as sub_task_last_field_executor_user_id,
       st.last_field_executor_user_uid   as sub_task_last_field_executor_user_uid,
       st.last_field_supervisor_user_id  as sub_task_last_field_supervisor_user_id,
       st.last_field_supervisor_user_uid as sub_task_last_field_supervisor_user_uid,
       st.sub_task_type_id               as sub_task_type_id,
       vstt.name                         as sub_task_type_name,
       vstt.full_code                    as sub_task_type_full_code,
       vstt.full_name                    as sub_task_type_full_name,
       vt.id                             as task_id,
       vt.uid                            as task_uid,
       vt.code                           as task_code,
       vt.task_type_id,
       vt.task_type_code,
       vt.task_type_name,
       vt.status                         as task_status,
       vt.customer_id,
       vt.customer_uid,
       vt.customer_fullname,
       vt.customer_email,
       vt.customer_phonenumber,
       vt.customer_number,
       vt.customer_registration_number,
       vt.customer_identity_type,
       vt.customer_identity_number,
       vt.customer_npwp,
       vt.customer_sales_area_code,
       vt.customer_address_name,
       vt.customer_address_street,
       vt.customer_address_rt,
       vt.customer_address_rw,
       vt.customer_address_kelurahan_location_code,
       vt.customer_address_kecamatan_location_code,
       vt.customer_address_kabupaten_location_code,
       vt.customer_address_province_location_code,
       vt.customer_address_postal_code,
       vt.customer_address_kelurahan_location_name,
       vt.customer_address_kecamatan_location_name,
       vt.customer_address_kabupaten_location_name,
       vt.customer_address_province_location_name,
       vt.customer_address,
       vt.customer_longitude,
       vt.customer_latitude,
       vt.data1,
       vt.data2
from task_management.sub_task_report str
         left join task_management.sub_task st on str.sub_task_id = st.id
         left join task_management.v_sub_task_type vstt on st.sub_task_type_id = vstt.id
         left join task_management.v_task vt on st.task_id = vt.id;


create view task_management.v_sub_task_file as
select stf.*,
       strfg.nameid as sub_task_report_file_group_nameid,
       strfg.sub_task_type_id,
       vst.sub_task_type_name,
       vst.sub_task_type_full_code,
       vst.sub_task_type_full_name,
       vst.code     as sub_task_code,
       vst.task_id,
       vst.task_uid,
       vst.task_code,
       vst.task_type_id,
       vst.task_type_name,
       vst.task_type_code,
       vst.customer_id
from task_management.sub_task_file stf
         left join task_management.sub_task_report_file_group strfg on stf.sub_task_report_file_group_id = strfg.id
         left join task_management.v_sub_task vst on stf.sub_task_id = vst.id;

create view task_management.v_sub_task_report_file as
select strf.*,
       vstf.sub_task_report_file_group_id,
       vstf.sub_task_report_file_group_nameid,
       vstr.sub_task_id,
       vstr.sub_task_uid,
       vstr.sub_task_type_id,
       vstr.sub_task_type_name,
       vstr.sub_task_type_full_code,
       vstr.sub_task_type_full_name,
       vstr.sub_task_code,
       vstr.task_uid,
       vstr.task_id,
       vstr.task_code,
       vstr.task_type_id,
       vstr.task_type_name,
       vstr.task_type_code,
       vstr.customer_id
from task_management.sub_task_report_file strf
         left join task_management.v_sub_task_report vstr on strf.sub_task_report_id = vstr.id
         left join task_management.v_sub_task_file vstf on strf.sub_task_file_id = vstf.id;

create table task_management.sub_task_history_item
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    sub_task_id                  bigint                   not null references task_management.sub_task (id),
    sub_task_uid                 varchar(1024)            not null references task_management.sub_task (uid),
    timestamp                    timestamp with time zone,
    from_status                  varchar(255), -- WAITING_ASSIGNMENT, ASSIGNED, SCHEDULED, ON PROGRESS, ON REVISION, DONE, CANCEL
    to_status                    varchar(255), -- WAITING_ASSIGNMENT, ASSIGNED, SCHEDULED, ON PROGRESS, ON REVISION, DONE, CANCEL
    user_id                      bigint,
    user_uid                     varchar(1024),
    user_loginid                 varchar(255),
    user_fullname                varchar(255),
    organization_id              bigint,
    organization_uid             varchar(1024),
    organization_code            varchar(255),
    organization_name            varchar(255),
    organization_type            varchar(255),
    operation_nameid             varchar(255),
    sub_task_report_id           bigint references task_management.sub_task_report (id),
    sub_task_report_uid          varchar(1024) references task_management.sub_task_report (uid),
    input_parameters             jsonb,
    output_parameters            jsonb,

    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table task_management.task_history_item
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    task_id                      bigint                   not null references task_management.task (id),
    timestamp                    timestamp with time zone,
    from_status                  varchar(255), -- WAITING_ASSIGNMENT, ASSIGNED, SCHEDULED, ON PROGRESS, ON REVISION, DONE, CANCEL
    to_status                    varchar(255), -- WAITING_ASSIGNMENT, ASSIGNED, SCHEDULED, ON PROGRESS, ON REVISION, DONE, CANCEL
    user_id                      bigint,
    user_loginid                 varchar(255),
    user_fullname                varchar(255),
    organization_id              bigint,
    organization_code            varchar(255),
    organization_name            varchar(255),
    organization_type            varchar(255),
    operation_nameid             varchar(255),

    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);