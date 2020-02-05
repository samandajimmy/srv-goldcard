DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transactions_type_enum_default') THEN
        CREATE TYPE transactions_type_enum_default AS ENUM (
            'debit',
            'credit'
        );
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transactions_status_enum_default') THEN
        CREATE TYPE transactions_status_enum_default AS ENUM (
            'posted',
            'pending',
            'canceled'
        );
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transactions_methods_enum_default') THEN
        CREATE TYPE transactions_methods_enum_default AS ENUM (
            'payment',
            'adjustments'
        );
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY NOT NULL,
    account_id INTEGER REFERENCES accounts(id) UNIQUE,
    ref_trx_pgdn VARCHAR(100) NOT NULL,
    ref_trx VARCHAR(100) NOT NULL,
    nominal INTEGER,
    gold_nominal FLOAT,
    type transactions_type_enum_default DEFAULT NULL,
    status transactions_status_enum_default DEFAULT NULL,
    balance INTEGER,
    gold_balance FLOAT,
    methods transactions_methods_enum_default DEFAULT NULL,
    trx_date TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_transactions ON transactions (id, ref_trx_pgdn, ref_trx, created_at);

CREATE TABLE IF NOT EXISTS billings (
    id SERIAL PRIMARY KEY NOT NULL,
    account_id INTEGER REFERENCES accounts(id) UNIQUE,
    amount INTEGER,
    gold_amount FLOAT,
    billing_date TIMESTAMP DEFAULT NULL,
    depth_amount INTEGER,
    depth_gold FLOAT,
    stl INTEGER,
    depth_stl INTEGER,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_billings ON billings (id, account_id, billing_date, created_at);

CREATE TABLE IF NOT EXISTS billing_transactions (
    id SERIAL PRIMARY KEY NOT NULL,
    trx_id INTEGER REFERENCES transactions(id) NOT NULL,
    bill_id INTEGER REFERENCES billings(id) NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_bill_transactions ON billing_transactions (id);

CREATE TABLE IF NOT EXISTS billing_payments (
    id SERIAL PRIMARY KEY NOT NULL,
    trx_id INTEGER REFERENCES transactions(id) NOT NULL,
    bill_id INTEGER REFERENCES billings(id) NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_billing_payments ON billing_payments (id);