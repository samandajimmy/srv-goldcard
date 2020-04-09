-- create table pegadaian_billings
CREATE TABLE IF NOT EXISTS pegadaian_billings (
    id SERIAL PRIMARY KEY NOT NULL,
    ref_id VARCHAR(255),
    billing_date TIMESTAMP,
    file_name VARCHAR(100),
    file_base64 TEXT,
    file_extension VARCHAR(10),
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_pegadaian_billings ON pegadaian_billings (id, created_at);