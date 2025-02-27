-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_git_workspace_sources" table
CREATE TABLE `new_git_workspace_sources` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `created_at` datetime NULL,
  `updated_at` datetime NULL,
  `deleted_at` datetime NULL,
  `repository_url` text NOT NULL,
  `files` text NULL
);
-- Copy rows from old table "git_workspace_sources" to new temporary table "new_git_workspace_sources"
INSERT INTO `new_git_workspace_sources` (`id`, `created_at`, `updated_at`, `deleted_at`, `repository_url`, `files`) SELECT `id`, `created_at`, `updated_at`, `deleted_at`, `repository_url`, `files` FROM `git_workspace_sources`;
-- Drop "git_workspace_sources" table after copying rows
DROP TABLE `git_workspace_sources`;
-- Rename temporary table "new_git_workspace_sources" to "git_workspace_sources"
ALTER TABLE `new_git_workspace_sources` RENAME TO `git_workspace_sources`;
-- Create index "git_workspace_sources_files" to table: "git_workspace_sources"
CREATE UNIQUE INDEX `git_workspace_sources_files` ON `git_workspace_sources` (`files`);
-- Create index "git_workspace_sources_repository_url" to table: "git_workspace_sources"
CREATE UNIQUE INDEX `git_workspace_sources_repository_url` ON `git_workspace_sources` (`repository_url`);
-- Create index "idx_git_workspace_sources_deleted_at" to table: "git_workspace_sources"
CREATE INDEX `idx_git_workspace_sources_deleted_at` ON `git_workspace_sources` (`deleted_at`);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
