CREATE TABLE `activity` (
  `activity_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '活动Id',
  `activity_name` varchar(50) NOT NULL DEFAULT '' COMMENT '活动名称',
  `product_id` int(11) unsigned NOT NULL COMMENT '商品Id',
  `start_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '活动开始时间',
  `end_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '活动结束时间',
  `total` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '商品数量',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '活动状态',
  `sec_speed` int(5) unsigned NOT NULL DEFAULT '0' COMMENT '每秒限制多少个商品售出',
  `buy_limit` int(5) unsigned NOT NULL COMMENT '购买限制',
  `buy_rate` decimal(2,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '购买限制',
  PRIMARY KEY (`activity_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COMMENT='@活动数据表';
