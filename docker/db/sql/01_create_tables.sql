-- Table for tasks
DROP TABLE IF EXISTS `tasks`;
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(50) NOT NULL,
    `password` CHAR(60) NOT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
    `deleted_at` DATETIME(6) NULL,
    UNIQUE KEY `idx_users_name` (`name`)
) DEFAULT CHARSET=utf8mb4 ENGINE=InnoDB;

CREATE TABLE `tasks` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `user_id` BIGINT NOT NULL,
    `title` VARCHAR(50) NOT NULL,
    `is_done` BOOLEAN NOT NULL DEFAULT b'0',
    `created_at` DATETIME(6) NOT NULL DEFAULT NOW(6),
    `deleted_at` DATETIME(6) NULL,
    KEY `idx_tasks_user_id_id` (`user_id`, `id`),
    FOREIGN KEY `fk_tasks_user_id` (`user_id`) REFERENCES `users` (`id`)
) DEFAULT CHARSET=utf8mb4 ENGINE=InnoDB;
