CREATE TABLE `accounts`
(
        `id` integer PRIMARY KEY NOT NULL AUTO_INCREMENT,
        `name` varchar(255) NOT NULL,
        `description` TEXT,
        `saldo` FLOAT NOT NULL DEFAULT 0,
        `group_id` INTEGER
);

ALTER TABLE `accounts` ADD FOREIGN KEY (`group_id`) REFERENCES `account_groups`(`id`);