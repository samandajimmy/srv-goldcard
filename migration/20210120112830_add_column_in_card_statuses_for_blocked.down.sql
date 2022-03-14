ALTER TABLE card_statuses
DROP is_replaced,
DROP replaced_date,
DROP last_encrypted_card_number;
DROP TYPE is_replaced_enum;