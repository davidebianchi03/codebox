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
  `ssh_public_key` text NOT NULL
);
-- Create index "users_id" to table: "users"
CREATE UNIQUE INDEX `users_id` ON `users` (`id`);
-- Create index "users_email" to table: "users"
CREATE UNIQUE INDEX `users_email` ON `users` (`email`);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX `idx_users_deleted_at` ON `users` (`deleted_at`);
