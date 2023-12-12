CREATE TABLE `contact_apply`
(
    `id`         int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '申请ID',
    `user_id`    varchar(50) NOT NULL COMMENT '申请人ID',
    `friend_id`  varchar(50) NOT NULL  COMMENT '被申请人',
    `remark`     varchar(50) NOT NULL DEFAULT '' COMMENT '申请备注',
    `apply_at` datetime    NOT NULL COMMENT '申请时间',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY          `idx_user_id` (`user_id`) USING BTREE,
    KEY          `idx_friend_id` (`friend_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户添加好友申请表';