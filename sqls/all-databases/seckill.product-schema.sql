CREATE TABLE `product` (
  `product_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '商品Id',
  `product_name` varchar(50) NOT NULL DEFAULT '' COMMENT '商品名称',
  `total` int(5) unsigned NOT NULL DEFAULT '0' COMMENT '商品数量',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '商品状态',
  PRIMARY KEY (`product_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COMMENT='@商品数据表';
