ALTER TABLE applications rename column status to status_old;
ALTER TABLE applications add status status_enum_default NOT NULL DEFAULT 'inactive';
UPDATE applications SET status = status_old::TEXT::status_enum_default;
ALTER TABLE applications drop column status_old;

DROP TYPE application_status_enum

ALTER TABLE applications
DROP COLUMN application_processed_date,
DROP COLUMN card_processed_date,
DROP COLUMN card_sent_date,
DROP COLUMN failed_date,
DROP COLUMN card_send_date;