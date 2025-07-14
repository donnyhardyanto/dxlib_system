-- DROP SCHEMA upload_data;

CREATE SCHEMA upload_data;

-- upload_data.organization definition

-- Drop table

-- DROP TABLE upload_data.organization;

CREATE TABLE upload_data.organization
(
    id                           bigserial                                        NOT NULL,
    session_id                   varchar(255)                                     NULL,
    row_no                       int8                                             NULL,
    organization_id              int8                                             NULL,
    code                         varchar(255)                                     NULL,
    parent_code                  varchar(255)                                     NULL,
    "name"                       varchar(1024)                                    NULL,
    parent_id                    int8                                             NULL,
    "type"                       varchar(1024)                                    NULL,
    address                      varchar(1024)                                    NULL,
    npwp                         varchar(255)                                     NULL,
    email                        varchar(255)                                     NULL,
    phonenumber                  varchar(255)                                     NULL,
    status                       varchar(255) DEFAULT 'ACTIVE'::character varying NULL,
    auth_source1                 varchar(255)                                     NULL,
    auth_source2                 varchar(255)                                     NULL,
    attribute1                   varchar(1024)                                    NULL,
    attribute2                   varchar(1024)                                    NULL,
    utag                         varchar(255)                                     NULL,
    tags                         varchar(1024)                                    NULL,
    roles                        varchar(255)                                     NULL,
    field_executor_area          varchar(255)                                     NULL,
    field_executor_expertise     varchar(255)                                     NULL,
    field_executor_location      varchar(255)                                     NULL,
    field_supervisor_area        varchar(255)                                     NULL,
    field_supervisor_expertise   varchar(255)                                     NULL,
    field_supervisor_location    varchar(255)                                     NULL,
    row_status                   varchar(255)                                     NULL,
    row_message                  varchar                                          NULL,
    is_deleted                   bool         DEFAULT false                       NOT NULL,
    created_at                   timestamptz  DEFAULT now()                       NOT NULL,
    created_by_user_id           varchar(255) DEFAULT ''::character varying       NOT NULL,
    created_by_user_nameid       varchar(255) DEFAULT ''::character varying       NOT NULL,
    last_modified_at             timestamptz  DEFAULT now()                       NOT NULL,
    last_modified_by_user_id     varchar(255) DEFAULT ''::character varying       NOT NULL,
    last_modified_by_user_nameid varchar(255) DEFAULT ''::character varying       NOT NULL,
    CONSTRAINT organization_pkey PRIMARY KEY (id)
);

-- user
-- upload_data."user" definition

-- Drop table

-- DROP TABLE upload_data."user";

CREATE TABLE upload_data."user"
(
    id                           bigserial                                         NOT NULL,
    session_id                   varchar(255)                                      NULL,
    row_no                       int8                                              NULL,
    user_id                      int8                                              NULL,
    loginid                      varchar(255)                                      NOT NULL,
    email                        varchar(255)  DEFAULT ''::character varying       NULL,
    fullname                     varchar(255)  DEFAULT ''::character varying       NULL,
    phonenumber                  varchar(255)  DEFAULT ''::character varying       NULL,
    status                       varchar(255)  DEFAULT 'ACTIVE'::character varying NULL,
    "attribute"                  varchar(1024) DEFAULT ''::character varying       NULL,
    identity_number              varchar(255)                                      NULL,
    identity_type                varchar(255)  DEFAULT ''::character varying       NULL,
    gender                       varchar(1)                                        NULL,
    address_on_identity_card     varchar(1024)                                     NULL,
    utag                         varchar(255)                                      NULL,
    organization_code            varchar(255)                                      NULL,
    "role"                       varchar(255)                                      NULL,
    expertise                    varchar(255)                                      NULL,
    area_code                    varchar(255)                                      NULL,
    "location"                   varchar(255)                                      NULL,
    role_id                      int8                                              NULL,
    organization_id              int8                                              NULL,
    row_status                   varchar(10)                                       NULL,
    row_message                  varchar                                           NULL,
    is_deleted                   bool          DEFAULT false                       NOT NULL,
    created_at                   timestamptz   DEFAULT now()                       NOT NULL,
    created_by_user_id           varchar(255)  DEFAULT ''::character varying       NOT NULL,
    created_by_user_nameid       varchar(255)  DEFAULT ''::character varying       NOT NULL,
    last_modified_at             timestamptz   DEFAULT now()                       NOT NULL,
    last_modified_by_user_id     varchar(255)  DEFAULT ''::character varying       NOT NULL,
    last_modified_by_user_nameid varchar(255)  DEFAULT ''::character varying       NOT NULL,
    CONSTRAINT user_pkey PRIMARY KEY (id)
);

