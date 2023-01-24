DROP DATABASE IF EXISTS `heroku_1d78b3a1b883e9c`;
CREATE DATABASE `heroku_1d78b3a1b883e9c`; 
USE `heroku_1d78b3a1b883e9c`;

CREATE TABLE `users` (
  `user_id` int(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL UNIQUE,
  `email` varchar(50) NOT NULL UNIQUE,
  `password` varchar(70) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

CREATE TABLE `card_ownerships` (
    `card_id` int(11) PRIMARY KEY NOT NULL AUTO_INCREMENT, 
    `user_id` int(11) NOT NULL,
    `version_id` varchar(50) NOT NULL, /*ID to identyfy a card version*/
    `oracle_id` varchar(50) NOT NULL, /* ID to identify a card. All versions of a single card have the same oracle_id*/
    `count`	int(11) NOT NULL,
    `extras`	varchar(50),
    `condi`	varchar(50),
    KEY `FK_user_id` (`user_id`),
	CONSTRAINT `FK_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION 
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;

CREATE TABLE `trades` (
    `trade_id` int(11) PRIMARY KEY NOT NULL AUTO_INCREMENT, 
    `user_id_origin` int(11) NOT NULL,
    `user_id_owner` int(11) NOT NULL,
    `card_id` int(11) NOT NULL,
    #`version_id` varchar(50) NOT NULL,
    `card_select` int(11) NOT NULL,
    #`extras`	varchar(50),
    #`condi`	    varchar(50),
	`status`	TINYINT SIGNED,
    KEY `FK_card_id` (`card_id`),
	CONSTRAINT `FK_card_id` FOREIGN KEY (`card_id`) REFERENCES `card_ownerships` (`card_id`) ON DELETE NO ACTION ON UPDATE NO ACTION 
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;