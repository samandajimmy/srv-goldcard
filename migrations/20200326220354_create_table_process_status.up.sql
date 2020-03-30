-- create table pegadaian_billings
CREATE TABLE IF NOT EXISTS process_status (
    id SERIAL PRIMARY KEY NOT NULL,
    process_id VARCHAR(20),
    process_type varchar(50),
    status TEXT,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_process_status ON process_status (id);

ALTER TABLE applications
ADD COLUMN process_id VARCHAR(50) DEFAULT NULL;

ALTER TABLE applications
ADD COLUMN error BOOLEAN DEFAULT FALSE;