-- upload_data.customer definition

-- Drop table

-- DROP TABLE upload_data.customer;

CREATE TABLE upload_data.customer
(
    id                              bigserial                                  NOT NULL,
    session_id                      varchar(255)                               NULL,
    row_no                          int8                                       NULL,
    customer_id                     int8                                       NULL,
    registration_number             varchar(255)                               NULL,
    customer_number                 varchar(255)                               NULL,
    fullname                        varchar(255)                               NULL,
    status                          varchar(255)                               NULL,
    special_flag_sk_primer          bool         DEFAULT false                 NULL,
    email                           varchar(255)                               NULL,
    phonenumber                     varchar(255)                               NULL,
    korespondensi_media             varchar(255)                               NULL,
    identity_type                   varchar(255)                               NULL,
    identity_number                 varchar(255)                               NULL,
    npwp                            varchar(255)                               NULL,
    customer_segment_code           varchar(255)                               NULL,
    customer_type_code              varchar(255)                               NULL,
    jenis_anggaran                  varchar(255)                               NULL,
    rs_customer_sector_code         varchar(255)                               NULL,
    sales_area_code                 varchar(255)                               NULL,
    latitude                        float8                                     NULL,
    longitude                       float8                                     NULL,
    geom                            public.geometry(point, 4326)               NULL,
    address_name                    varchar(255)                               NULL,
    address_street                  varchar(255)                               NULL,
    address_rt                      varchar(5)                                 NULL,
    address_rw                      varchar(5)                                 NULL,
    address_kelurahan_location_code varchar(255)                               NULL,
    address_kecamatan_location_code varchar(255)                               NULL,
    address_kabupaten_location_code varchar(255)                               NULL,
    address_province_location_code  varchar(255)                               NULL,
    address_postal_code             varchar(255)                               NULL,
    register_at                     timestamptz                                NULL,
    jenis_bangunan                  varchar(255)                               NULL,
    program_pelanggan               varchar(255)                               NULL,
    payment_scheme_code             varchar(255)                               NULL,
    kategory_wilayah                varchar(255)                               NULL,
    cancellation_submission_status  int4         DEFAULT 0                     NULL,
    is_create_construction          bool         DEFAULT false                 NULL,
    row_status                      varchar(10)                                NULL,
    row_message                     varchar                                    NULL,
    is_deleted                      bool         DEFAULT false                 NOT NULL,
    created_at                      timestamptz  DEFAULT now()                 NOT NULL,
    created_by_user_id              varchar(255) DEFAULT ''::character varying NOT NULL,
    created_by_user_nameid          varchar(255) DEFAULT ''::character varying NOT NULL,
    last_modified_at                timestamptz  DEFAULT now()                 NOT NULL,
    last_modified_by_user_id        varchar(255) DEFAULT ''::character varying NOT NULL,
    last_modified_by_user_nameid    varchar(255) DEFAULT ''::character varying NOT NULL,
    CONSTRAINT customer_pkey PRIMARY KEY (id)
);

-- upload_data.arrears definition

-- Drop table

-- DROP TABLE upload_data.arrears;

