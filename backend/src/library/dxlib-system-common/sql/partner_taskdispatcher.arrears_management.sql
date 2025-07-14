
CREATE SCHEMA arrears_management AUTHORIZATION postgres;

-- Create customer_piutang table in task_management schema
create table arrears_management.task_arrears
(
    id                           bigserial primary key,
    task_id                      bigint                   not null references task_management.task (id),
    amount_usage_bill            bigint                   null,
    amount_fine                  bigint                   null,
    amount_payment_guarantee     bigint                   null,
    amount_reinstallation_cost   bigint                   null,
    amount_reflow_cost           bigint                   null,
    amount_bill_total            bigint                   null,
    period_begin                 date                     null,
    period_end                   date                     null,
    status                       varchar(255)             null, -- status piutang: LUNAS, MENUNGGAK
    is_deleted                   boolean                  not null        default false,
    created_at                   timestamp with time zone not null        default now(),
    created_by_user_id           varchar(255)             not null        default '',
    created_by_user_nameid       varchar(255)             not null        default '',
    last_modified_at             timestamp with time zone not null        default now(),
    last_modified_by_user_id     varchar(255)             not null        default '',
    last_modified_by_user_nameid varchar(255)             not null        default ''
);

-- Column comment

COMMENT ON COLUMN arrears_management.task_arrears.status IS 'Status Piutang:
- LUNAS
- MENUNGGAK';
COMMENT ON COLUMN arrears_management.task_arrears.amount_usage_bill IS 'User session';
COMMENT ON COLUMN arrears_management.task_arrears.amount_reinstallation_cost IS 'Biaya Pasang Kembali';
COMMENT ON COLUMN arrears_management.task_arrears.amount_reflow_cost IS 'Biaya Alir Kembali';
COMMENT ON COLUMN arrears_management.task_arrears.amount_bill_total IS 'Total Tagihan';

-- arrears_management.upload_piutang_row definition

-- Drop table

-- DROP TABLE arrears_management.upload_piutang_row;

CREATE TABLE arrears_management.upload_arrears_row (
   id bigserial NOT NULL,
   session_id varchar(256) NULL, -- User session
   row_no int4 NULL, -- Row number in excel file
   nomor_surat varchar(255),
   spk_no varchar(255) NULL, -- Data from Excel
   id_pelanggan varchar(255) NULL, -- Data from Excel
   nama_pelanggan varchar(255) NULL, -- Data from Excel
   kode_area varchar(255) NULL, -- Data from Excel
   jenis_pelanggan varchar(255) NULL, -- Data from Excel
   area varchar(255) NULL, -- Data from Excel
   kelompok_pelanggan varchar NULL, -- Data from Excel
   customer_management varchar(255) NULL, -- Data from Excel
   no_tel_customer_management varchar(20) NULL, -- Data from Excel
   no_telp_pelanggan varchar(20) NULL, -- Data from Excel
   alamat_pelanggan varchar(256) NULL, -- Data from Excel
   periode_awal_tunggakan varchar(30) NULL, -- Data from Excel
   periode_akhir_tunggakan varchar(30) NULL, -- Data from Excel
   bulan_tunggakan varchar(256) NULL, -- Data from Excel
   tagihan_pemakaian_gas int8 NULL, -- Data from Excel
   jumlah_tagihan_rp int8 NULL, -- Data from Excel
   denda int8 NULL, -- Data from Excel
   jaminan int8 NULL, -- Data from Excel
   biaya_pasang_kembali int8 NULL, -- Data from Excel
   biaya_alir_kembali int8 NULL, -- Data from Excel
   "action" varchar(30) NULL, -- Data from Excel
   "no" int4 NULL, -- Data from Excel
   tgl_surat varchar(30) NULL, -- Data from Excel
   customer_id int8 NULL, -- Mapping from ID_PELANGGAN column
   task_id int8 NULL, -- Reference after insert to task
   sub_task_id int8 NULL, -- Reference after insert to sub_task
   sub_task_type_id int4 NULL, -- Mapping from ACTION column
   period_begin date NULL, -- Mapping from periode awal tunggakan
   period_end date NULL,  -- Mapping from periode akhir tunggakan
   sales_area_code varchar(255), -- Mapping from area
   date_issued date NULL, -- Mapping from tgl_surat
   row_status varchar(10) NULL, -- Status data processing
   is_deleted bool NULL default false,
   created_by_user_id varchar(256) NULL,
   created_by_user_nameid varchar(256) NULL,
   created_at timestamptz NULL,
   last_modified_at timestamptz NULL,
   last_modified_by_user_id varchar(256) NULL,
   last_modified_by_user_nameid varchar(256) NULL,
   CONSTRAINT upload_piutang_row_pk PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN arrears_management.upload_arrears_row.session_id IS 'User session';
COMMENT ON COLUMN arrears_management.upload_arrears_row.row_no IS 'Row number in excel file';
COMMENT ON COLUMN arrears_management.upload_arrears_row.spk_no IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.id_pelanggan IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.nama_pelanggan IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.kode_area IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.jenis_pelanggan IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.area IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.kelompok_pelanggan IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.customer_management IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.no_tel_customer_management IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.no_telp_pelanggan IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.alamat_pelanggan IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.periode_awal_tunggakan IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.periode_akhir_tunggakan IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.bulan_tunggakan IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.tagihan_pemakaian_gas IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.jumlah_tagihan_rp IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.denda IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.jaminan IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.biaya_pasang_kembali IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.biaya_alir_kembali IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row."action" IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row."no" IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.tgl_surat IS 'Data from Excel';
COMMENT ON COLUMN arrears_management.upload_arrears_row.customer_id IS 'Mapping from ID_PELANGGAN column';
COMMENT ON COLUMN arrears_management.upload_arrears_row.task_id IS 'Reference after insert to task';
COMMENT ON COLUMN arrears_management.upload_arrears_row.sub_task_id IS 'Reference after insert to sub_task';
COMMENT ON COLUMN arrears_management.upload_arrears_row.sub_task_type_id IS 'Mapping from ACTION column';
COMMENT ON COLUMN arrears_management.upload_arrears_row.period_begin IS 'Mapping from periode awal tunggakan';
COMMENT ON COLUMN arrears_management.upload_arrears_row.period_begin IS 'Mapping from periode akhir tunggakan';
COMMENT ON COLUMN arrears_management.upload_arrears_row.sales_area_code IS 'Mapping from area';
COMMENT ON COLUMN arrears_management.upload_arrears_row.date_issued IS 'Mapping from tgl_surat';
COMMENT ON COLUMN arrears_management.upload_arrears_row.row_status IS 'Status data processing';
