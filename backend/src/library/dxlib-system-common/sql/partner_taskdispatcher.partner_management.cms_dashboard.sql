CREATE INDEX IF NOT EXISTS idx_task_status_date ON task_management.task (created_at, status, is_deleted);
CREATE INDEX IF NOT EXISTS idx_subtask_executor ON task_management.sub_task (last_field_executor_user_id, status, is_deleted);
CREATE INDEX IF NOT EXISTS idx_customer_location ON task_management.customer (is_deleted, address_kelurahan_location_code, address_kecamatan_location_code, address_kabupaten_location_code,
                                                                              address_province_location_code);

CREATE MATERIALIZED VIEW partner_management.mv_task_status_summary AS
SELECT task_type_id,
       COUNT(*)                                                as total_tasks,
       COUNT(*) FILTER (WHERE status = 'COMPLETED')            as completed_tasks,
       COUNT(*) FILTER (WHERE status = 'IN_PROGRESS')          as in_progress_tasks,
       COUNT(*) FILTER (WHERE status = 'WAITING_ASSIGNMENT')   as not_started_tasks,
       COUNT(*) FILTER (WHERE status = 'CANCELED_BY_CUSTOMER') as canceled_tasks
FROM task_management.task
WHERE is_deleted = false
GROUP BY task_type_id
ORDER BY task_type_id;

CREATE UNIQUE INDEX idx_mv_task_status_summary ON partner_management.mv_task_status_summary (task_type_id);

CREATE MATERIALIZED VIEW partner_management.mv_task_location_distribution AS
WITH location_data AS (SELECT t.id   as task_id,
                              k.name as kelurahan,
                              c.name as kecamatan,
                              b.name as kabupaten,
                              p.name as province
                       FROM task_management.task t
                                JOIN task_management.customer cust ON t.customer_id = cust.id
                                LEFT JOIN master_data.location k ON cust.address_kelurahan_location_code = k.code
                                LEFT JOIN master_data.location c ON cust.address_kecamatan_location_code = c.code
                                LEFT JOIN master_data.location b ON cust.address_kabupaten_location_code = b.code
                                LEFT JOIN master_data.location p ON cust.address_province_location_code = p.code
                       WHERE t.is_deleted = false),
     location_counts AS (SELECT COALESCE(province, '')  as province,
                                COALESCE(kabupaten, '') as kabupaten,
                                COALESCE(kecamatan, '') as kecamatan,
                                COALESCE(kelurahan, '') as kelurahan,
                                COUNT(*)                as task_count
                         FROM location_data
                         GROUP BY province, kabupaten, kecamatan, kelurahan)
SELECT province,
       kabupaten,
       kecamatan,
       kelurahan,
       task_count,
       ROUND(100.0 * task_count / NULLIF(SUM(task_count) OVER (), 0), 2)::float as percentage
FROM location_counts
ORDER BY task_count DESC;

CREATE UNIQUE INDEX idx_mv_task_location_distribution ON partner_management.mv_task_location_distribution (province, kabupaten, kecamatan, kelurahan);

CREATE MATERIALIZED VIEW partner_management.mv_task_time_series AS
WITH RECURSIVE
    dates AS (SELECT date_trunc('day', make_date(2024, 1, 1))::date as day
              UNION ALL
              SELECT (day + interval '1 day')::date
              FROM dates
              WHERE day < date_trunc('day', CURRENT_DATE)::date),
    daily_task_stats AS (SELECT date_trunc('day', created_at)::date                     as day,
                                COUNT(*)                                                as total_tasks,
                                COUNT(*) FILTER (WHERE status = 'COMPLETED')            as completed_tasks,
                                COUNT(*) FILTER (WHERE status = 'IN_PROGRESS')          as in_progress_tasks,
                                COUNT(*) FILTER (WHERE status = 'WAITING_ASSIGNMENT')   as not_started_tasks,
                                COUNT(*) FILTER (WHERE status = 'CANCELED_BY_CUSTOMER') as canceled_tasks
                         FROM task_management.task
                         WHERE is_deleted = false
                         GROUP BY 1)
