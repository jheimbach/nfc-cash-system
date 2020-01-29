CREATE TABLE `accounts`
(
    `id`          INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT,
    `name`        VARCHAR(255)        NOT NULL,
    `description` TEXT,
    `saldo`       DECIMAL(15,2)               NOT NULL DEFAULT 0,
    `group_id`    INTEGER,
    # according to ISO 14443-3A for nfc tags, uids are 4-10 bytes long (hex 2 chars per byte are 20 max).
    # source: https://www.nxp.com/docs/en/application-note/AN10927.pdf
    `nfc_chip_uid`     char(20) UNIQUE  NOT NULL
);

ALTER TABLE `accounts`
    ADD FOREIGN KEY (`group_id`) REFERENCES `account_groups` (`id`);

CREATE INDEX idx_id ON accounts(`id`);
CREATE INDEX idx_nfc_chip_uid ON accounts(`nfc_chip_uid`)