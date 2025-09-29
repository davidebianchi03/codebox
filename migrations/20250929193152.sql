-- Modify "runners" table
ALTER TABLE `runners` ADD COLUMN `deletion_in_progress` bool NOT NULL DEFAULT 0 AFTER `version`;
-- Modify "users" table
ALTER TABLE `users` ADD COLUMN `deletion_in_progress` bool NOT NULL DEFAULT 0 AFTER `is_template_manager`;
