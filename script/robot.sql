CREATE TABLE `robot`
(
    `id`         int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '机器人ID',
    `user_id`    varchar(50) NOT NULL COMMENT '关联用户ID',
    `robot_name` varchar(20)  NOT NULL DEFAULT '' COMMENT '机器人名称',
    `describe`   varchar(255) NOT NULL DEFAULT '' COMMENT '描述信息',
    `logo`       varchar(255) NOT NULL DEFAULT '' COMMENT '机器人logo',
    `is_talk`    tinyint(4) NOT NULL DEFAULT '0' COMMENT '可发送消息[0:否;1:是;]',
    `status`     tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '状态[-1:已删除;0:正常;1:已禁用;]',
    `type`       tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '机器人类型',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_type` (`type`) USING HASH,
    UNIQUE KEY `uk_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='聊天机器人表';