CREATE TABLE `user_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `real_name` varchar(256) DEFAULT '' COMMENT '用户名',
  `alia_name` varchar(128) DEFAULT '' COMMENT '微信昵称',
  `phone` varchar(20) DEFAULT '' COMMENT '电话号码',
  `company` varchar(256) DEFAULT '' COMMENT '公司或者学校',
  `address` varchar(256) DEFAULT '' COMMENT '家庭住址',
  `wei_id` varchar(128) DEFAULT '' COMMENT '微信号',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='用户信息';
