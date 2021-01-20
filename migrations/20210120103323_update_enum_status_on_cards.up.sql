CREATE TYPE cards_status_enum AS ENUM (
    'inactive',
    'active',
    'blocked'
);

ALTER TABLE cards
ALTER COLUMN status DROP DEFAULT,
ALTER COLUMN status TYPE cards_status_enum USING status::TEXT::cards_status_enum,
ALTER COLUMN status SET DEFAULT 'inactive'::cards_status_enum;