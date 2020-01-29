CREATE TABLE `users`
(
    `id`              integer PRIMARY KEY NOT NULL AUTO_INCREMENT,
    `name`            varchar(255)        NOT NULL,
    `email`           varchar(255) UNIQUE NOT NULL,
    `hashed_password` char(60)            NOT NULL,
    `created`         datetime            NOT NULL
);