CREATE TABLE `comments` (
                            `id` bigint NOT NULL AUTO_INCREMENT,
                            `content` text NOT NULL,
                            `user_id` bigint NOT NULL,
                            `post_id` bigint NOT NULL,
                            `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                            PRIMARY KEY (`id`),
                            KEY `idx_post_id` (`post_id`),
                            KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;