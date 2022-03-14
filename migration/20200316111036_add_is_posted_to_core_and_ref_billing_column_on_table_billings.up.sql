DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'posted_to_core_enum') THEN
        CREATE TYPE posted_to_core_enum AS ENUM (
            'yes',
            'no'
        );
    END IF;
END
$$;

ALTER TABLE billings
ADD COLUMN posted_to_core posted_to_core_enum DEFAULT 'no';

ALTER TABLE billings
ADD COLUMN ref_billing VARCHAR(100) DEFAULT NULL