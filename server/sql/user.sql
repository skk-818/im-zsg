-- 用户表
CREATE TABLE `im_users` (
  `user_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码(加密后)',
  `nickname` varchar(100) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar` varchar(500) NOT NULL DEFAULT '' COMMENT '头像URL',
  `email` varchar(100) NOT NULL DEFAULT '' COMMENT '邮箱',
  `phone` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
  `gender` tinyint(1) NOT NULL DEFAULT 0 COMMENT '性别: 0-未知, 1-男, 2-女',
  `birthday` date NOT NULL DEFAULT '1900-01-01' COMMENT '生日',
  `signature` varchar(500) NOT NULL DEFAULT '' COMMENT '个性签名',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '账号状态: 0-禁用, 1-正常',
  `last_login_time` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(45) NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '删除标记: 0-未删除, 1-已删除',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_email` (`email`),
  KEY `idx_phone` (`phone`),
  KEY `idx_status` (`status`),
  KEY `idx_delete_flag` (`delete_flag`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 用户设备表 (用于多设备登录管理)
CREATE TABLE `im_user_devices` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `device_id` varchar(100) NOT NULL COMMENT '设备ID',
  `device_type` varchar(20) NOT NULL COMMENT '设备类型: ios, android, web, desktop',
  `device_name` varchar(100) NOT NULL DEFAULT '' COMMENT '设备名称',
  `push_token` varchar(500) NOT NULL DEFAULT '' COMMENT '推送token',
  `is_online` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否在线: 0-离线, 1-在线',
  `last_active_time` timestamp NULL DEFAULT NULL COMMENT '最后活跃时间',
  `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '删除标记: 0-未删除, 1-已删除',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_device` (`user_id`, `device_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_is_online` (`is_online`),
  KEY `idx_delete_flag` (`delete_flag`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户设备表';

-- 用户token表 (用于JWT token管理)
CREATE TABLE `im_user_tokens` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `device_id` varchar(100) NOT NULL COMMENT '设备ID',
  `token` varchar(500) NOT NULL COMMENT 'JWT token',
  `refresh_token` varchar(500) NOT NULL DEFAULT '' COMMENT '刷新token',
  `expires_at` timestamp NOT NULL COMMENT 'token过期时间',
  `is_revoked` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已撤销: 0-有效, 1-已撤销',
  `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '删除标记: 0-未删除, 1-已删除',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_token` (`token`(255)),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_expires_at` (`expires_at`),
  KEY `idx_is_revoked` (`is_revoked`),
  KEY `idx_delete_flag` (`delete_flag`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户token表';

-- 好友关系表
CREATE TABLE `im_user_friends` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `friend_id` bigint(20) unsigned NOT NULL COMMENT '好友ID',
  `remark` varchar(100) NOT NULL DEFAULT '' COMMENT '好友备注',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '关系状态: 0-已删除, 1-正常, 2-拉黑',
  `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '删除标记: 0-未删除, 1-已删除',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_friend` (`user_id`, `friend_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_friend_id` (`friend_id`),
  KEY `idx_status` (`status`),
  KEY `idx_delete_flag` (`delete_flag`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='好友关系表';

-- 好友申请表
CREATE TABLE `im_friend_requests` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `from_user_id` bigint(20) unsigned NOT NULL COMMENT '申请人ID',
  `to_user_id` bigint(20) unsigned NOT NULL COMMENT '被申请人ID',
  `message` varchar(500) NOT NULL DEFAULT '' COMMENT '申请消息',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '申请状态: 0-待处理, 1-已同意, 2-已拒绝, 3-已过期',
  `processed_at` timestamp NULL DEFAULT NULL COMMENT '处理时间',
  `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '删除标记: 0-未删除, 1-已删除',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_from_user` (`from_user_id`),
  KEY `idx_to_user` (`to_user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_delete_flag` (`delete_flag`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='好友申请表';

-- 群组表
CREATE TABLE `im_groups` (
  `group_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '群组ID',
  `group_name` varchar(100) NOT NULL COMMENT '群组名称',
  `group_avatar` varchar(500) NOT NULL DEFAULT '' COMMENT '群组头像',
  `description` varchar(500) NOT NULL DEFAULT '' COMMENT '群组描述',
  `owner_id` bigint(20) unsigned NOT NULL COMMENT '群主ID',
  `max_members` int(11) NOT NULL DEFAULT 500 COMMENT '最大成员数',
  `member_count` int(11) NOT NULL DEFAULT 0 COMMENT '当前成员数',
  `is_public` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否公开: 0-私有, 1-公开',
  `join_approval` tinyint(1) NOT NULL DEFAULT 0 COMMENT '加群是否需要审批: 0-不需要, 1-需要',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '群组状态: 0-解散, 1-正常',
  `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '删除标记: 0-未删除, 1-已删除',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`group_id`),
  KEY `idx_owner_id` (`owner_id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_public` (`is_public`),
  KEY `idx_delete_flag` (`delete_flag`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群组表';

-- 群组成员表
CREATE TABLE `im_group_members` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `group_id` bigint(20) unsigned NOT NULL COMMENT '群组ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `role` tinyint(1) NOT NULL DEFAULT 0 COMMENT '角色: 0-普通成员, 1-管理员, 2-群主',
  `nickname` varchar(100) DEFAULT NULL COMMENT '群内昵称',
  `mute_until` timestamp NULL DEFAULT NULL COMMENT '禁言到期时间',
  `delete_flag` tinyint(1) NOT NULL DEFAULT 0 COMMENT '删除标记: 0-未删除, 1-已删除',
  `joined_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_user` (`group_id`, `user_id`),
  KEY `idx_group_id` (`group_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role` (`role`),
  KEY `idx_delete_flag` (`delete_flag`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群组成员表';

-- 插入测试数据
INSERT INTO `im_users` (`username`, `password`, `nickname`, `email`, `status`) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKVjzieMwkOBX.bVARjAWpOqOGka', '管理员', 'admin@example.com', 1),
('user1', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKVjzieMwkOBX.bVARjAWpOqOGka', '用户1', 'user1@example.com', 1),
('user2', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKVjzieMwkOBX.bVARjAWpOqOGka', '用户2', 'user2@example.com', 1);