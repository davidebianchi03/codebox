-- Rename a column from "container_user" to "container_user_name"
ALTER TABLE `workspace_containers` RENAME COLUMN `container_user` TO `container_user_name`;
-- Add column "container_image" to table: "workspace_containers"
ALTER TABLE `workspace_containers` ADD COLUMN `container_image` text NULL;
