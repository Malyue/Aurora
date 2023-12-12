CREATE TABLE `article_tag`
(
    `id`         int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '文章分类ID',
    `user_id`       varchar(50) NOT NULL COMMENT '用户id',
    `tag_name`   varchar(20) NOT NULL DEFAULT '' COMMENT '标签名',
    `sort`       tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`    datetime              DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY          `idx_userid` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章标签表';