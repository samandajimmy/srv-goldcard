DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'cards_status_enum') THEN
        ALTER TYPE cards_status_enum RENAME TO cards_status_enum_old;
        CREATE TYPE cards_status_enum AS ENUM (
            'inactive',
            'active'
        );
        
        ALTER TABLE cards
        ALTER COLUMN status DROP DEFAULT,
        ALTER COLUMN status TYPE cards_status_enum USING status::TEXT::cards_status_enum,
        ALTER COLUMN status SET DEFAULT 'inactive'::cards_status_enum;
        DROP TYPE cards_status_enum_old;
    END IF;
END
$$;