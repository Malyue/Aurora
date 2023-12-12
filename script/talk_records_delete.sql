CREATE TABLE `talk_records_delete`
(
    `id`         int(11) unsigned NOT NULL AUTO_INCREMENT,
    `record_id`  int(11) unsigned NOT NULL DEFAULT '0' COMMENT '聊天记录ID',
    `user_id`    int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_record_user_id` (`record_id`,`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='聊天记录删除记录表';