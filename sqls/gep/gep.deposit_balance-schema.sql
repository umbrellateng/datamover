CREATE TABLE `deposit_balance` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) DEFAULT '0',
  `txid` varchar(70) DEFAULT NULL,
  `from_address` varchar(50) DEFAULT NULL,
  `to_address` varchar(50) DEFAULT NULL,
  `value` decimal(50,0) DEFAULT NULL,
  `height` int(11) DEFAULT NULL,
  `total_status` int(11) DEFAULT '0',
  `transfer_status` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=384 DEFAULT CHARSET=latin1;
