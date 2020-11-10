ALTER TABLE accounts 
DROP COLUMN correspondence_id;

DROP INDEX index_correspondences;
DROP TABLE IF EXISTS correspondences;
