DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_scheduler_enum') THEN
        CREATE TYPE status_scheduler_enum AS ENUM (
            'succeeded',
            'failed'
        );
    END IF;
END
$$;

ALTER TABLE billings
ADD COLUMN status_scheduler status_scheduler_enum DEFAULT NULL;