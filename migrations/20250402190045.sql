-- Modify "git_workspace_sources" table
ALTER TABLE `git_workspace_sources` ADD COLUMN `ref_name` varchar(255) NOT NULL;
