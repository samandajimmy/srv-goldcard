
CREATE TABLE IF NOT EXISTS gte_payments (
    id SERIAL PRIMARY KEY NOT NULL,
    account_id INTEGER REFERENCES accounts(id),
    trx_id VARCHAR(100) NOT NULL,
    gold_amount FLOAT NOT NULL,
    trx_amount INTEGER NOT NULL,
    bri_updated BOOLEAN DEFAULT 'no',
    pds_notified BOOLEAN DEFAULT 'no',
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_gte_payments ON gte_payments (id, trx_id, created_at);