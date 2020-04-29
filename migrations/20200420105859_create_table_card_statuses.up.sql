
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'is_reactivated_enum') THEN
        CREATE TYPE is_reactivated_enum AS ENUM (
            'yes',
            'no'
        );
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'reason_code_enum') THEN
        CREATE TYPE reason_code_enum AS ENUM (
            'lost',
            'stolen',
            'other'
        );
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS card_statuses (
    id SERIAL PRIMARY KEY NOT NULL,
    card_id INTEGER REFERENCES cards(id) NOT NULL,
    reason VARCHAR(255),
    reason_code reason_code_enum DEFAULT NULL,
    blocked_date TIMESTAMP DEFAULT NULL,
    is_reactivated is_reactivated_enum DEFAULT NULL,
    reactivated_date TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_card_statuses ON card_statuses (id, created_at);
