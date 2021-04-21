SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;


DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`         bigint(20)  NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `mobile`     varchar(32) NOT NULL COMMENT '手机号,唯一标识 ',
    `username`   varchar(64) NOT NULL COMMENT '用户名称',
    `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `mobile_idx` (`mobile`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 171
  DEFAULT CHARSET = utf8
  ROW_FORMAT = DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;