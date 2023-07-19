CREATE TABLE `withdraw_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `uid` int(11) DEFAULT NULL,
  `txid` varchar(70) DEFAULT NULL,
  `from_address` varchar(50) NOT NULL,
  `to_address` varchar(50) NOT NULL,
  `total_value` decimal(32,8) DEFAULT NULL,
  `actual_value` decimal(32,8) DEFAULT NULL,
  `total_fee` decimal(32,8) DEFAULT NULL,
  `tx_fee` decimal(32,8) DEFAULT NULL,
  `balance_fee` decimal(32,8) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=latin1;
