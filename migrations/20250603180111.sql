-- Modify "users" table
ALTER TABLE `users` ADD COLUMN `is_template_manager` bool NULL DEFAULT 0;
