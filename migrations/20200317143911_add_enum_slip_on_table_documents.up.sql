DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'doc_type_enum') THEN
        ALTER TYPE doc_type_enum rename TO doc_type_enum_old;
        CREATE TYPE doc_type_enum AS ENUM (
            'ktp',
            'npwp',
            'selfie',
            'slip_te',
            'undefined'
        );
        ALTER TABLE documents
        ALTER COLUMN type DROP DEFAULT,
        ALTER COLUMN type SET DATA TYPE doc_type_enum USING type::TEXT::doc_type_enum,
        ALTER COLUMN type SET DEFAULT 'undefined';
        DROP TYPE doc_type_enum_old;
    END IF;
END
$$;
