ALTER TABLE applications
DROP expired_at;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'application_status_enum') THEN
        ALTER TYPE application_status_enum rename TO application_status_enum_old;
        CREATE TYPE application_status_enum AS ENUM (
            'application_ongoing',
            'application_processed',
            'card_processed',
            'card_send',
            'card_sent',
            'card_suspended',
            'rejected',
            'inactive',
            'active'
        );
        ALTER TABLE applications
        ALTER COLUMN status DROP DEFAULT,
        ALTER COLUMN status TYPE application_status_enum USING status::TEXT::application_status_enum,
        ALTER COLUMN status SET DEFAULT 'inactive'::application_status_enum;
        DROP TYPE application_status_enum_old;
    END IF;
END
$$;
