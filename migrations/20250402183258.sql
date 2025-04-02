-- Modify "git_workspace_sources" table
ALTER TABLE `git_workspace_sources` ADD COLUMN `config_file_path` text NULL;
-- Modify "workspace_template_versions" table
ALTER TABLE `workspace_template_versions` ADD COLUMN `config_file_path` text NULL;
-- Modify "workspaces" table
ALTER TABLE `workspaces` DROP COLUMN `config_source_file_path`;
