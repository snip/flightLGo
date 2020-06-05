SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `flightlog`;
CREATE TABLE `flightlog` (
	  `id` bigint NOT NULL AUTO_INCREMENT,
	  `lastIp` varchar(255) CHARACTER SET latin1 COLLATE latin1_general_ci NOT NULL,
	  `aircraftId` varchar(255) CHARACTER SET latin1 COLLATE latin1_general_ci NOT NULL,
	  `aircraftReg` varchar(255) CHARACTER SET latin1 COLLATE latin1_general_ci DEFAULT NULL,
	  `aircraftCN` varchar(3) CHARACTER SET latin1 COLLATE latin1_general_ci DEFAULT NULL,
	  `aircraftTracked` char(1) CHARACTER SET latin1 COLLATE latin1_general_ci NOT NULL DEFAULT '',
	  `aircraftIdentified` char(1) CHARACTER SET latin1 COLLATE latin1_general_ci NOT NULL DEFAULT '',
	  `takeoffTimestamp` timestamp NULL DEFAULT NULL,
	  `takeoffAirfield` varchar(255) CHARACTER SET latin1 COLLATE latin1_general_ci DEFAULT NULL,
	  `landingTimestamp` timestamp NULL DEFAULT NULL,
	  `landingAirfield` varchar(255) CHARACTER SET latin1 COLLATE latin1_general_ci DEFAULT NULL,
	  PRIMARY KEY (`id`),
	  UNIQUE KEY `aircraftIdTakeoff` (`aircraftId`,`takeoffTimestamp`),
	  UNIQUE KEY `aircraftIdLanding` (`aircraftId`,`landingTimestamp`),
	  KEY `aircraftReg` (`aircraftReg`),
	  KEY `takeoffAirfield` (`takeoffAirfield`),
	  KEY `landingAirfield` (`landingAirfield`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1 COLLATE=latin1_general_ci;

