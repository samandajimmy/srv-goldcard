UPDATE parameters SET "key" = 'PUSH_PROCESS_START_TIME', description = 'Jam start push process' where "key" = 'PUSH_PROCESS_TIME'
;

INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('PUSH_PROCESS_STOP_TIME', '22:00', 'Jam stop push process', now(), NULL)
;

INSERT INTO parameters
("key", value, description, created_at, updated_at)
VALUES('PUSH_PROCESS_SEC_INTERVAL', '3600', 'Jumlah interval setiap detik untuk push proses', now(), NULL)
;