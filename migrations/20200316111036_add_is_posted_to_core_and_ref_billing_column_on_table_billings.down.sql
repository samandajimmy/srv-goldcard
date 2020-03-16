DROP TYPE posted_to_core_enum;

ALTER TABLE billings
DROP COLUMN posted_to_core;

ALTER TABLE billings
DROP COLUMN ref_billing;