CREATE TABLE `group_apply`
(
    `id`         int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `group_id`   int(11) unsigned NOT NULL DEFAULT '0' COMMENT '群组ID',
    `user_id`    int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
    `status`     int(11) NOT NULL DEFAULT '1' COMMENT '申请状态',
    `remark`     varchar(255) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '备注信息',
    `reason`     varchar(255)                       NOT NULL DEFAULT '' COMMENT '拒绝原因',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY          `idx_group_id_user_id` (`group_id`,`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='群聊成员';