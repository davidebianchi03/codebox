-- Create "git_workspace_sources" table
CREATE TABLE `git_workspace_sources` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `repository_url` text NOT NULL,
  `files` text NOT NULL
);
-- Create index "git_workspace_sources_repository_url" to table: "git_workspace_sources"
CREATE UNIQUE INDEX `git_workspace_sources_repository_url` ON `git_workspace_sources` (`repository_url`);
-- Create index "git_workspace_sources_files" to table: "git_workspace_sources"
CREATE UNIQUE INDEX `git_workspace_sources_files` ON `git_workspace_sources` (`files`);
-- Create index "idx_git_workspace_sources_deleted_at" to table: "git_workspace_sources"
CREATE INDEX `idx_git_workspace_sources_deleted_at` ON `git_workspace_sources` (`deleted_at`);
-- Create "groups" table
CREATE TABLE `groups` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `name` text NULL
);
-- Create index "groups_name" to table: "groups"
CREATE UNIQUE INDEX `groups_name` ON `groups` (`name`);
-- Create index "idx_groups_deleted_at" to table: "groups"
CREATE INDEX `idx_groups_deleted_at` ON `groups` (`deleted_at`);
-- Create "runners" table
CREATE TABLE `runners` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `name` text NOT NULL,
  `token` text NOT NULL,
  `restricted` numeric NULL DEFAULT false,
  `use_public_url` numeric NULL DEFAULT false,
  `public_url` text NULL,
  `last_contact` datetime NULL
);
-- Create index "runners_name" to table: "runners"
CREATE UNIQUE INDEX `runners_name` ON `runners` (`name`);
-- Create index "runners_token" to table: "runners"
CREATE UNIQUE INDEX `runners_token` ON `runners` (`token`);
-- Create index "idx_runners_deleted_at" to table: "runners"
CREATE INDEX `idx_runners_deleted_at` ON `runners` (`deleted_at`);
-- Create "runner_allowed_groups" table
CREATE TABLE `runner_allowed_groups` (
  `runner_id` integer NULL,
  `group_id` integer NULL,
  PRIMARY KEY (`runner_id`, `group_id`),
  CONSTRAINT `fk_runner_allowed_groups_group` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_runner_allowed_groups_runner` FOREIGN KEY (`runner_id`) REFERENCES `runners` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "users" table
CREATE TABLE `users` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `email` text NOT NULL,
  `password` text NOT NULL,
  `first_name` text NULL,
  `last_name` text NULL,
  `ssh_private_key` text NOT NULL,
  `ssh_public_key` text NOT NULL,
  `is_superuser` numeric NULL DEFAULT false
);
-- Create index "users_email" to table: "users"
CREATE UNIQUE INDEX `users_email` ON `users` (`email`);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX `idx_users_deleted_at` ON `users` (`deleted_at`);
-- Create "tokens" table
CREATE TABLE `tokens` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `token` text NULL,
  `expiration_date` datetime NULL,
  `user_id` integer NULL,
  CONSTRAINT `fk_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "tokens_token" to table: "tokens"
CREATE UNIQUE INDEX `tokens_token` ON `tokens` (`token`);
-- Create index "idx_tokens_deleted_at" to table: "tokens"
CREATE INDEX `idx_tokens_deleted_at` ON `tokens` (`deleted_at`);
-- Create "user_groups" table
CREATE TABLE `user_groups` (
  `user_id` integer NULL,
  `group_id` integer NULL,
  PRIMARY KEY (`user_id`, `group_id`),
  CONSTRAINT `fk_user_groups_group` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_user_groups_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "workspace_templates" table
CREATE TABLE `workspace_templates` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `name` text NOT NULL,
  `description` text NULL,
  `icon` text NULL
);
-- Create index "workspace_templates_name" to table: "workspace_templates"
CREATE UNIQUE INDEX `workspace_templates_name` ON `workspace_templates` (`name`);
-- Create index "idx_workspace_templates_deleted_at" to table: "workspace_templates"
CREATE INDEX `idx_workspace_templates_deleted_at` ON `workspace_templates` (`deleted_at`);
-- Create "workspace_template_versions" table
CREATE TABLE `workspace_template_versions` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `template_id` integer NULL,
  `name` text NOT NULL,
  `files` text NOT NULL,
  CONSTRAINT `fk_workspace_template_versions_template` FOREIGN KEY (`template_id`) REFERENCES `workspace_templates` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_workspace_template_versions_deleted_at" to table: "workspace_template_versions"
CREATE INDEX `idx_workspace_template_versions_deleted_at` ON `workspace_template_versions` (`deleted_at`);
-- Create "workspaces" table
CREATE TABLE `workspaces` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `name` text NOT NULL,
  `user_id` integer NULL,
  `runner_id` integer NULL,
  `config_source` text NOT NULL,
  `template_version_id` integer NULL,
  `git_source_id` integer NULL,
  `environment_variables` text NULL,
  CONSTRAINT `fk_workspaces_template_version` FOREIGN KEY (`template_version_id`) REFERENCES `workspace_template_versions` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `fk_workspaces_runner` FOREIGN KEY (`runner_id`) REFERENCES `runners` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `fk_workspaces_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `fk_workspaces_git_source` FOREIGN KEY (`git_source_id`) REFERENCES `git_workspace_sources` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_workspaces_deleted_at" to table: "workspaces"
CREATE INDEX `idx_workspaces_deleted_at` ON `workspaces` (`deleted_at`);
-- Create "workspace_containers" table
CREATE TABLE `workspace_containers` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `workspace_id` integer NULL,
  `container_id` text NULL,
  `container_name` text NULL,
  `container_user` text NULL,
  `agent_last_contact` datetime NULL,
  CONSTRAINT `fk_workspace_containers_workspace` FOREIGN KEY (`workspace_id`) REFERENCES `workspaces` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_workspace_containers_deleted_at" to table: "workspace_containers"
CREATE INDEX `idx_workspace_containers_deleted_at` ON `workspace_containers` (`deleted_at`);
-- Create "workspace_container_ports" table
CREATE TABLE `workspace_container_ports` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `container_id` integer NULL,
  `port_number` integer NOT NULL,
  `public` numeric NULL DEFAULT false,
  CONSTRAINT `fk_workspace_container_ports_container` FOREIGN KEY (`container_id`) REFERENCES `workspace_containers` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_workspace_container_ports_deleted_at" to table: "workspace_container_ports"
CREATE INDEX `idx_workspace_container_ports_deleted_at` ON `workspace_container_ports` (`deleted_at`);
