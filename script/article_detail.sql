CREATE TABLE `article_detail`
(
    `id`         int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '文章详情ID',
    `article_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '文章ID',
    `md_content` longtext CHARACTER SET utf8mb4 NOT NULL COMMENT 'Markdown 内容',
    `content`    longtext CHARACTER SET utf8mb4 NOT NULL COMMENT 'Markdown 解析HTML内容',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`    datetime              DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_article_id` (`article_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章详情表';