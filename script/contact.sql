CREATE TABLE `contact`
(
    `id`         int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '关系ID',
    `user_id`    varchar(50) NOT NULL COMMENT '用户id',
    `friend_id`  varchar(50) NOT NULL COMMENT '好友id',
    `remark`     varchar(20) NOT NULL DEFAULT '' COMMENT '好友的备注',
    `status`     tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '好友状态 [0:否;1:是]',
    `group_id`   int(11) unsigned NOT NULL DEFAULT '0' COMMENT '好友分组',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY          `idx_user1_user2` (`user_id`,`friend_id`) USING BTREE,
    KEY          `idx_user2_user1` (`friend_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户好友关系表';