ALTER TABLE billings
    RENAME COLUMN debt_amount TO depth_amount;

ALTER TABLE billings
    RENAME COLUMN debt_gold TO depth_gold;

ALTER TABLE billings
    RENAME COLUMN debt_stl TO depth_stl;

ALTER TABLE billing_payments
    DROP COLUMN source;