CREATE TABLE upload_data.arrears
(
    id                           bigserial          NOT NULL,
    session_id                   varchar(256)       NULL, -- User session
    row_no                       int4               NULL, -- Row number in excel file
    nomor_surat                  varchar(255)       NULL,
    spk_no                       varchar(255)       NULL, -- Data from Excel
    id_pelanggan                 varchar(255)       NULL, -- Data from Excel
    nama_pelanggan               varchar(255)       NULL, -- Data from Excel
    kode_area                    varchar(255)       NULL, -- Data from Excel
    jenis_pelanggan              varchar(255)       NULL, -- Data from Excel
    area                         varchar(255)       NULL, -- Data from Excel
    kelompok_pelanggan           varchar            NULL, -- Data from Excel
    customer_management          varchar(255)       NULL, -- Data from Excel
    no_tel_customer_management   varchar(20)        NULL, -- Data from Excel
    no_telp_pelanggan            varchar(20)        NULL, -- Data from Excel
    alamat_pelanggan             varchar(256)       NULL, -- Data from Excel
    periode_awal_tunggakan       varchar(30)        NULL, -- Data from Excel
    periode_akhir_tunggakan      varchar(30)        NULL, -- Data from Excel
    bulan_tunggakan              varchar(256)       NULL, -- Data from Excel
    tagihan_pemakaian_gas        int8               NULL, -- Data from Excel
    jumlah_tagihan_rp            int8               NULL, -- Data from Excel
    denda                        int8               NULL, -- Data from Excel
    jaminan                      int8               NULL, -- Data from Excel
    biaya_pasang_kembali         int8               NULL, -- Data from Excel
    biaya_alir_kembali           int8               NULL, -- Data from Excel
    "action"                     varchar(30)        NULL, -- Data from Excel
    "no"                         int4               NULL, -- Data from Excel
    tgl_surat                    varchar(30)        NULL, -- Data from Excel
    customer_id                  int8               NULL, -- Mapping from ID_PELANGGAN column
    task_id                      int8               NULL, -- Reference after insert to task
    sub_task_id                  int8               NULL, -- Reference after insert to sub_task
    sub_task_type_id             int4               NULL, -- Mapping from ACTION column
    period_begin                 date               NULL, -- Mapping from periode akhir tunggakan
    period_end                   date               NULL,
    sales_area_code              varchar(255)       NULL, -- Mapping from area
    date_issued                  date               NULL, -- Mapping from tgl_surat
    row_status                   varchar(10)        NULL, -- Status data processing
    row_message                  varchar            NULL, -- Message data processing
    is_deleted                   bool DEFAULT false NULL,
    created_by_user_id           varchar(256)       NULL,
    created_by_user_nameid       varchar(256)       NULL,
    created_at                   timestamptz        NULL,
    last_modified_at             timestamptz        NULL,
    last_modified_by_user_id     varchar(256)       NULL,
    last_modified_by_user_nameid varchar(256)       NULL,
    CONSTRAINT upload_piutang_row_pk PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN upload_data.arrears.session_id IS 'User session';
COMMENT ON COLUMN upload_data.arrears.row_no IS 'Row number in excel file';
COMMENT ON COLUMN upload_data.arrears.spk_no IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.id_pelanggan IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.nama_pelanggan IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.kode_area IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.jenis_pelanggan IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.area IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.kelompok_pelanggan IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.customer_management IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.no_tel_customer_management IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.no_telp_pelanggan IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.alamat_pelanggan IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.periode_awal_tunggakan IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.periode_akhir_tunggakan IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.bulan_tunggakan IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.tagihan_pemakaian_gas IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.jumlah_tagihan_rp IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.denda IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.jaminan IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.biaya_pasang_kembali IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.biaya_alir_kembali IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears."action" IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears."no" IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.tgl_surat IS 'Data from Excel';
COMMENT ON COLUMN upload_data.arrears.customer_id IS 'Mapping from ID_PELANGGAN column';
COMMENT ON COLUMN upload_data.arrears.task_id IS 'Reference after insert to task';
COMMENT ON COLUMN upload_data.arrears.sub_task_id IS 'Reference after insert to sub_task';
COMMENT ON COLUMN upload_data.arrears.sub_task_type_id IS 'Mapping from ACTION column';
COMMENT ON COLUMN upload_data.arrears.period_begin IS 'Mapping from periode akhir tunggakan';
COMMENT ON COLUMN upload_data.arrears.sales_area_code IS 'Mapping from area';
COMMENT ON COLUMN upload_data.arrears.date_issued IS 'Mapping from tgl_surat';
COMMENT ON COLUMN upload_data.arrears.row_status IS 'Status data processing';
COMMENT ON COLUMN upload_data.arrears.row_message IS 'Message status data processing';