CREATE TABLE `users_refreshkeys`
(
    `id`          integer PRIMARY KEY NOT NULL AUTO_INCREMENT,
    `user_id`     integer UNIQUE      NOT NULL,
    `refresh_key` char(64) UNIQUE     NOT NULL
);

ALTER TABLE `users_refreshkeys`
    ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

CREATE INDEX idx_id ON users_refreshkeys (`id`);