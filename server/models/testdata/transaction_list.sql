INSERT INTO `accounts` (id, name, saldo, group_id,nfc_chip_uid)
VALUES (2, 'testaccount1', 120,1,'testchipid2');

INSERT INTO `transactions` (new_saldo, old_saldo, amount, account_id, created)
VALUES (115, 120, -5, 1, UTC_TIMESTAMP),
       (110, 115, -5, 1, UTC_TIMESTAMP),
       (105, 110, -5, 1, UTC_TIMESTAMP),
       (100, 105, -5, 1, UTC_TIMESTAMP),
       (95, 100, -5, 1, UTC_TIMESTAMP),
       (115, 120, -5, 2, UTC_TIMESTAMP),
       (110, 115, -5, 2, UTC_TIMESTAMP),
       (105, 110, -5, 2, UTC_TIMESTAMP),
       (100, 105, -5, 2, UTC_TIMESTAMP);
