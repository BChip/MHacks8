CREATE TABLE `travelListings` (
  `id` mediumint(8) unsigned NOT NULL auto_increment,
  `firstN` varchar(255) default NULL,
  `lastN` varchar(255) default NULL,
  `age` mediumint default NULL,
  `gender` varchar(1) default NULL,
  `city` varchar(255) default NULL,
  `startDate` varchar(255),
  `endDate` varchar(255),
  `interests` varchar(255) default NULL,
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=1;

