INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('CORE_EOD_INTERVAL_HEALTH_CHECK', '3600', 'Interval waktu untuk melakukan healtcheck dalam satuan detik', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('CORE_EOD_HEALTH', 'true', 'status core service', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('UPDATE_LIMIT_EMAIL_ADDRESS', 'email@pegadaian.co.id', 'Alamat email untuk kirim batch csv update limit', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('BILLING_MIN_PAYMENT', '0.1', 'Minimum pembayaran user', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('BILLING_STATEMENT_TIME', '06:00', 'Waktu Cetak Tagihan ketika membuat cetak tagihan pada Tanggal Tagihan', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('BILLING_INTERVAL_DUE_DATE', '14', 'Jangka waktu jatuh tempo pembayaran sejak tanggal penagihan', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('CORE_EOD_TIME', '22:00', 'Waktu EOD dimulai', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('UPDATE_LIMIT_TO_BRI_TIME', '06:00', 'Jam untuk mengajukan limit kartu baru ke BRI', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('BILLING_PRINT_DATE', '02', 'Tanggal Jatuh Tempo Tagihan', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('REMIND_ACT_NOTIF_TIME', '10:00', 'Jam reminder aktivasi', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('REMIND_ACT_NOTIF_DATE', '1', 'Tanggal reminder aktivasi', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('DAILY_SYNC_TRANSACTION_TIME', '09:00', 'Waktu mensingkronisasi pending dan posted transaction ke BRI', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('CONVERT_GTE_DATE', '16', 'Tanggal convert GTE', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('SIGNATORY_NIP', 'Assistant Vice President', 'SIGNATORY_NIP (Untuk dokumen Slip TE)', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('SIGNATORY_NAME', 'Heri Prasongko', 'SIGNATORY_NAME', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('BILLING_STATEMENT_DATE', '03', 'Tanggal Jatuh Tempo Tagihan (Sepertinya duplicate dengan billing_print_date)', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('UPDATE_LIMIT_INQUIRIES_CLOSED_DATE', '15,16', 'UPDATE_LIMIT_INQUIRIES_CLOSED_DATE', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('CONVERT_GTE_TIME', '09:00', 'Jam convert GTE', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('PUBLIC_HOLIDAY_DATE', '', 'Tanggal Merah dan scheduller UPDATE_LIMIT_TO_BRI tidak dijalankan', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('CHECK_STATUS_UPDATE_LIMIT_TO_BRI_TIME', '07:00', 'Jam untuk mengecek status limit kartu yang sudah diajukan ke BRI', now(), NULL);
INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('PUSH_PROCESS_TIME', '07:00', 'Jam Push Process', now(), NULL);
