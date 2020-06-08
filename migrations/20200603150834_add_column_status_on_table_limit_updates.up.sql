DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_limit_update_enum') THEN
        CREATE TYPE status_limit_update_enum AS ENUM (
            'pending',
            'applied',
            'approved',
            'rejected'
        );
    END IF;
END
$$;

ALTER TABLE limit_updates ADD status status_limit_update_enum DEFAULT NULL;
