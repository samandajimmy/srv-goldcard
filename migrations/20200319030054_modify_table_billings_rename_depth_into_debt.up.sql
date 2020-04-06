ALTER TABLE billings
    RENAME COLUMN depth_amount TO debt_amount;

ALTER TABLE billings
    RENAME COLUMN depth_gold TO debt_gold;

ALTER TABLE billings
    RENAME COLUMN depth_stl TO debt_stl;

ALTER TABLE billing_payments
    ADD COLUMN source VARCHAR(20);