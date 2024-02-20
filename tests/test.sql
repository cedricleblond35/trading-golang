CREATE DATABASE IF NOT EXISTS trading;
use trading;
DROP TABLE IF EXISTS `candleUS100`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `candleUS100` (
  `id` int NOT NULL AUTO_INCREMENT,
  `ctm` int NOT NULL,
  `date`  varchar(255)  NOT NULL,
  `close`  decimal(8,2) NOT NULL  DEFAULT '0.00',
  `high` decimal(8,2) NOT NULL DEFAULT '0.00',
  `low` decimal(8,2) NOT NULL DEFAULT '0.00',
  `open` decimal(8,2) NOT NULL DEFAULT '0.00',
  `vol` int NOT NULL,
  `period`  SMALLINT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `candleuniq_id` (`ctm`,`period`),
  INDEX name (`ctm`,`period`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1 CHECKSUM=1;

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `email` varchar(255) NOT NULL,
  `userbrocker` varchar(255) NOT NULL,
  `passBrocker` varchar(255) NOT NULL,
  `lastConnection` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1 CHECKSUM=1;