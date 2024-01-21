CREATE TABLE `group`
(
    `id`           int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '群ID',
    `creator_id`   varchar(50)  NOT NULL COMMENT '创建者ID(群主ID)',
    `type`         tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '群类型[1:普通群;2:企业群;]',
    `group_name`   varchar(30) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '群名称',
    `profile`      varchar(100)                      NOT NULL DEFAULT '' COMMENT '群介绍',
    `avatar`       varchar(255)                      NOT NULL DEFAULT '' COMMENT '群头像',
    `max_num`      smallint(5) unsigned NOT NULL DEFAULT '200' COMMENT '最大群成员数量',
    `is_overt`     tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否公开可见[0:否;1:是;]',
    `is_mute`      tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否全员禁言 [0:否;1:是;]，提示:不包含群主或管理员',
    `is_apply`     tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '是否需要审批[0:否;1:是]',
    `is_allow_invite` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '是否允许邀请[0:否;1:是],不包括群主和管理员',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime                                   DEFAULT NULL COMMENT '解散时间',
    PRIMARY KEY (`id`),
    KEY            `idx_created_at` (`created_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户聊天群';