DROP TABLE IF EXISTS process_status;

ALTER TABLE applications
DROP COLUMN process_id;

ALTER TABLE applications
DROP COLUMN error;