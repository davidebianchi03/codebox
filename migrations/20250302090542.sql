-- Add column "service_name" to table: "workspace_container_ports"
ALTER TABLE `workspace_container_ports` ADD COLUMN `service_name` text NOT NULL;
