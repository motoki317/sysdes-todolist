INSERT INTO `users` (`id`, `name`, `password`) VALUES (1, 'test', '$2a$10$WloBSQLMX2qf0qIqiJdbrOP3LxdY.5OimGkmp1o3O0D8FEX6BB4JC');

INSERT INTO `tasks` (`user_id`, `title`) VALUES (1, 'sample-task-01');
INSERT INTO `tasks` (`user_id`, `title`) VALUES (1, 'sample-task-02');
INSERT INTO `tasks` (`user_id`, `title`) VALUES (1, 'sample-task-03');
INSERT INTO `tasks` (`user_id`, `title`) VALUES (1, 'sample-task-04');
INSERT INTO `tasks` (`user_id`, `title`, `is_done`) VALUES (1, 'sample-task-05', true);
