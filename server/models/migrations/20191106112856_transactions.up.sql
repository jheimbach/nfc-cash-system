CREATE TABLE transactions
(
    `id`         integer PRIMARY KEY NOT NULL AUTO_INCREMENT,
    `new_saldo`  float               NOT NULL,
    `old_saldo`  float               NOT NULL,
    `amount`     float               NOT NULL,
    `account_id` integer             NOT NULL,
    `created`    datetime            NOT NULL
);

ALTER TABLE `transactions`
    ADD FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`);