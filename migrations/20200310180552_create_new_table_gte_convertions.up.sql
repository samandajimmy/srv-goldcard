CREATE TABLE IF NOT EXISTS gte_convertions (
    id SERIAL PRIMARY KEY NOT NULL,
    billing_id INTEGER REFERENCES billings(id),
    amount INTEGER,
    gold_amount FLOAT,
    gte_date TIMESTAMP DEFAULT NULL,
    credit_id VARCHAR(100) DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);