CREATE TABLE `article` (
    `id` BIGINT,
    `title` VARCHAR(256) NULL,
    `description` VARCHAR(256) NULL,
    `body` TEXT NULL,
    `author` VARCHAR(64),
    `date` DATE NULL,
    PRIMARY KEY (`id`)
);