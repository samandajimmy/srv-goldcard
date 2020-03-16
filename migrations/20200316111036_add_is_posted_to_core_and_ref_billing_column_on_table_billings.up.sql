DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'is_posted_to_core_enum') THEN
        CREATE TYPE is_posted_to_core_enum AS ENUM (
            'yes',
            'no'
        );
    END IF;
END
$$;

ALTER TABLE billings
ADD COLUMN is_posted_to_core is_posted_to_core_enum DEFAULT 'no';

ALTER TABLE billings
ADD COLUMN ref_billing VARCHAR(100) DEFAULT NULL