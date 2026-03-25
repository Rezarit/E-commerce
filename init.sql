-- init.sql

-- 用户表 (users)
CREATE TABLE `users` (
                         `user_id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
                         `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
                         `nickname` varchar(255) DEFAULT NULL COMMENT '昵称',
                         `introduction` text DEFAULT NULL COMMENT '个人简介',
                         `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
                         `qq` varchar(20) DEFAULT NULL COMMENT 'QQ号',
                         `gender` varchar(10) DEFAULT NULL COMMENT '性别',
                         `email` varchar(255) DEFAULT NULL COMMENT '邮箱',
                         `birthday` varchar(20) DEFAULT NULL COMMENT '生日',
                         `username` varchar(255) NOT NULL UNIQUE COMMENT '用户名',
                         `password` varchar(255) NOT NULL COMMENT '密码',
                         `role` int NOT NULL DEFAULT '0' COMMENT '角色 (0: 普通用户, 1: 商家)',
                         `created_at` datetime(3) DEFAULT NULL,
                         `updated_at` datetime(3) DEFAULT NULL,
                         PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 商家表 (merchants)
CREATE TABLE `merchants` (
                             `merchant_id` bigint NOT NULL AUTO_INCREMENT COMMENT '商家ID',
                             `user_id` bigint NOT NULL UNIQUE COMMENT '关联的用户ID',
                             `merchant_name` varchar(100) NOT NULL UNIQUE COMMENT '商家名称',
                             `business_license` varchar(50) NOT NULL COMMENT '营业执照',
                             `contact_phone` varchar(20) NOT NULL COMMENT '联系电话',
                             `address` text NOT NULL COMMENT '地址',
                             `shop_description` text COMMENT '店铺描述',
                             `status` varchar(20) DEFAULT 'active' COMMENT '商家状态',
                             `create_time` datetime(3) DEFAULT NULL,
                             `update_time` datetime(3) DEFAULT NULL,
                             PRIMARY KEY (`merchant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 商家入驻申请表 (merchant_applications)
CREATE TABLE `merchant_applications` (
                                         `application_id` bigint NOT NULL AUTO_INCREMENT,
                                         `user_id` bigint NOT NULL,
                                         `merchant_name` varchar(100) NOT NULL,
                                         `business_license` varchar(50) NOT NULL,
                                         `contact_phone` varchar(20) NOT NULL,
                                         `address` text NOT NULL,
                                         `status` varchar(20) DEFAULT 'pending',
                                         `apply_time` datetime(3) DEFAULT NULL,
                                         `audit_time` datetime(3) DEFAULT NULL,
                                         `audit_admin` varchar(50) DEFAULT NULL,
                                         `reject_reason` text DEFAULT NULL,
                                         PRIMARY KEY (`application_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 商品表 (products)
CREATE TABLE `products` (
                            `product_id` bigint NOT NULL AUTO_INCREMENT,
                            `merchant_id` bigint NOT NULL,
                            `product_name` varchar(255) NOT NULL,
                            `description` text NOT NULL,
                            `comment_num` int DEFAULT '0',
                            `price` decimal(10,2) DEFAULT '0.00',
                            `stock` int DEFAULT '0',
                            `cover` varchar(255) NOT NULL,
                            `publish_time` datetime(3) DEFAULT NULL,
                            `link` varchar(255) NOT NULL,
                            PRIMARY KEY (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 购物车表 (carts)
CREATE TABLE `carts` (
                         `cart_id` bigint NOT NULL AUTO_INCREMENT,
                         `user_id` bigint DEFAULT NULL,
                         `product_id` bigint DEFAULT NULL,
                         `quantity` int DEFAULT '1',
                         PRIMARY KEY (`cart_id`),
                         KEY `idx_carts_user_id` (`user_id`),
                         KEY `idx_carts_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 订单表 (orders)
CREATE TABLE `orders` (
                          `order_id` bigint NOT NULL AUTO_INCREMENT,
                          `user_id` bigint DEFAULT NULL,
                          `address` text NOT NULL,
                          `total` decimal(10,2) DEFAULT '0.00',
                          `status` varchar(255) DEFAULT 'pending',
                          `created_at` datetime(3) DEFAULT NULL,
                          PRIMARY KEY (`order_id`),
                          KEY `idx_orders_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 订单项表 (order_items)
CREATE TABLE `order_items` (
                               `order_item_id` bigint NOT NULL AUTO_INCREMENT,
                               `order_id` bigint DEFAULT NULL,
                               `product_id` bigint NOT NULL,
                               `product_name` varchar(255) NOT NULL,
                               `quantity` int NOT NULL,
                               `price` decimal(10,2) NOT NULL,
                               PRIMARY KEY (`order_item_id`),
                               KEY `idx_order_items_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;