ALTER TABLE cards
    DROP COLUMN stl_balance;

ALTER TABLE cards
    RENAME COLUMN stl_limit TO current_stl;

ALTER TABLE billings
    DROP COLUMN minimum_payment,
    DROP COLUMN status,
    DROP COLUMN billing_due_date;

DROP TYPE billing_status_enum;