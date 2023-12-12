CREATE TABLE `talk_records`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '聊天记录ID',
    `msg_id`      varchar(50) NOT NULL DEFAULT '',
    `sequence`    int(10) unsigned NOT NULL DEFAULT '0',
    `talk_type`   tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '对话类型[1:私信;2:群聊;]',
    `msg_type`    int(11) unsigned NOT NULL DEFAULT '1' COMMENT '消息类型[1:文本消息;2:文件消息;3:会话消息;4:代码消息;5:投票消息;6:群公告;7:好友申请;8:登录通知;9:入群消息/退群消息;]',
    `user_id`     varchar(50) NOT NULL DEFAULT '0' COMMENT '发送者ID（0:代表系统消息 >0: 用户ID）',
    `scope` tinyint(4) unsigned NOT NULL COMMENT '接收者类型(0:私聊;1:群聊)'
    `receiver_id` varchar(50) unsigned NOT NULL DEFAULT '0' COMMENT '接收者ID（用户ID）',
    `group_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '接收群ID',
    `is_revoke`   tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否撤回消息[0:否;1:是;]',
    `is_mark`     tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否重要消息[0:否;1:是;]',
    `is_read`     tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否已读[0:否;1:是;]',
    `quote_id`    varchar(50) NOT NULL COMMENT '引用消息ID',
    `content`     text CHARACTER SET utf8mb4 COMMENT '文本消息 {@nickname@}',
    `extra`       json        NOT NULL COMMENT '消息扩展信息',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_at` datetime NOT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_id_receiver_id` (`user_id`,`receiver_id`,`sequence`) USING BTREE,
    UNIQUE KEY `un_msgid` (`msg_id`) USING HASH
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户聊天记录表';