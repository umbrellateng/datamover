CREATE TABLE `user_temperature` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `alia_name` varchar(128) DEFAULT '' COMMENT '用户昵称',
  `temperature` decimal(32,4) DEFAULT '0.0000' COMMENT '用户体温',
  `health_status` int(11) DEFAULT '-1',
  `others` varchar(1024) DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COMMENT='用户温度表';
