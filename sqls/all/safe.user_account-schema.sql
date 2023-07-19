CREATE TABLE `user_account` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `real_name` varchar(128) DEFAULT '',
  `alia_name` varchar(128) NOT NULL,
  `phone` varchar(20) DEFAULT '',
  `email` varchar(128) DEFAULT '',
  `community` varchar(256) DEFAULT '',
  `building_number` int(11) DEFAULT '0',
  `building_uint` int(11) DEFAULT '0',
  `house_number` int(11) DEFAULT '0',
  `province` varchar(128) DEFAULT '',
  `city` varchar(128) DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_account_alia_name_uindex` (`alia_name`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
