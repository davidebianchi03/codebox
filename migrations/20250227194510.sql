-- Add column "type" to table: "runners"
ALTER TABLE `runners` ADD COLUMN `type` text NULL;
-- Add column "config_source_file_path" to table: "workspaces"
ALTER TABLE `workspaces` ADD COLUMN `config_source_file_path` text NULL;
