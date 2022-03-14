ALTER TABLE cards
    ADD COLUMN stl_balance INTEGER;

ALTER TABLE cards
    RENAME COLUMN current_stl TO stl_limit;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'billing_status_enum') THEN
        CREATE TYPE billing_status_enum AS ENUM (
            'unpaid',
            'paid',
            'gte_converted'
        );
    END IF;
END
$$;

ALTER TABLE billings
    ADD COLUMN minimum_payment FLOAT,
    ADD COLUMN status billing_status_enum DEFAULT 'unpaid',
    ADD COLUMN billing_due_date TIMESTAMP DEFAULT NULL;