SELECT d.day,
       COALESCE(ds.total_tasks, 0)::integer                                    as new_tasks,
       SUM(COALESCE(ds.total_tasks, 0)) OVER (ORDER BY d.day)::integer         as cumulative_total_tasks,
       COALESCE(ds.completed_tasks, 0)::integer                                as new_completed_tasks,
       SUM(COALESCE(ds.completed_tasks, 0)) OVER (ORDER BY d.day)::integer     as cumulative_completed_tasks,
       COALESCE(ds.in_progress_tasks, 0)::integer                              as new_in_progress_tasks,
       (SUM(COALESCE(ds.in_progress_tasks, 0)) OVER (ORDER BY d.day) -
        SUM(COALESCE(ds.completed_tasks, 0)) OVER (ORDER BY d.day))::integer   as current_in_progress_tasks,
       COALESCE(ds.not_started_tasks, 0)::integer                              as new_not_started_tasks,
       (SUM(COALESCE(ds.not_started_tasks, 0)) OVER (ORDER BY d.day) -
        SUM(COALESCE(ds.in_progress_tasks, 0)) OVER (ORDER BY d.day))::integer as current_not_started_tasks,
       COALESCE(ds.canceled_tasks, 0)::integer                                 as new_canceled_tasks,
       SUM(COALESCE(ds.canceled_tasks, 0)) OVER (ORDER BY d.day)::integer      as cumulative_canceled_tasks
FROM dates d
         LEFT JOIN daily_task_stats ds ON d.day = ds.day
ORDER BY d.day;

CREATE UNIQUE INDEX idx_mv_task_time_series_day ON partner_management.mv_task_time_series (day);

CREATE MATERIALIZED VIEW partner_management.mv_field_executor_status_time_series AS
WITH RECURSIVE
    dates AS (SELECT date_trunc('day', make_date(2024, 1, 1))::date as day
              UNION ALL
              SELECT (day + interval '1 day')::date
              FROM dates
              WHERE day < date_trunc('day', CURRENT_DATE)::date),
    active_executors AS (SELECT user_id, created_at, last_modified_at, is_deleted, user_status
                         FROM partner_management.v_field_executor
                         WHERE is_deleted = false),
    executor_tasks AS (SELECT DISTINCT last_field_executor_user_id, created_at, completed_at
                       FROM task_management.sub_task
                       WHERE status IN ('ASSIGNED', 'WORKING', 'FIXING', 'REWORKING')
                         AND is_deleted = false)
SELECT d.day,
       COUNT(DISTINCT fe.user_id) FILTER (WHERE fe.created_at <= d.day AND (fe.is_deleted = false OR fe.last_modified_at > d.day))::integer                                  as total_executors,
       COUNT(DISTINCT fe.user_id) FILTER (WHERE fe.created_at <= d.day AND (fe.is_deleted = false OR fe.last_modified_at > d.day) AND fe.user_status = 'ACTIVE')::integer    as active_executors,
       COUNT(DISTINCT fe.user_id) FILTER (WHERE fe.created_at <= d.day AND (fe.is_deleted = false OR fe.last_modified_at > d.day) AND fe.user_status = 'SUSPENDED')::integer as suspended_executors,
       COUNT(DISTINCT fe.user_id) FILTER (WHERE fe.created_at <= d.day AND (fe.is_deleted = false OR fe.last_modified_at > d.day) AND fe.user_status = 'DELETED')::integer   as deleted_executors,
       COUNT(DISTINCT et.last_field_executor_user_id) FILTER (WHERE et.created_at <= d.day AND (et.completed_at IS NULL OR et.completed_at > d.day))::integer                as executors_with_tasks,
       ROUND(100.0 * COUNT(DISTINCT fe.user_id) FILTER (WHERE fe.user_status = 'ACTIVE') / NULLIF(COUNT(DISTINCT fe.user_id), 0), 2)::float                                  as active_percentage,
       ROUND(100.0 * COUNT(DISTINCT fe.user_id) FILTER (WHERE fe.user_status = 'SUSPENDED') / NULLIF(COUNT(DISTINCT fe.user_id), 0), 2)::float                               as suspended_percentage,
       ROUND(100.0 * COUNT(DISTINCT fe.user_id) FILTER (WHERE fe.user_status = 'DELETED') / NULLIF(COUNT(DISTINCT fe.user_id), 0), 2)::float                                 as deleted_percentage,
       ROUND(100.0 * COUNT(DISTINCT et.last_field_executor_user_id) / NULLIF(COUNT(DISTINCT fe.user_id), 0), 2)::float                                                       as executors_with_tasks_percentage
