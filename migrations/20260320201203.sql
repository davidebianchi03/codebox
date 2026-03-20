-- Create "analytics_configs" table
CREATE TABLE `analytics_configs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `send_analytics_data` bool NULL DEFAULT 0,
  `analytics_banner_sent` bool NULL DEFAULT 0,
  `last_attempt` datetime(3) NULL,
  `last_successfull_attempt` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_analytics_configs_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
