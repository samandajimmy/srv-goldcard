ALTER TABLE cards
ADD previous_card_balance INTEGER DEFAULT NULL,
ADD previous_card_balance_date TIMESTAMP DEFAULT NULL,
ADD previous_card_limit INTEGER DEFAULT NULL,
ADD previous_card_limit_date TIMESTAMP DEFAULT NULL;

ALTER TABLE applications
ADD card_limit INTEGER DEFAULT NULL;