CREATE TABLE `user_account` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `type` int(11) NOT NULL COMMENT '账户类型:0个人,1企业',
  `alias` varchar(50) NOT NULL COMMENT '用户名',
  `mobile` varchar(50) NOT NULL DEFAULT '' COMMENT '绑定手机',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '绑定邮箱',
  `password` varchar(200) NOT NULL COMMENT '账户密码',
  `paycode` varchar(200) NOT NULL DEFAULT '' COMMENT '支付密码',
  `coin_address` varchar(100) NOT NULL DEFAULT '' COMMENT '链账户地址',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '0普通,1实名中,2实名失败,3实名成功',
  PRIMARY KEY (`id`),
  UNIQUE KEY `alias` (`alias`),
  KEY `ix_mobile` (`mobile`),
  KEY `ix_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=10009 DEFAULT CHARSET=utf8;
