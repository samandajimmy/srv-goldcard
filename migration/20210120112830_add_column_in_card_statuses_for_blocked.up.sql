CREATE TYPE is_replaced_enum AS ENUM (
    'yes',
    'no'
);

ALTER TABLE card_statuses
ADD is_replaced is_replaced_enum DEFAULT 'no'::is_replaced_enum,
ADD replaced_date TIMESTAMP DEFAULT NULL,
ADD last_encrypted_card_number VARCHAR(100) DEFAULT NULL;