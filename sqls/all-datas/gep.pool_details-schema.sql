CREATE TABLE `pool_details` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `txid` varchar(70) DEFAULT NULL,
  `from_address` varchar(50) DEFAULT NULL,
  `to_address` varchar(50) DEFAULT NULL,
  `value` decimal(32,8) DEFAULT NULL,
  `height` int(10) unsigned DEFAULT NULL,
  `is_withdraw` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `txid` (`txid`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=latin1;
