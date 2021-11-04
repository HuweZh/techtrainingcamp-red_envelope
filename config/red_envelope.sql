/*
 Navicat Premium Data Transfer

 Source Server         : huhusw
 Source Server Type    : MySQL
 Source Server Version : 80025
 Source Host           : localhost:3306
 Source Schema         : red_envelope

 Target Server Type    : MySQL
 Target Server Version : 80025
 File Encoding         : 65001

 Date: 04/11/2021 15:35:55
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for envelope
-- ----------------------------
DROP TABLE IF EXISTS `envelope`;
CREATE TABLE `envelope`  (
  `envelope_id` bigint NOT NULL COMMENT '红包id',
  `user_id` int NOT NULL COMMENT '用户id',
  `value` int NULL DEFAULT NULL COMMENT '红包金额',
  `opened` int NULL DEFAULT NULL COMMENT '红包的状态，0：未打开，1：已打开',
  `snatch_time` int NULL DEFAULT NULL COMMENT '红包获取时间，UNIX时间戳',
  PRIMARY KEY (`envelope_id`, `user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_cs_0900_ai_ci COMMENT = '红包表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `user_id` int NOT NULL COMMENT '用户id',
  `max_count` int NULL DEFAULT NULL COMMENT '最多抢max_count次',
  `cur_count` int NULL DEFAULT NULL COMMENT '当前第几次抢',
  `create_time` int NULL DEFAULT NULL COMMENT '创建时间，unix时间戳',
  `amount` int NULL DEFAULT 0 COMMENT '用户当前钱包的金额，不算未拆红包',
  PRIMARY KEY (`user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_cs_0900_ai_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
