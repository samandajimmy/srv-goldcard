-- add api_requests table
drop index index_api_requests;
drop table api_requests;
drop type request_status_enum;

-- add documents table one to many with applications
drop index index_documents;
drop table documents;
drop type doc_type_enum;

alter table applications
add column ktp_image_base64 TEXT;

alter table applications
add column npwp_image_base64 TEXT;

alter table applications
add column selfie_image_base64 TEXT;

alter table applications
add column ktp_doc_id varchar(100);

alter table applications
add column npwp_doc_id varchar(100);

alter table applications
add column selfie_doc_id varchar(100);

-- add column last_request_id on table accounts
alter table accounts
drop column last_request_id;
