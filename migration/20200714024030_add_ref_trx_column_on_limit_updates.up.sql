ALTER TABLE limit_updates
ADD ref_trx VARCHAR(100) NULL;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_limit_update_enum') THEN
        ALTER TYPE status_limit_update_enum rename TO status_limit_update_enum_old;
        CREATE TYPE status_limit_update_enum AS ENUM (
            'pending',
            'applied',
            'approved',
            'rejected',
            'inquired',
            'force_applied'
        );
        ALTER TABLE limit_updates
        ALTER COLUMN status TYPE status_limit_update_enum USING status::TEXT::status_limit_update_enum;
        DROP TYPE status_limit_update_enum_old;
    END IF;
END
$$;