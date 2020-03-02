ALTER TABLE billings
ADD COLUMN status status_enum_default NOT NULL DEFAULT 'inactive';