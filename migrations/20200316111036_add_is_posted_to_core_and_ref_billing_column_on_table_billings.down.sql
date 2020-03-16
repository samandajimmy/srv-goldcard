DROP TYPE is_posted_to_core_enum;

ALTER TABLE billings
DROP COLUMN is_posted_to_core;

ALTER TABLE billings
DROP COLUMN ref_billing;