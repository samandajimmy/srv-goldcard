-- create table pegadaian_billings
CREATE TABLE IF NOT EXISTS process_statuses (
    id SERIAL PRIMARY KEY NOT NULL,
    process_id INT,
    process_type varchar(50),
    tbl_name varchar(100),
    reason TEXT,
    error_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_process_statuses ON process_statuses (id);

ALTER TABLE applications
ADD COLUMN core_open BOOLEAN DEFAULT FALSE;