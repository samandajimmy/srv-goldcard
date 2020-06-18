ALTER TABLE gte_convertions
ADD account_id INTEGER REFERENCES accounts(id) DEFAULT NULL;