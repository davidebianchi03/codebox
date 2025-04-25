-- Create "authorization_codes" table
CREATE TABLE `authorization_codes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NULL,
  `token_id` bigint unsigned NULL,
  `expires_at` datetime(3) NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_authorization_codes_token` (`token_id`),
  INDEX `idx_authorization_codes_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_authorization_codes_token` FOREIGN KEY (`token_id`) REFERENCES `tokens` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
