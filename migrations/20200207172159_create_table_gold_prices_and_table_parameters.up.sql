-- create table gold_prices
CREATE TABLE IF NOT EXISTS gold_prices (
    id SERIAL PRIMARY KEY NOT NULL,
    price FLOAT,
    valid_date TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_gold_prices ON gold_prices (id, created_at);

-- create table parameters
CREATE TABLE IF NOT EXISTS parameters (
    id SERIAL PRIMARY KEY NOT NULL,
    key VARCHAR(100) NOT NULL,
    value VARCHAR(100) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_parameters ON parameters (id, created_at);