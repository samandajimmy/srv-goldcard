DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'cards_status_enum') THEN
        ALTER TABLE cards
        ALTER COLUMN status DROP DEFAULT,
        ALTER COLUMN status TYPE status_enum_default USING status::TEXT::status_enum_default,
        ALTER COLUMN status SET DEFAULT 'inactive'::status_enum_default;
        DROP TYPE cards_status_enum;
    END IF;
END
$$;