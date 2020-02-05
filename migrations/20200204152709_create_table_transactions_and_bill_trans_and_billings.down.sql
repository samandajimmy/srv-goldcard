DROP INDEX index_transactions;
DROP TABLE IF EXISTS transactions;

DROP INDEX index_billing_transactions;
DROP TABLE IF EXISTS billing_transactions;

DROP INDEX index_billings;
DROP TABLE IF EXISTS billings;

DROP TYPE transactions_type_enum_default;
DROP TYPE transactions_status_enum_default;
DROP TYPE transactions_methods_enum_default;
