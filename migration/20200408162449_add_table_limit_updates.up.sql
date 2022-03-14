CREATE TABLE IF NOT EXISTS limit_updates (
    id SERIAL PRIMARY KEY NOT NULL,
    ref_id VARCHAR(100) NOT NULL,
    limit_date TIMESTAMP,
    account_id INTEGER REFERENCES accounts(id),
    card_limit INTEGER,
    gold_limit FLOAT,
    stl_limit INTEGER,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_limit_updates ON limit_updates (id, ref_id, created_at);