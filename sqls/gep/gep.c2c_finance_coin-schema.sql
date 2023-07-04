CREATE TABLE `c2c_finance_coin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `uid` int(11) DEFAULT '0' COMMENT '用户id',
  `type` tinyint(4) DEFAULT '0' COMMENT '资产分类，1：数字源货币',
  `sum_coin` decimal(32,8) DEFAULT '0.00000000' COMMENT '总的数字源币',
  `freeze_coin` decimal(32,8) DEFAULT '0.00000000' COMMENT '冻结的数字源币',
  `balance_coin` decimal(32,8) DEFAULT '0.00000000' COMMENT '余额的数字源币',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='挂单表';
