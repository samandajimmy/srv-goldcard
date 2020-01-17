-- rename the existing type
ALTER TYPE application_status_enum RENAME TO application_status_enum_old;

-- create the new type
CREATE TYPE application_status_enum AS ENUM (
    'application_ongoing',
    'application_processed',
    'card_processed',
    'card_send',
    'card_sent',
    'failed',
    'active',
    'inactive'
);

ALTER TABLE applications rename column status to status_old;
ALTER TABLE applications add status application_status_enum NOT NULL DEFAULT 'inactive';
UPDATE applications SET status = status_old::TEXT::application_status_enum;
ALTER TABLE applications drop column status_old;

-- if you get an error, see bottom of post
-- remove the old type
DROP TYPE application_status_enum_old;

ALTER TABLE applications
ADD COLUMN current_step SMALLINT DEFAULT 0;