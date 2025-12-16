CREATE DATABASE IF NOT EXISTS test CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE test;

CREATE TABLE IF NOT EXISTS `sale` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
  `merchant_id` bigint unsigned NOT NULL COMMENT 'Merchant ID',
  `store_id` bigint unsigned NOT NULL COMMENT 'Store ID',
  `cashier_terminal_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Cashier Terminal ID',
  `cashier_id` bigint unsigned NOT NULL COMMENT 'Cashier ID',
  `cashier_employee_no` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Cashier Employee No.',
  `cashier_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Cashier Name',
  `order_total_amount` decimal(18,2) NOT NULL DEFAULT '0.00' COMMENT 'Order Total Amount',
  `order_discount_amount` decimal(18,2) NOT NULL DEFAULT '0.00' COMMENT 'Order Discount Amount',
  `order_paid_amount` decimal(18,2) NOT NULL DEFAULT '0.00' COMMENT 'Order Paid Amount',
  `order_no` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Order No., Unique Identifier',
  `order_status` tinyint NOT NULL DEFAULT '0' COMMENT 'Order Status: 0-Pending Payment, 1-Paid, 2-Cancelled, 3-Completed',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Order Creation Time',
  `ext_json` json DEFAULT NULL COMMENT 'JSON extension field, used to store dynamic attributes',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Sales Table';
