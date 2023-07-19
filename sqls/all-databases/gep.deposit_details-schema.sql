CREATE TABLE `deposit_details` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `uid` int(10) unsigned DEFAULT NULL,
  `txid` varchar(70) DEFAULT NULL,
  `from_address` varchar(50) DEFAULT NULL,
  `to_address` varchar(50) DEFAULT NULL,
  `value` decimal(32,8) DEFAULT '0.00000000',
  `height` int(10) unsigned DEFAULT '0',
  `status` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
