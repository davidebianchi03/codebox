-- Modify "users" table
ALTER TABLE `users` ADD COLUMN `approved` bool NULL DEFAULT 0 AFTER `is_template_manager`, ADD COLUMN `email_verified` bool NOT NULL DEFAULT 0 AFTER `deletion_in_progress`;
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
  `users_must_be_approved` bool NULL DEFAULT 0,
  `approved_by_default_email_regex` text NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_authentication_settings_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "email_verification_codes" table
CREATE TABLE `email_verification_codes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NOT NULL,
  `expiration` datetime(3) NULL,
  `user_id` bigint unsigned NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_email_verification_codes_user` (`user_id`),
  INDEX `idx_email_verification_codes_deleted_at` (`deleted_at`),
  UNIQUE INDEX `uni_email_verification_codes_code` (`code`),
  CONSTRAINT `fk_email_verification_codes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
