CREATE TABLE `metadata` (
  `object_id` char(26) NOT NULL COMMENT '对象id',
  `parent_id` char(26) NOT NULL COMMENT '父对象id',
  `name` varchar(64) NOT NULL COMMENT '名称',
  `basic_attr` int(11) DEFAULT NULL COMMENT '元数据类型。1: 文件夹 2:目录',
  PRIMARY KEY (`object_id`),
  KEY `idx_parentId` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='元数据表';

CREATE TABLE `metadata_closure` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `ancestor` char(26) NOT NULL COMMENT '祖先',
  `descendant` char(26) NOT NULL COMMENT '后代',
  `depth` int(11) NOT NULL COMMENT '层级深度，从0开始',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='元数据路径闭包表';