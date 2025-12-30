-- Create "authentication_settings" table
CREATE TABLE `authentication_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `is_signup_open` bool NULL DEFAULT 0,
  `is_signup_restricted` bool NULL DEFAULT 0,
  `allowed_email_regex` text NULL,
  `blocked_email_regex` text NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_authentication_settings_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Drop "instance_settings" table
DROP TABLE `instance_settings`;
