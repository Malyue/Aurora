CREATE TABLE `group_member`
(
    `id`            int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `group_id`      int(11) unsigned NOT NULL DEFAULT '0' COMMENT '群组ID',
    `user_id`      varchar(50) NOT NULL COMMENT '用户ID',
    `role`        tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '成员属性[0:普通成员;1:管理员;2:群主;]',
    `user_card`     varchar(20) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '群名片',
    `is_mute`       tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否禁言[0:否;1:是;]',
    `min_record_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '可查看最大消息ID',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '退群时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_group_id_user_id` (`group_id`,`user_id`) USING BTREE,
    KEY             `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='群聊成员';