CREATE TABLE `deposit` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(20) DEFAULT NULL,
  `txid` varchar(70) DEFAULT NULL,
  `from_address` varchar(50) DEFAULT NULL,
  `to_address` varchar(50) DEFAULT NULL,
  `value` decimal(50,0) DEFAULT NULL,
  `height` int(11) DEFAULT NULL,
  `balance_status` int(11) DEFAULT '0',
  `transfer_status` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=26571 DEFAULT CHARSET=latin1;
