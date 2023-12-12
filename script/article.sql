CREATE TABLE `article`(
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '文章ID',
    `user_id` varchar(50) NOT NULL COMMENT '作者id',
    `group_id` int(20) NOT NULL COMMENT '组id',
    `class_id` int(11) unsigned NOT NULL COMMENT '分类id',
    `tags_id` varchar(50) NOT NULL DEFAULT '' COMMENT '关联标签',
    `title` varchar(50) NOT NULL DEFAULT '' COMMENT '文章标题',
    `abstract` varchar(200) NOT NULL DEFAULT COMMENT '文章摘要',
    `image` varchar(255) NOT NULL DEFAULT '' COMMENT '文章首图',
    `is_asterisk` tinyint(3) unsigned NOT NULL DEFAULT 0 COMMENT '是否为星标文章[0:否;1:是]',
    `scope_type` tinyint(3) unsigned NOT NULL DEFAULT 0 COMMENT '可见范围[0:私人文档;1:群组;2:所有人]'
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_userid_classid` (`user_id`,`class_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文档表';