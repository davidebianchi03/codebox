-- Create "password_reset_tokens" table
CREATE TABLE `password_reset_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `token` varchar(255) NOT NULL,
  `expiration` datetime(3) NOT NULL,
  `created_at` bigint NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_password_reset_tokens_user_id` (`user_id`),
  UNIQUE INDEX `uni_password_reset_tokens_token` (`token`),
  CONSTRAINT `fk_password_reset_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
