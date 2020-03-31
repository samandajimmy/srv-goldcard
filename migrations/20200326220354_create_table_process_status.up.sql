-- create table pegadaian_billings
CREATE TABLE IF NOT EXISTS process_statuses (
    id SERIAL PRIMARY KEY NOT NULL,
    process_id VARCHAR(50),
    process_type varchar(50),
    status TEXT,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_process_statuses ON process_statuses (id);

ALTER TABLE applications
ADD COLUMN process_id VARCHAR(50) DEFAULT NULL;

ALTER TABLE applications
ADD COLUMN error BOOLEAN DEFAULT FALSE;