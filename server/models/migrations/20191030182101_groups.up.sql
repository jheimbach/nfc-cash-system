CREATE TABLE `account_groups`
(
    `id` integer PRIMARY KEY NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `description` TEXT,
    `can_overdraw` BOOLEAN NOT NULL DEFAULT false
)