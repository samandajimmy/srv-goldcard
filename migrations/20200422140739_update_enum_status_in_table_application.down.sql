DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'application_status_enum') THEN
        ALTER TYPE application_status_enum rename TO application_status_enum_old;
        CREATE TYPE application_status_enum AS ENUM (
            'application_processed',
            'card_processed',
            'card_send',
            'card_sent',
            'failed',
            'active',
            'inactive'
        );
    END IF;
END
$$;

ALTER TABLE applications rename column status to status_old;
ALTER TABLE applications rename column rejected_date to failed_date;
ALTER TABLE applications add status application_status_enum NOT NULL DEFAULT 'inactive';
UPDATE applications SET status = status::TEXT::application_status_enum;
ALTER TABLE applications drop column status_old;
DROP TYPE application_status_enum_old;