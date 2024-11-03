
CREATE TABLE `auth_users` (
  `id` bigint NOT NULL,
  `email` varchar(256) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `auth_users_email_unq_idx` (`email`)
);

CREATE TABLE `auth_tokens` (
  `token` char(36) NOT NULL,
  `user_id` bigint DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`token`),
  KEY `auth_tokens_user_id_idx` (`user_id`),
  CONSTRAINT `auth_tokens_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `auth_users` (`id`)
);

CREATE TABLE `auth_magiclink` (
  `token` char(36) NOT NULL,
  `email` varchar(256) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`token`),
);

CREATE TABLE `people_people` (
  `id` bigint NOT NULL,
  `first_name` varchar(64) NOT NULL DEFAULT '',
  `last_name` varchar(64) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  CONSTRAINT `people_people_id_fkey` FOREIGN KEY (`id`) REFERENCES `auth_users` (`id`) ON DELETE CASCADE
);

CREATE TABLE `people_searches` (
  `id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `s3_key` varchar(128) NOT NULL DEFAULT '',
  `created_at` bigint NOT NULL,
  `country_code` varchar(16) NOT NULL DEFAULT '',
  `api_response` text NOT NULL,
  PRIMARY KEY (`id`),
  KEY `people_searches_user_id_idx` (`user_id`),
  CONSTRAINT `people_searches_people_people_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `people_people` (`id`) ON DELETE CASCADE
);

CREATE TABLE `searches_options`(
    `id` bigint NOT NULL,
    `search_id` bigint NOT NULL,
    `title` varchar(256) NOT NULL DEFAULT '',
    `link` varchar(128) NOT NULL DEFAULT '',
    `source` varchar(64) NOT NULL DEFAULT '',
    `source_icon` varchar(128) NOT NULL DEFAULT '',
    `in_stock` tinyint NOT NULL DEFAULT 0,
    `price` int NOT NULL DEFAULT 0,
    `currency` varchar(8) NOT NULL DEFAULT '',
    `rank` int NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    KEY `searches_options_search_id_idx` (`search_id`),
    CONSTRAINT `searches_options_search_id_fkey` FOREIGN KEY (`search_id`) REFERENCES `people_searches` (`id`) ON DELETE CASCADE
);

CREATE TABLE `templates_emails` (
  `id` bigint NOT NULL,
  `name` varchar(64) NOT NULL,
  `body` text NOT NULL,
  `subject` text NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unq_templates_emails_name` (`name`)
);

ALTER TABLE `searches_options` RENAME COLUMN `rank` TO `display_order`;
ALTER TABLE `searches_options` ADD COLUMN `image` varchar(256) NOT NULL DEFAULT '';