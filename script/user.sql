CREATE TABLE `user` (
    `id` varchar(50) NOT NULL COMMMENt '用户id',
    `account` varchar(20) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '账号',
    `username` varchar(20) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '用户名',
    `password` varchar(255) NOT NULL COMMENT '用户密码',
    `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '用户头像',
    `gender` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '用户性别[0:未知;1:男;2:女]',
    `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
    `email` varchar(30) NOT NULL DEFAULT '' COMMENT '用户邮箱',
    `introduce` varchar(255) NOT NULL DEFAULT '' COMMENT '用户个人介绍',
    `is_robot`   tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否机器人[0:否;1:是;]',
    `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态: 1:正常 2:停用',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `idx_account` (`account`) USING BTREE,
    UNIQUE KEY `idx_email` (`email`) USING BTREE,
    UNIQUE KEY `idx_mobile` (`mobile`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='用户表';