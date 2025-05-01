-- Modify "workspace_templates" table
ALTER TABLE `workspace_templates` ADD INDEX `idx_workspace_templates_created_at` (`created_at`), ADD INDEX `idx_workspace_templates_updated_at` (`updated_at`);
-- Create "files" table
CREATE TABLE `files` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `filepath` longtext NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_files_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Modify "git_workspace_sources" table
ALTER TABLE `git_workspace_sources` DROP COLUMN `files`, ADD COLUMN `sources_id` bigint unsigned NULL, ADD INDEX `fk_git_workspace_sources_sources` (`sources_id`), ADD CONSTRAINT `fk_git_workspace_sources_sources` FOREIGN KEY (`sources_id`) REFERENCES `files` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "workspace_template_versions" table
ALTER TABLE `workspace_template_versions` DROP COLUMN `files`, ADD COLUMN `sources_id` bigint unsigned NULL, ADD COLUMN `published` bool NULL DEFAULT 0, ADD COLUMN `edited_by_id` bigint unsigned NULL, ADD COLUMN `edited_on` datetime(3) NULL, ADD INDEX `fk_workspace_template_versions_edited_by` (`edited_by_id`), ADD INDEX `fk_workspace_template_versions_sources` (`sources_id`), ADD CONSTRAINT `fk_workspace_template_versions_edited_by` FOREIGN KEY (`edited_by_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT `fk_workspace_template_versions_sources` FOREIGN KEY (`sources_id`) REFERENCES `files` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
