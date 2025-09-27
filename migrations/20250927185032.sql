-- Modify "tokens" table
ALTER TABLE `tokens` ADD COLUMN `impersonated_user_id` bigint unsigned NULL AFTER `user_id`, ADD INDEX `fk_tokens_impersonated_user` (`impersonated_user_id`), ADD CONSTRAINT `fk_tokens_impersonated_user` FOREIGN KEY (`impersonated_user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE;
-- Create "impersonation_logs" table
CREATE TABLE `impersonation_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `token_id` bigint unsigned NULL,
  `impersonated_user_id` bigint unsigned NOT NULL,
  `impersonator_ip_address` longtext NOT NULL,
  `impersonation_started_at` datetime(3) NOT NULL,
  `impersonation_finished_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_impersonation_logs_impersonator` (`impersonated_user_id`),
  INDEX `fk_impersonation_logs_token` (`token_id`),
  INDEX `idx_impersonation_logs_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_impersonation_logs_impersonated_user` FOREIGN KEY (`impersonated_user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `fk_impersonation_logs_impersonator` FOREIGN KEY (`impersonated_user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `fk_impersonation_logs_token` FOREIGN KEY (`token_id`) REFERENCES `tokens` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
