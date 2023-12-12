CREATE TABLE `article_annex`
(
    `id`            int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '文件ID',
    `user_id`       varchar(50) NOT NULL COMMENT '上传文件的用户ID',
    `article_id`    int(11) unsigned NOT NULL DEFAULT '1' COMMENT '笔记ID',
    `drive`         tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '文件驱动[1:local;2:cos;]',
    `suffix`        varchar(10)  NOT NULL DEFAULT '' COMMENT '文件后缀名',
    `size`          bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '文件大小',
    `path`          varchar(500) NOT NULL DEFAULT '' COMMENT '文件地址（相对地址）',
    `original_name` varchar(100) NOT NULL DEFAULT '' COMMENT '原文件名',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`    datetime              DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY             `idx_userid_articleid` (`user_id`,`article_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章附件信息表';