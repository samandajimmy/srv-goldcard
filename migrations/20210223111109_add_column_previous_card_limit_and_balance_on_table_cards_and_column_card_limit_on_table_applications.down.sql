ALTER TABLE cards
DROP previous_card_balance,
DROP previous_card_balance_date,
DROP previous_card_limit,
DROP previous_card_limit_date;

ALTER TABLE applications
DROP card_limit;