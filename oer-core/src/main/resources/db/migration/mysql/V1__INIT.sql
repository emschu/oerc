CREATE TABLE IF NOT EXISTS `channel` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `adapter_family` int(11) DEFAULT NULL,
  `channel_key` int(11) DEFAULT NULL,
  `home_page` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `technical_id` varchar(12) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UK_channel_key_index` (`channel_key`),
  UNIQUE KEY `UK_technical_id_index` (`technical_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `image_link` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL,
  `url` varchar(2500) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `program_entry` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `adapter_family` int(11) NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `description` text DEFAULT NULL,
  `duration_in_minutes` int(11) DEFAULT NULL,
  `end_date_time` datetime DEFAULT NULL,
  `home_page` varchar(1000) CHARACTER SET utf8mb4 DEFAULT NULL,
  `start_date_time` datetime DEFAULT NULL,
  `technical_id` varchar(36) CHARACTER SET utf8mb4 DEFAULT NULL,
  `title` varchar(1000) CHARACTER SET utf8mb4 DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `url` varchar(1000) CHARACTER SET utf8mb4 DEFAULT NULL,
  `channel_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UK_technical_id_index` (`technical_id`),
  UNIQUE KEY `UK_technical_id_adapter_family` (`technical_id`,`adapter_family`),
  KEY `FK_channel_relation` (`channel_id`),
  CONSTRAINT `FK_channel_relation` FOREIGN KEY (`channel_id`) REFERENCES `channel` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `program_entry_image_links` (
  `program_entry_id` bigint(20) NOT NULL,
  `image_links_id` bigint(20) NOT NULL,
  KEY `FK_image_link_relation` (`image_links_id`),
  KEY `FK_program_entry_relation` (`program_entry_id`),
  CONSTRAINT `FK_program_entry_relation` FOREIGN KEY (`program_entry_id`) REFERENCES `program_entry` (`id`),
  CONSTRAINT `FK_image_link_relation` FOREIGN KEY (`image_links_id`) REFERENCES `image_link` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `tag` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL,
  `tag_name` varchar(120) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UK_tag_name_index` (`tag_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `program_entry_tags` (
  `program_entry_id` bigint(20) NOT NULL,
  `tags_id` bigint(20) NOT NULL,
  KEY `FK_tag_id_relation` (`tags_id`),
  KEY `FK_program_entry_tag_relation` (`program_entry_id`),
  CONSTRAINT `FK_tag_id_relation` FOREIGN KEY (`tags_id`) REFERENCES `tag` (`id`),
  CONSTRAINT `FK_program_entry_tag_relation` FOREIGN KEY (`program_entry_id`) REFERENCES `program_entry` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `tv_show` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `adapter_family` int(11) NOT NULL,
  `additional_id` varchar(1000) CHARACTER SET utf8mb4 DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `home_page` varchar(1500) CHARACTER SET utf8mb4 DEFAULT NULL,
  `image_url` varchar(1500) CHARACTER SET utf8mb4 DEFAULT NULL,
  `technical_id` varchar(32) CHARACTER SET utf8mb4 NOT NULL,
  `title` varchar(1000) CHARACTER SET utf8mb4 NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  `url` varchar(1500) CHARACTER SET utf8mb4 DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UK_adapter_tid_index` (`adapter_family`,`technical_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `tv_show_related_program_entries` (
  `tv_show_id` bigint(20) NOT NULL,
  `related_program_entry_id` bigint(20) NOT NULL,
  KEY `FK_tv_program_entry_relation` (`related_program_entry_id`),
  KEY `FK_tv_show_relation` (`tv_show_id`),
  CONSTRAINT `FK_tv_program_entry_relation` FOREIGN KEY (`related_program_entry_id`) REFERENCES `program_entry` (`id`),
  CONSTRAINT `FK_tv_show_relation` FOREIGN KEY (`tv_show_id`) REFERENCES `tv_show` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `tv_show_tags` (
  `tv_show_id` bigint(20) NOT NULL,
  `tag_id` bigint(20) NOT NULL,
  KEY `FK_tag_show_relation` (`tag_id`),
  KEY `FK_tv_show_tag_relation` (`tv_show_id`),
  CONSTRAINT `FK_tv_show_tag_relation` FOREIGN KEY (`tv_show_id`) REFERENCES `tv_show` (`id`),
  CONSTRAINT `FK_tag_show_relation` FOREIGN KEY (`tag_id`) REFERENCES `tag` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
