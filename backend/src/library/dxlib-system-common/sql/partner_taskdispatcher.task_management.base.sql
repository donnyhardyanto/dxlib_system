CREATE EXTENSION IF NOT EXISTS postgis;

create schema task_management;

create table task_management.task_type
(
    id                           bigint primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    code                         varchar(10) unique,
    name                         varchar(255) unique,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table task_management.sub_task_type
(
    id                           bigint primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    task_type_id                 bigint references task_management.task_type (id),
    code                         varchar(100),
    unique (task_type_id, code),
    name                         varchar(255),
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create view task_management.v_sub_task_type as
select stt.*,
       tt.name                          as task_type_name,
       tt.code                          as task_type_code,
       concat(tt.code, '-', stt.code)   as full_code,
       concat(tt.name, ' - ', stt.name) as full_name
from task_management.sub_task_type stt
         left join task_management.task_type tt on stt.task_type_id = tt.id;

create table task_management.customer
(
    id                              bigserial primary key,
    uid                             varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    registration_number             varchar(255) unique,                                      -- no ref pelanggan/no registration pelanggan
    customer_number                 varchar(255) unique,                                      -- id pelanggan/no pelanggan

    fullname                        varchar(255),-- IN FIGMA
    status                          varchar(255),                                             -- Status: "CALON_PELANGGAN, PELANGGAN, TIDAK_BERLANGGANAN_LAGI"
    special_flag_sk_primer          bool                                     default false,
    email                           varchar(255),-- IN FIGMA
    phonenumber                     varchar(255),-- IN FIGMA
    korespondensi_media             varchar(255),                                             -- email, Whatsapp, SMS, phone

    identity_type                   varchar(255),-- IN FIGMA: KTP
    identity_number                 varchar(255),-- IN FIGMA: KTP(NIK)
    npwp                            varchar(255),-- IN FIGMA

    customer_segment_code           varchar(255),                                             -- KI, PK, RT, Transmisi, Perusahaan, Individu --> master_data.vw_customer_segment
    customer_type_code              varchar(255),
    -- --> master_data.vw_customer_segmentBronze 1, Bronze 2, Bronze 3, Silver, Gold, Platinum, IJK1, IJK2, IJK3, IJKK, IJKU, IMP1, IMP2, IN FIGMA: bila customer_segment = RT,
    -- then possible value: RT1, RT2 CUSTOMER TYPE IN MASTER_DATA: GPiR
    jenis_anggaran                  varchar(255),                                             -- IN FIGMA: APBN
    rs_customer_sector_code         varchar(255) references master_data.rs_customer_sector (code),
    sales_area_code                 varchar(255) references master_data.area (code),          -- SOR

    latitude                        float,                                                    -- IN FIGMA
    longitude                       float,                                                    -- IN FIGMA
    geom                            geometry(Point, 4326),

    address_name                    varchar(255),                                             -- IN FIGMA
    address_street                  varchar(255),                                             -- IN FIGMA
    address_rt                      varchar(5),-- IN FIGMA
    address_rw                      varchar(5),-- IN FIGMA
    address_kelurahan_location_code varchar(255) references master_data.location (code),-- IN FIGMA
    address_kecamatan_location_code varchar(255) references master_data.location (code),-- IN FIGMA
    address_kabupaten_location_code varchar(255) references master_data.location (code),-- IN FIGMA
    address_province_location_code  varchar(255) references master_data.location (code),-- IN FIGMA

    address_postal_code             varchar(255),-- IN FIGMA

    register_at                     timestamp with time zone,

    jenis_bangunan                  varchar(255),

    program_pelanggan               varchar(255),                                             -- GasKita COCO 2022 SOR 2
    payment_scheme_code             varchar(255) references master_data.global_lookup (code), -- PASCABAYAR, PRABAYAR
    kategory_wilayah                varchar(255),                                             -- Non APBN

    cancellation_submission_status  int                                      default 0,
    -- 0: not canceled,
    -- 1: cancellation_submitted,
    -- 2: check data
    -- 3: publish Berita Acara Batal
    -- 4: Canceled Done

    is_deleted                      boolean                  not null        default false,
    created_at                      timestamp with time zone not null        default now(),
    created_by_user_id              varchar(255)             not null        default '',
    created_by_user_nameid          varchar(255)             not null        default '',
    last_modified_at                timestamp with time zone not null        default now(),
    last_modified_by_user_id        varchar(255)             not null        default '',
    last_modified_by_user_nameid    varchar(255)             not null        default ''
);

CREATE INDEX idx_customer_geom ON task_management.customer USING GIST (geom);

CREATE OR REPLACE FUNCTION task_management.update_customer_geom()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.geom = ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude), 4326);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to call the function
CREATE TRIGGER trigger_update_customer_geom
    BEFORE INSERT OR UPDATE OF latitude, longitude
    ON task_management.customer
    FOR EACH ROW
