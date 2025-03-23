-- Create "groups" table
CREATE TABLE `groups` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `name` varchar(255) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_groups_deleted_at` (`deleted_at`),
  UNIQUE INDEX `uni_groups_name` (`name`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "runners" table
CREATE TABLE `runners` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `token` varchar(255) NOT NULL,
  `type` varchar(255) NULL,
  `restricted` bool NULL DEFAULT 0,
  `use_public_url` bool NULL DEFAULT 0,
  `public_url` text NULL,
  `last_contact` datetime(3) NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_runners_deleted_at` (`deleted_at`),
  UNIQUE INDEX `uni_runners_name` (`name`),
  UNIQUE INDEX `uni_runners_token` (`token`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "runner_allowed_groups" table
CREATE TABLE `runner_allowed_groups` (
  `runner_id` bigint unsigned NOT NULL,
  `group_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`runner_id`, `group_id`),
  INDEX `fk_runner_allowed_groups_group` (`group_id`),
  CONSTRAINT `fk_runner_allowed_groups_group` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_runner_allowed_groups_runner` FOREIGN KEY (`runner_id`) REFERENCES `runners` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "users" table
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `password` longtext NOT NULL,
  `first_name` varchar(255) NULL,
  `last_name` varchar(255) NULL,
  `ssh_private_key` longtext NOT NULL,
  `ssh_public_key` longtext NOT NULL,
  `is_superuser` bool NULL DEFAULT 0,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_users_deleted_at` (`deleted_at`),
  UNIQUE INDEX `uni_users_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "tokens" table
CREATE TABLE `tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `token` varchar(255) NULL,
  `expiration_date` datetime(3) NULL,
  `user_id` bigint unsigned NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_tokens_user` (`user_id`),
  INDEX `idx_tokens_deleted_at` (`deleted_at`),
  UNIQUE INDEX `uni_tokens_token` (`token`),
  CONSTRAINT `fk_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "user_groups" table
CREATE TABLE `user_groups` (
  `user_id` bigint unsigned NOT NULL,
  `group_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`user_id`, `group_id`),
  INDEX `fk_user_groups_group` (`group_id`),
  CONSTRAINT `fk_user_groups_group` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_user_groups_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "git_workspace_sources" table
CREATE TABLE `git_workspace_sources` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `repository_url` text NOT NULL,
  `files` text NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_git_workspace_sources_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "workspace_templates" table
CREATE TABLE `workspace_templates` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NULL,
  `description` longtext NULL,
  `icon` text NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_workspace_templates_deleted_at` (`deleted_at`),
  UNIQUE INDEX `uni_workspace_templates_name` (`name`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "workspace_template_versions" table
CREATE TABLE `workspace_template_versions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `template_id` bigint unsigned NULL,
  `name` varchar(255) NOT NULL,
  `files` text NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_workspace_template_versions_template` (`template_id`),
  INDEX `idx_workspace_template_versions_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_workspace_template_versions_template` FOREIGN KEY (`template_id`) REFERENCES `workspace_templates` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "workspaces" table
CREATE TABLE `workspaces` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `user_id` bigint unsigned NULL,
  `status` varchar(30) NOT NULL,
  `type` varchar(255) NOT NULL,
  `runner_id` bigint unsigned NULL,
  `config_source` varchar(20) NOT NULL,
  `template_version_id` bigint unsigned NULL,
  `git_source_id` bigint unsigned NULL,
  `config_source_file_path` text NULL,
  `environment_variables` longtext NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_workspaces_git_source` (`git_source_id`),
  INDEX `fk_workspaces_runner` (`runner_id`),
  INDEX `fk_workspaces_template_version` (`template_version_id`),
  INDEX `fk_workspaces_user` (`user_id`),
  INDEX `idx_workspaces_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_workspaces_git_source` FOREIGN KEY (`git_source_id`) REFERENCES `git_workspace_sources` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `fk_workspaces_runner` FOREIGN KEY (`runner_id`) REFERENCES `runners` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `fk_workspaces_template_version` FOREIGN KEY (`template_version_id`) REFERENCES `workspace_template_versions` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `fk_workspaces_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "workspace_containers" table
CREATE TABLE `workspace_containers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `workspace_id` bigint unsigned NULL,
  `container_id` varchar(255) NULL,
  `container_name` varchar(255) NULL,
  `container_image` varchar(255) NULL,
  `container_user_id` bigint unsigned NULL,
  `container_user_name` varchar(255) NULL,
  `agent_last_contact` datetime(3) NULL,
  `workspace_path` varchar(255) NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_workspace_containers_workspace` (`workspace_id`),
  INDEX `idx_workspace_containers_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_workspace_containers_workspace` FOREIGN KEY (`workspace_id`) REFERENCES `workspaces` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "workspace_container_ports" table
CREATE TABLE `workspace_container_ports` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `container_id` bigint unsigned NULL,
  `service_name` varchar(255) NOT NULL,
  `port_number` bigint unsigned NOT NULL,
  `public` bool NULL DEFAULT 0,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_workspace_container_ports_container` (`container_id`),
  INDEX `idx_workspace_container_ports_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_workspace_container_ports_container` FOREIGN KEY (`container_id`) REFERENCES `workspace_containers` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
