CREATE TABLE IF NOT EXISTS correspondences (
    id SERIAL PRIMARY KEY NOT NULL,
    address_line_1 VARCHAR(255),
    address_line_2 VARCHAR(255),
    address_line_3 VARCHAR(255),
    address_city VARCHAR(100),
    zipcode VARCHAR(50),
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_correspondences ON correspondences (id, created_at);

ALTER TABLE accounts
ADD COLUMN correspondence_id INTEGER REFERENCES correspondences(id) UNIQUE;
