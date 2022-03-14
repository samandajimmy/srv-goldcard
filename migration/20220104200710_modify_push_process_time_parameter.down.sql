UPDATE parameters SET "key" = 'PUSH_PROCESS_TIME', description = 'Jam push process' where "key" = 'PUSH_PROCESS_START_TIME'
;

DELETE FROM parameters WHERE "key" = 'PUSH_PROCESS_STOP_TIME'
;

DELETE FROM parameters WHERE "key" = 'PUSH_PROCESS_SEC_INTERVAL'
;