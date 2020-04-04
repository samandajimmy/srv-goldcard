DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_inquiry_status_enum') THEN
        CREATE TYPE payment_inquiry_status_enum AS ENUM (
            'unpaid',
            'paid'
        );
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS payment_inquiries (
    id SERIAL PRIMARY KEY NOT NULL,
    account_id INTEGER REFERENCES accounts(id),
    billing_id INTEGER REFERENCES billings(id),
    ref_trx VARCHAR(100) NOT NULL,
    nominal INTEGER,
    status payment_inquiry_status_enum NOT NULL DEFAULT 'unpaid',
    inquiry_date TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_payment_inquiries ON payment_inquiries (id, ref_trx, created_at);