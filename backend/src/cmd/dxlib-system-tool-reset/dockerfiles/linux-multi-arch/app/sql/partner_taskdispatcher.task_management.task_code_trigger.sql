-- Create function to generate task code from task ID
CREATE OR REPLACE FUNCTION task_management.generate_task_code()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.code IS NULL OR NEW.code = '' THEN
        NEW.code := 'TASK-' || LPAD(NEW.id::text, 8, '0');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically set task code
DROP TRIGGER IF EXISTS set_task_code ON task_management.task;
CREATE TRIGGER set_task_code
    BEFORE INSERT ON task_management.task
    FOR EACH ROW
    EXECUTE FUNCTION task_management.generate_task_code(); 