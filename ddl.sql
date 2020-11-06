CREATE TABLE `g_account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `user_id` bigint(20) NOT NULL,
  `currency` varchar(255) NOT NULL,
  `hold` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `available` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uid_currency` (`user_id`,`currency`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `g_bill` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `user_id` bigint(20) NOT NULL,
  `currency` varchar(255) NOT NULL,
  `available` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `hold` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `type` varchar(255) NOT NULL,
  `settled` tinyint(1) NOT NULL DEFAULT '0',
  `notes` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_gsoci` (`user_id`,`currency`,`settled`,`id`),
  KEY `idx_s` (`settled`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `g_config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `key` varchar(255) NOT NULL,
  `value` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `g_fill` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `trade_id` bigint(20) NOT NULL DEFAULT '0',
  `order_id` bigint(20) NOT NULL DEFAULT '0',
  `product_id` varchar(255) NOT NULL,
  `size` decimal(32,16) NOT NULL,
  `price` decimal(32,16) NOT NULL,
  `funds` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `fee` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `liquidity` varchar(255) NOT NULL,
  `settled` tinyint(1) NOT NULL DEFAULT '0',
  `side` varchar(255) NOT NULL,
  `done` tinyint(1) NOT NULL DEFAULT '0',
  `done_reason` varchar(255) NOT NULL,
  `message_seq` bigint(20) NOT NULL,
  `log_offset` bigint(20) NOT NULL DEFAULT '0',
  `log_seq` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `o_m` (`order_id`,`message_seq`),
  KEY `idx_gsoi` (`order_id`,`settled`,`id`),
  KEY `idx_si` (`settled`,`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `g_order` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `product_id` varchar(255) NOT NULL,
  `user_id` bigint(20) NOT NULL,
  `size` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `funds` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `filled_size` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `executed_value` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `price` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `fill_fees` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `type` varchar(255) NOT NULL,
  `side` varchar(255) NOT NULL,
  `time_in_force` varchar(255) DEFAULT NULL,
  `status` varchar(255) NOT NULL,
  `settled` tinyint(1) NOT NULL DEFAULT '0',
  `client_oid` varchar(32) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_uspsi` (`user_id`,`product_id`,`status`,`side`,`id`),
  KEY `idx_uid_coid` (`user_id`,`client_oid`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `g_product` (
  `id` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `base_currency` varchar(255) NOT NULL,
  `quote_currency` varchar(255) NOT NULL,
  `base_min_size` decimal(32,16) NOT NULL,
  `base_max_size` decimal(32,16) NOT NULL,
  `base_scale` int(11) NOT NULL,
  `quote_scale` int(11) NOT NULL,
  `quote_increment` double NOT NULL,
  `quote_min_size` decimal(32,16) NOT NULL,
  `quote_max_size` decimal(32,16) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `g_tick` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `product_id` varchar(255) NOT NULL,
  `granularity` bigint(20) NOT NULL,
  `time` bigint(20) NOT NULL,
  `open` decimal(32,16) NOT NULL,
  `high` decimal(32,16) NOT NULL,
  `low` decimal(32,16) NOT NULL,
  `close` decimal(32,16) NOT NULL,
  `volume` decimal(32,16) NOT NULL,
  `log_offset` bigint(20) NOT NULL DEFAULT '0',
  `log_seq` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `p_g_t` (`product_id`,`granularity`,`time`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `g_trade` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `product_id` varchar(255) NOT NULL,
  `taker_order_id` bigint(20) NOT NULL,
  `maker_order_id` bigint(20) NOT NULL,
  `price` decimal(32,16) NOT NULL,
  `size` decimal(32,16) NOT NULL,
  `side` varchar(255) NOT NULL,
  `time` timestamp NOT NULL,
  `log_offset` bigint(20) NOT NULL DEFAULT '0',
  `log_seq` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `g_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


insert into `g_product`(`id`,`created_at`,`updated_at`,`base_currency`,`quote_currency`,`base_min_size`,`base_max_size`,`base_scale`,`quote_scale`,`quote_increment`,`quote_min_size`,`quote_max_size`) values
('BITC-USDT',null,null,'BITC','USDT',0.0001000000000000,10000.0000000000000000,4,2,0.01,0E-16,0E-16),
('BCH-USDT',null,null,'BCH','USDT',0.0001000000000000,10000.0000000000000000,4,2,0.01,0E-16,0E-16),
('BTC-USDT',null,null,'BTC','USDT',0.0000010000000000,10000000.0000000000000000,6,2,0.01,0E-16,0E-16),
('EOS-USDT',null,null,'EOS','USDT',0.0001000000000000,1000.0000000000000000,4,3,0,0E-16,0E-16),
('ETH-USDT',null,null,'ETH','USDT',0.0001000000000000,10000.0000000000000000,4,2,0.01,0E-16,0E-16),
('LTC-USDT',null,null,'LTC','USDT',0.0001000000000000,1000.0000000000000000,4,2,0.01,0E-16,0E-16);


----- BitEx Tables
CREATE TABLE `g_address` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `address` varchar(255) NOT NULL,
  `public_key` varchar(255) NOT NULL,
  `private_key` varchar(255) NOT NULL,
  `mnemonic` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_address` (`address`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `g_account_asset` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `user_id` bigint(20) NOT NULL,
  `currency` varchar(255) NOT NULL,
  `hold` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `available` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uid_currency` (`user_id`,`currency`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `g_account_pool` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `user_id` bigint(20) NOT NULL,
  `currency` varchar(255) NOT NULL,
  `hold` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `available` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uid_currency` (`user_id`,`currency`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `g_account_shop` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `user_id` bigint(20) NOT NULL,
  `currency` varchar(255) NOT NULL,
  `hold` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `available` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uid_currency` (`user_id`,`currency`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


insert into `g_config`(`created_at`,`updated_at`,`key`,`value`) values
(NOW(),NOW(),'ACTIVE_NUMBER','0.1'),
(NOW(),NOW(),'ISSUE_NUMBER','100000'),
(NOW(),NOW(),'ENERGY_NUMBER','5000'),
(NOW(),NOW(),'HOLD_NUMBER','8000'),
(NOW(),NOW(),'PROMOTE_NUMBER','4000');

