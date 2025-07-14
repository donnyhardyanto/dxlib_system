-- Create report schema if not exists
CREATE SCHEMA IF NOT EXISTS report;

-- Create the base view for task status monthly reporting
CREATE OR REPLACE VIEW report.v_task_status_monthly_report AS
WITH monthly_tasks AS (SELECT DATE_TRUNC('month', t.created_at)                                                   as report_month,
                              tt.code                                                                             as task_type_code,
                              tt.name                                                                             as task_type_name,
                              t.status                                                                            as task_status,
                              c.customer_segment_code,
                              c.customer_type_code,
                              c.rs_customer_sector_code,
                              c.sales_area_code,
                              COUNT(*)                                                                            as task_count,
                              COUNT(DISTINCT c.id)                                                                as unique_customer_count,
                              AVG(EXTRACT(EPOCH FROM (t.last_modified_at - t.created_at)) / 3600)::numeric(10, 2) as avg_task_duration_hours
                       FROM task_management.task t
                                JOIN task_management.task_type tt ON t.task_type_id = tt.id
                                JOIN task_management.customer c ON t.customer_id = c.id
                       WHERE t.is_deleted = false
                       GROUP BY DATE_TRUNC('month', t.created_at),
                                tt.code,
                                tt.name,
                                t.status,
                                c.customer_segment_code,
                                c.customer_type_code,
                                c.rs_customer_sector_code,
                                c.sales_area_code),
     subtask_metrics AS (SELECT DATE_TRUNC('month', st.created_at)                                                        as report_month,
                                t.task_type_id,
                                COUNT(*)                                                                                  as total_subtasks,
                                -- Status counts
                                SUM(CASE WHEN st.status = 'DONE' THEN 1 ELSE 0 END)                                       as completed_subtasks,
                                SUM(CASE WHEN st.status = 'CANCELED' THEN 1 ELSE 0 END)                                   as canceled_subtasks,
                                SUM(CASE WHEN st.status = 'ON_PROGRESS' THEN 1 ELSE 0 END)                                as in_progress_subtasks,
                                SUM(CASE WHEN st.status = 'ON_REVISION' THEN 1 ELSE 0 END)                                as revision_subtasks,
                                SUM(CASE WHEN st.status = 'WAITING_ASSIGNMENT' THEN 1 ELSE 0 END)                         as waiting_assignment_subtasks,
                                SUM(CASE WHEN st.status = 'ASSIGNED' THEN 1 ELSE 0 END)                                   as assigned_subtasks,
                                SUM(CASE WHEN st.status = 'SCHEDULED' THEN 1 ELSE 0 END)                                  as scheduled_subtasks,
                                -- Timing metrics
                                AVG(EXTRACT(EPOCH FROM (st.completed_at - st.created_at)) / 3600)::numeric(10, 2)         as avg_completion_hours,
                                AVG(EXTRACT(EPOCH FROM (st.working_end_at - st.working_start_at)) / 3600)::numeric(10, 2) as avg_working_hours,
                                -- Verification metrics
                                SUM(CASE WHEN st.is_verification_success = true THEN 1 ELSE 0 END)                        as verification_success_count,
                                SUM(CASE WHEN st.is_cgp_verification_success = true THEN 1 ELSE 0 END)                    as cgp_verification_success_count,
                                -- Fix/revision metrics
                                AVG(st.fix_count)::numeric(10, 2)                                                         as avg_fix_count,
                                AVG(EXTRACT(EPOCH FROM (st.last_fixing_end_at - st.first_fixing_start_at)) /
                                    3600)::numeric(10, 2)                                                                 as avg_fixing_duration_hours
                         FROM task_management.sub_task st
                                  JOIN task_management.task t ON st.task_id = t.id
                         WHERE st.is_deleted = false
                         GROUP BY DATE_TRUNC('month', st.created_at),
                                  t.task_type_id)