FROM dates d
         CROSS JOIN active_executors fe
         LEFT JOIN executor_tasks et ON fe.user_id = et.last_field_executor_user_id
GROUP BY d.day
ORDER BY d.day;

CREATE UNIQUE INDEX idx_mv_field_executor_status_time_series_day ON partner_management.mv_field_executor_status_time_series (day);

CREATE MATERIALIZED VIEW partner_management.mv_dashboard_time_series AS
SELECT t.day::date,
       t.cumulative_total_tasks     as total_tasks,
       t.cumulative_completed_tasks as completed_tasks,
       t.current_in_progress_tasks  as in_progress_tasks,
       e.total_executors,
       e.active_executors,
       e.suspended_executors,
       e.deleted_executors,
       e.executors_with_tasks,
       e.active_percentage,
       e.suspended_percentage,
       e.deleted_percentage,
       e.executors_with_tasks_percentage
FROM partner_management.mv_task_time_series t
         JOIN partner_management.mv_field_executor_status_time_series e USING (day)
ORDER BY day;

CREATE UNIQUE INDEX idx_mv_dashboard_time_series ON partner_management.mv_dashboard_time_series (day);

CREATE OR REPLACE FUNCTION partner_management.refresh_dashboard_views() RETURNS TRIGGER AS
$$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY partner_management.mv_task_status_summary;
    REFRESH MATERIALIZED VIEW CONCURRENTLY partner_management.mv_task_location_distribution;
    REFRESH MATERIALIZED VIEW CONCURRENTLY partner_management.mv_task_time_series;
    REFRESH MATERIALIZED VIEW CONCURRENTLY partner_management.mv_field_executor_status_time_series;
    REFRESH MATERIALIZED VIEW CONCURRENTLY partner_management.mv_dashboard_time_series;
    RETURN NULL;
EXCEPTION
    WHEN OTHERS THEN
        RAISE WARNING 'Failed to refresh views: %', SQLERRM;
        RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER refresh_views_on_task_change
    AFTER INSERT OR UPDATE OR DELETE
    ON task_management.task
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_dashboard_views();

CREATE TRIGGER refresh_views_on_subtask_change
    AFTER INSERT OR UPDATE OR DELETE
    ON task_management.sub_task
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_dashboard_views();

CREATE TRIGGER refresh_views_on_customer_location_change
    AFTER UPDATE OF address_kelurahan_location_code,
        address_kecamatan_location_code,
        address_kabupaten_location_code,
        address_province_location_code
    ON task_management.customer
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_dashboard_views();

CREATE TRIGGER refresh_views_on_user_change
    AFTER UPDATE OF status
    ON user_management.user
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_dashboard_views();

CREATE TRIGGER refresh_views_on_role_change
    AFTER INSERT OR UPDATE OR DELETE
    ON user_management.user_role_membership
    FOR EACH STATEMENT
EXECUTE FUNCTION partner_management.refresh_dashboard_views();
