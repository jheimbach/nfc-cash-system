CREATE TABLE transactions
(
    `id`         integer PRIMARY KEY NOT NULL AUTO_INCREMENT,
    `new_saldo`  decimal(15, 2)      NOT NULL,
    `old_saldo`  decimal(15, 2)      NOT NULL,
    `amount`     decimal(15, 2)      NOT NULL,
    `account_id` integer             NOT NULL,
    `created`    datetime            NOT NULL
);

ALTER TABLE `transactions`
    ADD FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`);


CREATE INDEX idx_created ON transactions (`created`);
CREATE INDEX idx_id ON transactions (`id`)