EXECUTE FUNCTION task_management.update_customer_geom();

create table task_management.customer_meter
(
    id                       bigserial primary key,
    uid                      varchar(1024) not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    register_timestamp       timestamp with time zone,
    customer_id              bigint        not null references task_management.customer (id),
    meter_appliance_type_id  bigint, -- known also as "meter_id" in API parameter
    meter_brand              varchar(255),
    sn_meter                 varchar(255),
    g_size_id                bigint,
    qmin                     double precision,
    qmax                     double precision,
    start_calibration_month  int,
    start_calibration_year   int,
    gas_in_date              date,
    meter_location_longitude double precision,
    meter_location_latitude  double precision
);

create table task_management.task
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    code                         varchar(255) unique,
    task_type_id                 bigint                   not null references task_management.task_type (id),
    status                       varchar(255), -- NOT_STARTED, STARTED, EXECUTION_DONE, REVISED, DONE, CANCELED, CANCELED_ACK_SALES, CANCELED_CLOSING_STATEMENT_SENT, CANCELLED_DONE
    customer_id                  bigint                   not null references task_management.customer (id),
    spk_no                       varchar(255),
    data1                        varchar(255),
    data2                        varchar(255),
    last_relyon_sync_at          timestamp with time zone,
    last_relyon_sync_status_code int,
    last_relyon_sync_message     varchar(255),
    last_relyon_sync_success_at  timestamp with time zone,
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table task_management.sub_task
(
    id                                                  bigserial primary key,
    uid                                                 varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    task_id                                             bigint                   not null references task_management.task (id),
    sub_task_type_id                                    bigint                   not null references task_management.sub_task_type (id),
    code                                                varchar(255) unique,
    status                                              varchar(255), -- WAITING_ASSIGNMENT, ASSIGNED, SCHEDULED, ON PROGRESS, ON REVISION, DONE, CANCELED
    priority                                            integer,
    last_field_executor_user_id                         bigint references user_management.user (id),
    last_field_executor_user_uid                        varchar(1024) references user_management.user (uid),
    last_field_executor_user_loginid                    varchar(255),
    last_field_executor_user_fullname                   varchar(255),

    last_field_supervisor_user_id                       bigint references user_management.user (id),
    last_field_supervisor_user_uid                      varchar(1024) references user_management.user (uid),
    last_field_supervisor_user_loginid                  varchar(255),
    last_field_supervisor_user_fullname                 varchar(255),

    last_cgp_user_id                                    bigint references user_management.user (id),
    last_cgp_user_uid                                   varchar(1024) references user_management.user (uid),
    last_cgp_user_loginid                               varchar(255),
    last_cgp_user_fullname                              varchar(255),

    last_sub_task_report_id                             bigint,       -- soft references task_management.sub_task_report (id), so it not circular reference
    last_sub_task_report_uid                            varchar(1024),

    last_form_sub_task_report_id                        bigint,       -- on last WAITING_VERIFICATION - soft references task_management.sub_task_report (id), so it not circular reference
    last_form_sub_task_report_uid                       varchar(1024),

    is_working_finish                                   boolean                  not null        default false,
    is_verification_success                             boolean                  not null        default false,
    is_cgp_verification_success                         boolean                  not null        default false,

    assigned_at                                         timestamp with time zone,
    expired_at                                          timestamp with time zone,

    scheduled_start_date                                date,
    scheduled_end_date                                  date,

    completed_at                                        timestamp with time zone,

    working_start_at                                    timestamp with time zone,
    working_end_at                                      timestamp with time zone,
    last_working_end_sub_task_report_id                 bigint,       -- soft references task_management.sub_task_report (id), so it not circular reference
    last_working_end_sub_task_report_uid                varchar(255),

    last_verification_end_at                            timestamp with time zone,
    last_verification_sub_task_report_id                bigint,       -- soft references task_management.sub_task_report (id), so it not circular reference
    last_verification_sub_task_report_uid               varchar(1024),

    last_cgp_verification_end_at                        timestamp with time zone,
    last_cgp_verification_sub_task_report_id            bigint,       -- soft references task_management.sub_task_report (id), so it not circular reference
    last_cgp_verification_sub_task_report_uid           varchar(1024),

    fix_count                                           integer                  not null        default 0,
    first_fixing_start_at                               timestamp with time zone,
    last_fixing_end_at                                  timestamp with time zone,
    last_fixing_end_sub_task_report_id                  bigint,       -- soft references task_management.sub_task_report (id), so it not circular reference
    last_fixing_end_sub_task_report_uid                 varchar(1024),

    last_reworking_end_at                               timestamp with time zone,
    last_reworking_end_sub_task_report_id               bigint,       -- soft references task_management.sub_task_report (id), so it not circular reference
    last_reworking_end_sub_task_report_uid              varchar(1024),

    last_start_pause_at                                 timestamp with time zone,
    last_end_pause_at                                   timestamp with time zone,
    last_pause_sub_task_report_id                       bigint,       -- soft references task_management.sub_task_report (id), so it not circular reference
    last_pause_sub_task_report_uid                      varchar(1024),

    last_canceled_by_field_executor_at                  timestamp with time zone,
    last_canceled_by_field_executor_sub_task_report_id  bigint,       -- soft references task_management.sub_task_report (id), so it not circular reference
    last_canceled_by_field_executor_sub_task_report_uid varchar(1024),

    last_canceled_by_customer_at                        timestamp with time zone,
    last_canceled_by_customer_sub_task_report_id        bigint,       -- soft references task_management.sub_task_report (id), so it not circular reference
    last_canceled_by_customer_sub_task_report_uid       varchar(1024),

    is_deleted                                          boolean                  not null        default false,
    created_at                                          timestamp with time zone not null        default now(),
    created_by_user_id                                  varchar(255)             not null        default '',
    created_by_user_nameid                              varchar(255)             not null        default '',
    last_modified_at                                    timestamp with time zone not null        default now(),
    last_modified_by_user_id                            varchar(255)             not null        default '',
    last_modified_by_user_nameid                        varchar(255)             not null        default ''
);

