-- Create "instance_settings" table
CREATE TABLE `instance_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `allow_user_sign_up` bool NULL DEFAULT 0,
  `sign_up_restricted` bool NULL DEFAULT 0,
  `allowed_email_regex` text NULL,
  `blocked_email_regex` text NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_instance_settings_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
