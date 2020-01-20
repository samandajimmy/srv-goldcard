-- add api_requests table
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'request_status_enum') THEN
        CREATE TYPE request_status_enum AS ENUM (
            'success',
            'error'
        );
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS api_requests (
    id SERIAL PRIMARY KEY NOT NULL,
    request_id VARCHAR(255),
    host_name VARCHAR(255),
    endpoint VARCHAR(255),
    status request_status_enum NOT NULL DEFAULT 'success',
	request_data JSONB,
	response_data JSONB,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_api_requests ON api_requests (id, request_id, created_at);

-- add documents table one to many with applications
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'doc_type_enum') THEN
        CREATE TYPE doc_type_enum AS ENUM (
            'ktp',
            'npwp',
            'selfie',
            'undefined'
        );
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY NOT NULL,
    file_name VARCHAR(255),
    file_base64 text,
    file_extension varchar(10),
    doc_id VARCHAR(255),
    last_request_id varchar(255),
    type doc_type_enum NOT NULL DEFAULT 'undefined',
    application_id INTEGER REFERENCES applications(id),
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_documents ON documents (id, last_request_id, created_at);

alter table applications
drop column ktp_image_base64;

alter table applications
drop column npwp_image_base64;

alter table applications
drop column selfie_image_base64;

alter table applications
drop column ktp_doc_id;

alter table applications
drop column npwp_doc_id;

alter table applications
drop column selfie_doc_id;

-- add column last_request_id on table accounts
alter table accounts
add column last_request_id varchar(255);