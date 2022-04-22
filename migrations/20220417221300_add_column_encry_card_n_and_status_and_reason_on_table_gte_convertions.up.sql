CREATE TYPE gte_convertions_status AS ENUM (
    'success',
    'error-billkey-not-found',
    'error-post-to-core',
    'error-calculate-card-limit',
    'error-bill-payment',
    'error-notification-pds',
    'error-internal-gc-api'
);

ALTER TABLE gte_convertions
ADD encrypted_card_number VARCHAR(20) DEFAULT NULL,
ADD status gte_convertions_status DEFAULT NULL,
ADD reason VARCHAR(255) DEFAULT NULL;