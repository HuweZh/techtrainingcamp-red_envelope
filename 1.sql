/*
SQLyog Community v13.1.5  (64 bit)
MySQL - 8.0.17 : Database - red_envelope
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`red_envelope` /*!40100 DEFAULT CHARACTER SET utf8 */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `red_envelope`;

/*Table structure for table `envelope` */

DROP TABLE IF EXISTS `envelope`;

CREATE TABLE `envelope` (
  `envelope_id` bigint(20) NOT NULL COMMENT '红包id',
  `uid` bigint(20) NOT NULL COMMENT '拥有者',
  `value` int(11) NOT NULL COMMENT '红包金额（分）',
  `opened` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否打开',
  `snatch_time` bigint(20) NOT NULL COMMENT '抢到时间',
  `opened_time` bigint(20) DEFAULT NULL COMMENT '打开时间',
  PRIMARY KEY (`envelope_id`,`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE envelope ADD INDEX uid_index(`uid`);

/*Data for the table `envelope` */

/*Table structure for table `wallet` */

DROP TABLE IF EXISTS `wallet`;

CREATE TABLE `wallet` (
  `uid` bigint(20) NOT NULL,
  `money` int(11) DEFAULT '0',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Data for the table `wallet` */

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
