/*
 Navicat Premium Data Transfer

 Source Server         : 短信分销
 Source Server Type    : MySQL
 Source Server Version : 50726
 Source Host           : localhost:3306
 Source Schema         : fusionsms

 Target Server Type    : MySQL
 Target Server Version : 50726
 File Encoding         : 65001

 Date: 29/10/2023 19:31:14
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for exchange_code
-- ----------------------------
DROP TABLE IF EXISTS `exchange_code`;
CREATE TABLE `exchange_code`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(16) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `quota` int(11) NULL DEFAULT NULL,
  `status` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '0' COMMENT '0未使用 1已使用',
  `create_time` datetime(0) NULL DEFAULT NULL,
  `use_time` datetime(0) NULL DEFAULT NULL,
  `use_uid` int(11) NULL DEFAULT NULL COMMENT '使用者UID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `code`(`code`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of exchange_code
-- ----------------------------
INSERT INTO `exchange_code` VALUES (3, 'iaXe2oMSRdSWsRr2', 10, '0', '2023-03-03 22:23:06', NULL, NULL);
INSERT INTO `exchange_code` VALUES (4, '5I4A7fTyfqCnrK47', 10, '0', '2023-03-03 22:23:06', NULL, NULL);
INSERT INTO `exchange_code` VALUES (5, 'AxqviADVil1O0wgW', 10, '0', '2023-03-03 22:23:06', NULL, NULL);
INSERT INTO `exchange_code` VALUES (6, 'mcnDLQ8k3OPFXvJy', 10, '0', '2023-03-03 22:23:06', NULL, NULL);
INSERT INTO `exchange_code` VALUES (7, 'LLhAWo8o0HrKIOQf', 10, '0', '2023-03-03 22:23:06', NULL, NULL);
INSERT INTO `exchange_code` VALUES (8, 'uj7uQj9NyqBjCqhU', 10, '0', '2023-03-03 22:23:06', NULL, NULL);
INSERT INTO `exchange_code` VALUES (9, '5YdljPuqztn1zMcx', 10, '0', '2023-03-03 22:23:06', NULL, NULL);
INSERT INTO `exchange_code` VALUES (10, 'IhTCczpbpaFwGXwu', 10, '0', '2023-03-03 22:23:06', NULL, NULL);
INSERT INTO `exchange_code` VALUES (11, 'SCyui6FXnDjpe12l', 10, '1', '2023-03-03 22:23:06', '2023-09-09 19:31:30', 15);
INSERT INTO `exchange_code` VALUES (12, 'osR4YMGIeards06E', 10, '1', '2023-03-03 22:23:06', '2023-03-03 22:48:10', 12);

-- ----------------------------
-- Table structure for exchange_log
-- ----------------------------
DROP TABLE IF EXISTS `exchange_log`;
CREATE TABLE `exchange_log`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NULL DEFAULT NULL,
  `content` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '',
  `time` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of exchange_log
-- ----------------------------
INSERT INTO `exchange_log` VALUES (1, 12, '兑换额度:10', '2023-03-03 16:29:29');
INSERT INTO `exchange_log` VALUES (2, 12, '兑换额度:10', '2023-03-03 22:48:10');
INSERT INTO `exchange_log` VALUES (3, 15, '兑换额度:10', '2023-09-09 19:31:30');

-- ----------------------------
-- Table structure for limit
-- ----------------------------
DROP TABLE IF EXISTS `limit`;
CREATE TABLE `limit`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sid` int(11) NULL DEFAULT NULL,
  `uid` int(11) NULL DEFAULT NULL,
  `num` int(11) NULL DEFAULT NULL,
  `time` varchar(3) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `uid`(`uid`) USING BTREE,
  INDEX `limitofsign`(`sid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for log
-- ----------------------------
DROP TABLE IF EXISTS `log`;
CREATE TABLE `log`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `error` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `time` datetime(0) NULL DEFAULT NULL,
  `uid` int(11) NULL DEFAULT NULL,
  `sid` int(11) NULL DEFAULT NULL,
  `content` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 41 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for reg_temp
-- ----------------------------
DROP TABLE IF EXISTS `reg_temp`;
CREATE TABLE `reg_temp`  (
  `ip` varchar(16) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `code` varchar(4) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`ip`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for sent_log
-- ----------------------------
DROP TABLE IF EXISTS `sent_log`;
CREATE TABLE `sent_log`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '本地id',
  `uid` int(11) NOT NULL COMMENT '用户id',
  `sid` int(11) NOT NULL COMMENT '签名头id',
  `content` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '发信内容',
  `phone` varchar(11) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '发送的手机号',
  `time` datetime(0) NULL DEFAULT NULL COMMENT '发信时间',
  `task_id` int(11) NULL DEFAULT -1 COMMENT '任务id',
  `receive` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '回执内容',
  `receive_time` datetime(0) NULL DEFAULT NULL COMMENT '回执传递时间',
  `decrease_num` int(11) NULL DEFAULT NULL COMMENT '扣除额度数量',
  `status` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '0' COMMENT '发信状态码 0正在发送 1发送成功 2发送失败 3待回执',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `filter_content_time`(`content`, `time`, `status`) USING BTREE,
  INDEX `task`(`task_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 23 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for sign
-- ----------------------------
DROP TABLE IF EXISTS `sign`;
CREATE TABLE `sign`  (
  `sid` int(11) NOT NULL AUTO_INCREMENT COMMENT '签名头 id',
  `uid` int(11) NULL DEFAULT NULL COMMENT '用户 id',
  `content` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '签名',
  `key` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '密钥',
  `status` int(1) NULL DEFAULT 0 COMMENT '状态 0待审核 1审核通过 2拒绝',
  `super_id` int(11) NOT NULL COMMENT '上游id',
  `md5` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'md5',
  PRIMARY KEY (`sid`) USING BTREE,
  UNIQUE INDEX `unique`(`content`) USING BTREE,
  INDEX `uid`(`uid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sign
-- ----------------------------
INSERT INTO `sign` VALUES (9, 12, '【云耀科技】', 'afhaioejfoia', 1, 0, NULL);
INSERT INTO `sign` VALUES (12, 16, '【测试】', 'sdfsdhgiosjeofisef', 0, 231, '77b93a458dc6341470d9f3b381fb8392');

-- ----------------------------
-- Table structure for system
-- ----------------------------
DROP TABLE IF EXISTS `system`;
CREATE TABLE `system`  (
  `email_check` int(1) NULL DEFAULT 1 COMMENT '邮箱验证 0关闭 1开启',
  `smtp_server` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'smtp服务器需包括端口',
  `smtp_username` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'smtp登录用户名',
  `smtp_password` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'smtp登录密码',
  `smtp_nickname` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'smtp发信时的昵称',
  `smtp_ssl` int(1) NULL DEFAULT 0 COMMENT 'smtp是否是ssl 0不是 1是'
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of system
-- ----------------------------
INSERT INTO `system` VALUES (1, 'smtp.qq.com:465', 'test@qq.com', 'password', '发信名', 1);

-- ----------------------------
-- Table structure for template
-- ----------------------------
DROP TABLE IF EXISTS `template`;
CREATE TABLE `template`  (
  `tid` int(11) NOT NULL AUTO_INCREMENT COMMENT '模板id',
  `sid` int(11) NULL DEFAULT NULL COMMENT '签名id',
  `uid` int(11) NULL DEFAULT NULL COMMENT '用户id',
  `content` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '模板内容',
  `status` int(1) NULL DEFAULT 0 COMMENT '审核状态 0待审核 1审核通过 2审核拒绝',
  `reason` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '原因',
  `super_tid` int(11) NULL DEFAULT NULL COMMENT '上游模板id',
  PRIMARY KEY (`tid`) USING BTREE,
  INDEX `sid`(`sid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of template
-- ----------------------------
INSERT INTO `template` VALUES (3, 9, 12, '【云耀科技】你的验证码为：@，5分钟有效，请勿泄露他人', 1, '', 1639);

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `uid` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户uid 唯一',
  `username` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '用户名 唯一',
  `password` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '用户密码',
  `group` int(11) NOT NULL DEFAULT 0 COMMENT '用户权限组 默认0',
  `quota` int(11) NULL DEFAULT 0 COMMENT '短信剩余条数',
  `key` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `jwt` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'jwt 前后端分离必须',
  `status` int(1) NULL DEFAULT 0 COMMENT '0 正常 1封禁 2禁止发信（默认0）',
  `code` varchar(5) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '找回密码时使用的验证码',
  `phone` varchar(11) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '手机号',
  `limit_quota_phone` int(11) NULL DEFAULT 0 COMMENT '短信数少于多少时',
  `black_phone` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '黑名单手机号',
  `white_ip` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '白名单IP',
  PRIMARY KEY (`uid`) USING BTREE,
  UNIQUE INDEX `username`(`username`) USING BTREE,
  INDEX `jwt`(`jwt`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 17 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (12, '123456', 'e10adc3949ba59abbe56e057f20f883e', 6, 130, '0WV7Xl3IIMHcbQpTXobgtxJv4KfINg5s', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEyLCJ1c2VybmFtZSI6IjEyMzQ1NiIsInBhc3N3b3JkIjoiZTEwYWRjMzk0OWJhNTlhYmJlNTZlMDU3ZjIwZjg4M2UiLCJzdGF0dXMiOjAsImdyb3VwIjo2LCJiYWxhbmNlIjoxMzAsImV4cCI6MTY5NDk1NjEwMiwibmJmIjoxNjk0MzUxMzAyLCJpYXQiOjE2OTQzNTEzMDJ9.Nspswcl6RW5wCZzOdyPENgoZ-zYOh3ZObw3aNVFZUs0', 0, '', '', 0, '', '1.1.1.3');

SET FOREIGN_KEY_CHECKS = 1;