SELECT mt.report_month,
       mt.task_type_code,
       mt.task_type_name,
       mt.task_status,
       mt.customer_segment_code,
       mt.customer_type_code,
       mt.rs_customer_sector_code,
       mt.sales_area_code,
       mt.task_count,
       mt.unique_customer_count,
       mt.avg_task_duration_hours,
       -- Subtask counts
       COALESCE(sm.total_subtasks, 0)                 as total_subtasks,
       COALESCE(sm.completed_subtasks, 0)             as completed_subtasks,
       COALESCE(sm.canceled_subtasks, 0)              as canceled_subtasks,
       COALESCE(sm.in_progress_subtasks, 0)           as in_progress_subtasks,
       COALESCE(sm.revision_subtasks, 0)              as revision_subtasks,
       COALESCE(sm.waiting_assignment_subtasks, 0)    as waiting_assignment_subtasks,
       COALESCE(sm.assigned_subtasks, 0)              as assigned_subtasks,
       COALESCE(sm.scheduled_subtasks, 0)             as scheduled_subtasks,
       -- Timing metrics
       COALESCE(sm.avg_completion_hours, 0)           as avg_completion_hours,
       COALESCE(sm.avg_working_hours, 0)              as avg_working_hours,
       -- Verification metrics
       COALESCE(sm.verification_success_count, 0)     as verification_success_count,
       COALESCE(sm.cgp_verification_success_count, 0) as cgp_verification_success_count,
       -- Fix/revision metrics
       COALESCE(sm.avg_fix_count, 0)                  as avg_fix_count,
       COALESCE(sm.avg_fixing_duration_hours, 0)      as avg_fixing_duration_hours,
       -- Calculated rates
       CASE
           WHEN COALESCE(sm.total_subtasks, 0) > 0
               THEN (COALESCE(sm.completed_subtasks, 0)::float / sm.total_subtasks * 100)::numeric(5, 2)
           ELSE 0
           END                                        as completion_rate,
       CASE
           WHEN COALESCE(sm.completed_subtasks, 0) > 0
               THEN ((sm.verification_success_count + sm.cgp_verification_success_count)::float /
                     sm.completed_subtasks * 100)::numeric(5, 2)
           ELSE 0
           END                                        as verification_success_rate
FROM monthly_tasks mt
         JOIN task_management.task_type tt2 ON mt.task_type_code = tt2.code
         LEFT JOIN subtask_metrics sm ON
    mt.report_month = sm.report_month AND
    tt2.id = sm.task_type_id;

-- Create materialized view for better performance
CREATE MATERIALIZED VIEW report.mv_task_status_monthly_report AS
SELECT *
FROM report.v_task_status_monthly_report;

-- Create indexes for the materialized view
CREATE UNIQUE INDEX idx_mv_task_monthly_report_pk ON
    report.mv_task_status_monthly_report (
                                          report_month,
                                          task_type_code,
                                          task_status,
                                          customer_segment_code,
                                          customer_type_code,
                                          rs_customer_sector_code,
                                          sales_area_code
        );

CREATE INDEX idx_mv_task_monthly_report_date ON
    report.mv_task_status_monthly_report (report_month);

CREATE INDEX idx_mv_task_monthly_report_type ON
    report.mv_task_status_monthly_report (task_type_code);

CREATE INDEX idx_mv_task_monthly_report_status ON
    report.mv_task_status_monthly_report (task_status);

CREATE INDEX idx_mv_task_monthly_report_area ON
    report.mv_task_status_monthly_report (sales_area_code);

-- Create refresh function
CREATE OR REPLACE FUNCTION report.refresh_task_monthly_report_mv()
    RETURNS TRIGGER AS
$$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY report.mv_task_status_monthly_report;
    RETURN NULL;
EXCEPTION
    WHEN OTHERS THEN
        RAISE WARNING 'Failed to refresh task_status_monthly_report materialized view: %', SQLERRM;
        RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create triggers to refresh the materialized view
