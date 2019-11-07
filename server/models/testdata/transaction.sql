INSERT INTO `account_groups` (id, name, description)
VALUES (1, 'testgroup1', NULL);

INSERT INTO `accounts` (id, name, saldo, group_id, nfc_chip_uid)
VALUES (1, 'testaccount', 12,1,'testchipid');
