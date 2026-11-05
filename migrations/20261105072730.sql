-- Modify "runners" table
ALTER TABLE `runners` ADD COLUMN `description` text NULL AFTER `token`;
