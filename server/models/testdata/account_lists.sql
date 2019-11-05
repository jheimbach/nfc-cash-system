INSERT INTO `account_groups` (id, name, description)
VALUES (1, 'testgroup1', ''),
       (2, 'testgroup2', 'with description');

INSERT INTO accounts (name, group_id)
VALUES ('testuser1', 1),
       ('testuser2', 1),
       ('testuser3', 1),
       ('testuser4', 1),
       ('testuser5', 1),
       ('testuser6', 2),
       ('testuser7', 2),
       ('testuser8', 2),
       ('testuser9', 2);
