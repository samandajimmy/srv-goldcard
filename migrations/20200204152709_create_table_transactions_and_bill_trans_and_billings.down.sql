DROP INDEX index_transactions;
DROP TABLE IF EXISTS transactions;

DROP INDEX index_billing_transactions;
DROP TABLE IF EXISTS billing_transactions;

DROP INDEX index_billing_payments;
DROP TABLE IF EXISTS billing_payments;

DROP INDEX index_billings;
DROP TABLE IF EXISTS billings;

DROP TYPE transactions_type_enum;
DROP TYPE transactions_status_enum;
DROP TYPE transactions_methods_enum;
