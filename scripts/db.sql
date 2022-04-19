SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for logs
-- ----------------------------
DROP TABLE IF EXISTS `logs`;
CREATE TABLE `logs`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `job_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `command` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `output` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `start_time` bigint(20) UNSIGNED NOT NULL DEFAULT 0,
  `end_time` bigint(20) UNSIGNED NOT NULL DEFAULT 0,
  `plan_time` bigint(20) UNSIGNED NOT NULL DEFAULT 0,
  `schedule_time` bigint(20) UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