DROP TRIGGER IF EXISTS refresh_task_monthly_report_mv_on_task ON task_management.task;
CREATE TRIGGER refresh_task_monthly_report_mv_on_task
    AFTER INSERT OR UPDATE OR DELETE
    ON task_management.task
    FOR EACH STATEMENT
EXECUTE FUNCTION report.refresh_task_monthly_report_mv();

DROP TRIGGER IF EXISTS refresh_task_monthly_report_mv_on_subtask ON task_management.sub_task;
CREATE TRIGGER refresh_task_monthly_report_mv_on_subtask
    AFTER INSERT OR UPDATE OR DELETE
    ON task_management.sub_task
    FOR EACH STATEMENT
EXECUTE FUNCTION report.refresh_task_monthly_report_mv();

-- Create helper view for summary statistics
CREATE OR REPLACE VIEW report.v_task_monthly_summary AS
SELECT report_month,
       task_type_code,
       task_type_name,
       SUM(task_count)                as total_tasks,
       SUM(unique_customer_count)     as total_unique_customers,
       AVG(completion_rate)           as avg_completion_rate,
       AVG(verification_success_rate) as avg_verification_rate,
       AVG(avg_completion_hours)      as overall_avg_completion_hours,
       AVG(avg_fix_count)             as overall_avg_fix_count,
       SUM(completed_subtasks)        as total_completed_subtasks,
       SUM(canceled_subtasks)         as total_canceled_subtasks,
       SUM(in_progress_subtasks)      as total_in_progress_subtasks
FROM report.mv_task_status_monthly_report
GROUP BY report_month,
         task_type_code,
         task_type_name
ORDER BY report_month DESC,
         task_type_code;

-- Example queries:
COMMENT ON MATERIALIZED VIEW report.mv_task_status_monthly_report IS
    'Monthly task and subtask statistics including completion rates, verification success, and timing metrics.
    Example queries:

    -- Get monthly summary by task type:
    SELECT * FROM report.v_task_monthly_summary
    WHERE report_month >= NOW() - INTERVAL ''6 months'';

    -- Get detailed status breakdown for specific month:
    SELECT
        task_type_name,
        task_status,
        SUM(task_count) as count,
        AVG(completion_rate) as avg_completion_rate
    FROM report.mv_task_status_monthly_report
    WHERE report_month = DATE_TRUNC(''month'', NOW())
    GROUP BY task_type_name, task_status
    ORDER BY task_type_name, count DESC;

    -- Get verification success trends:
    SELECT
        report_month,
        task_type_name,
        SUM(verification_success_count) as total_verification_success,
        AVG(verification_success_rate) as avg_success_rate
    FROM report.mv_task_status_monthly_report
    GROUP BY report_month, task_type_name
    ORDER BY report_month DESC, task_type_name;';

-- report.v_customer_meter_smart source

CREATE OR REPLACE VIEW report.v_customer_meter_smart
AS SELECT vt.customer_address,
          (vt.customer_address_rt::text || '/'::text) || vt.customer_address_rw::text AS customer_rt_rw,
          vt.customer_address_postal_code,
          vt.customer_latitude,
          vt.customer_longitude,
          vt.customer_number,
          vt.customer_sales_area_code,
          ''::text AS customer_address_country,
          vt.customer_address_province_location_name,
          vt.customer_address_kabupaten_location_name,
          vt.customer_address_kecamatan_location_name,
          vt.customer_address_kelurahan_location_name,
          vt.code as task_code,
          cm.register_timestamp AS installation_date,
          cm.meter_brand,
          gs.name AS g_size,
          cm.qmin,
          cm.qmax,
          cm.start_calibration_year AS calibration_year,
          cm.gas_in_date,
          cm.sn_meter
   FROM task_management.v_task vt
            JOIN task_management.customer_meter cm ON vt.customer_id = cm.customer_id
            LEFT JOIN construction_management.g_size gs ON gs.id = cm.g_size_id
   WHERE cm.meter_appliance_type_id = 2;