create table task_management.sub_task_report
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    sub_task_id                  bigint                   not null references task_management.sub_task (id),
    sub_task_uid                 varchar(255)             not null references task_management.sub_task (uid),
    sub_task_status              varchar(255),                      -- WAITING_ASSIGNMENT, ASSIGNED, SCHEDULED, ON PROGRESS, ON REVISION, DONE, CANCEL
    code                         varchar(255) unique,
    timestamp                    timestamp with time zone,
    user_id                      bigint                   not null, -- Yang membuat report, bisa Admin, Pelaksana Lapangan atau Supervisor
    user_uid                     varchar(1024)            not null, -- Yang membuat report, bisa Admin, Pelaksana Lapangan atau Supervisor
    user_loginid                 varchar(255)             not null,
    user_fullname                varchar(255)             not null,
    user_phonenumber             varchar(255),
    organization_id              bigint                   not null,
    organization_uid             varchar(1024)            not null,
    organization_name            varchar(255)             not null,
    report                       jsonb,
    berita_acara_link            varchar(1024),
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

create table task_management.sub_task_report_file_group
(
    id                           bigserial primary key,
    uid                          varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    sub_task_type_id             bigint                   not null references task_management.sub_task_type (id),
    nameid                       varchar(255),
    /*
    * 1. SK - Pneumatik Akhir, Isometrik, Valve, Pneumatik Awal, Bubble Test, Other, Ttd Pelanggan, Ttd Petugas
    * 2. SR -  Pneumatik Akhir, Isometrik, Tapping Saddle, Pneumatik Awal, Crossing Jalan, Crossing Saluran, Pondasi Sambungan Rumah, Other
    * 3. Meter Gas - Meter Barcode, Meter, Meter dan Tiang, Other
    * 4. Gas In - Valve, Stand Meter, Kompor/Peralatan, Selang ke Kompor, Tampak Depan Rumah, Stiker Gas In, Meter Other, Regulator, Other Regulator, Ttd Petugas, Ttd Pelanggan
     */
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default '',
    unique (sub_task_type_id, nameid)
);

create table task_management.sub_task_file
(
    id                             bigserial primary key,
    uid                            varchar(1024)            not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    sub_task_id                    bigint                   not null references task_management.sub_task (id),
    sub_task_uid                   varchar(1024)            not null references task_management.sub_task (uid),
    sub_task_report_file_group_id  bigint                   not null references task_management.sub_task_report_file_group (id),
    sub_task_report_file_group_uid varchar(1024)            not null references task_management.sub_task_report_file_group (uid),
    at                             timestamp with time zone,
    longitude                      float,
    latitude                       float,
    is_deleted                     boolean                  not null        default false,
    created_at                     timestamp with time zone not null        default now(),
    created_by_user_id             varchar(255)             not null        default '',
    created_by_user_nameid         varchar(255)             not null        default '',
    last_modified_at               timestamp with time zone not null        default now(),
    last_modified_by_user_id       varchar(255)             not null        default '',
    last_modified_by_user_nameid   varchar(255)             not null        default ''
);

create table task_management.sub_task_report_file
(
    id                  bigserial primary key,
    uid                 varchar(1024) not null unique default CONCAT(
            to_hex((extract(epoch from now()) * 1000000)::bigint), gen_random_uuid()::text),
    sub_task_report_id  bigint        not null references task_management.sub_task_report (id),
    sub_task_report_uid varchar(1024) not null references task_management.sub_task_report (uid),
    sub_task_file_id    bigint        not null references task_management.sub_task_file (id),
    sub_task_file_uid   varchar(1024) not null references task_management.sub_task_file (uid),
    unique (sub_task_report_id, sub_task_file_id)
);