DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transactions_methods_enum') THEN
        ALTER TYPE transactions_methods_enum rename TO transactions_methods_enum_old;
        CREATE TYPE transactions_methods_enum AS ENUM (
            'payment',
            'adjustment'
        );
        ALTER TABLE transactions
        ALTER COLUMN methods TYPE transactions_methods_enum USING methods::TEXT::transactions_methods_enum;
        DROP TYPE transactions_methods_enum_old;
    END IF;
END
$$;
