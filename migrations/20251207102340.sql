-- Modify "users" table
ALTER TABLE `users` ADD COLUMN `email_verified` bool NOT NULL DEFAULT 0 AFTER `deletion_in_progress`;
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
