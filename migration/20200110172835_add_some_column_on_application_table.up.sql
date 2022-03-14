DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'application_status_enum') THEN
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
ALTER TABLE applications add status application_status_enum NOT NULL DEFAULT 'inactive';
UPDATE applications SET status = status_old::TEXT::application_status_enum;
ALTER TABLE applications drop column status_old;

ALTER TABLE applications
ADD COLUMN application_processed_date TIMESTAMP DEFAULT NULL,
ADD COLUMN card_processed_date TIMESTAMP DEFAULT NULL,
ADD COLUMN card_sent_date TIMESTAMP DEFAULT NULL,
ADD COLUMN failed_date TIMESTAMP DEFAULT NULL,
ADD COLUMN card_send_date TIMESTAMP DEFAULT NULL;