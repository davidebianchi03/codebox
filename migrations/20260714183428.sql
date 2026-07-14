-- Modify "workspace_containers" table
ALTER TABLE `workspace_containers` ADD COLUMN `order` bigint NULL AFTER `workspace_path`;
-- Create "container_port_backups" table
CREATE TABLE `container_port_backups` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `workspace_id` bigint unsigned NULL,
  `container_name` varchar(255) NOT NULL,
  `service_name` varchar(255) NOT NULL,
  `port_number` bigint unsigned NOT NULL,
  `public` bool NULL DEFAULT 0,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_container_port_backups_workspace` (`workspace_id`),
  INDEX `idx_container_port_backups_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_container_port_backups_workspace` FOREIGN KEY (`workspace_id`) REFERENCES `workspaces` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
