CREATE TABLE IF NOT EXISTS disputed_transactions (
    id SERIAL PRIMARY KEY NOT NULL,
    transaction_id INTEGER REFERENCES transactions(id) UNIQUE,
    disputed_date TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_disputed_transactions ON disputed_transactions (id, created_at);