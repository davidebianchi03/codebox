-- Create "instance_settings" table
CREATE TABLE `instance_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `is_signup_open` bool NULL DEFAULT 0,
  `is_sign_up_restricted` bool NULL DEFAULT 0,
  `allowed_email_regex` text NULL,
  `blocked_email_regex` text NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_instance_settings_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
