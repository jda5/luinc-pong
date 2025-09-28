-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema table_tennis
-- -----------------------------------------------------
DROP SCHEMA IF EXISTS `table_tennis` ;

-- -----------------------------------------------------
-- Schema table_tennis
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `table_tennis` DEFAULT CHARACTER SET utf8 ;
USE `table_tennis` ;

-- -----------------------------------------------------
-- Table `table_tennis`.`players`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `table_tennis`.`players` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(63) NOT NULL,
  `elo_rating` DOUBLE NOT NULL DEFAULT 1000,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC) VISIBLE,
  INDEX `idx_elo_rating` (`elo_rating` ASC) COMMENT 'For quick sorting' VISIBLE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `table_tennis`.`games`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `table_tennis`.`games` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `winner_id` INT NOT NULL,
  `loser_id` INT NOT NULL,
  `winner_score` TINYINT UNSIGNED NULL,
  `loser_score` TINYINT UNSIGNED NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `fk_player_winner_id_idx` (`winner_id` ASC) VISIBLE,
  INDEX `fk_player_loser_id_idx` (`loser_id` ASC) VISIBLE,
  INDEX `idx_created_at` (`created_at` ASC) COMMENT 'Speeds up look-up on game history' VISIBLE,
  INDEX `idx_game_players` (`winner_id` ASC, `loser_id` ASC) COMMENT 'For fast lookups on games involving two players.' VISIBLE,
  CONSTRAINT `fk_player_winner_id`
    FOREIGN KEY (`winner_id`)
    REFERENCES `table_tennis`.`players` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `fk_player_loser_id`
    FOREIGN KEY (`loser_id`)
    REFERENCES `table_tennis`.`players` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `table_tennis`.`achievement`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `table_tennis`.`achievement` (
  `id` INT NOT NULL,
  `title` VARCHAR(63) NOT NULL,
  `description` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `table_tennis`.`player_achievement`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `table_tennis`.`player_achievement` (
  `player_id` INT NOT NULL,
  `achievement_id` INT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `fk_player_achievements_players1_idx` (`player_id` ASC) VISIBLE,
  INDEX `fk_player_achievements_achievement1_idx` (`achievement_id` ASC) VISIBLE,
  UNIQUE INDEX `UNIQUE_ player_achievement_id` (`player_id` ASC, `achievement_id` ASC) VISIBLE,
  INDEX `idx_achievement_created_at` (`created_at` ASC) COMMENT 'For fast sorting!' VISIBLE,
  CONSTRAINT `fk_player_achievements_players1`
    FOREIGN KEY (`player_id`)
    REFERENCES `table_tennis`.`players` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_player_achievements_achievement1`
    FOREIGN KEY (`achievement_id`)
    REFERENCES `table_tennis`.`achievement` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;

-- -----------------------------------------------------
-- Data for table `table_tennis`.`achievement`
-- -----------------------------------------------------
START TRANSACTION;
USE `table_tennis`;
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (1, 'Warming Up', 'Play your first game');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (2, 'Minimum Viable Pong', 'Play 10 games');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (3, 'Regular', 'Play 50 games');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (4, 'Centurion', 'Play 100 games');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (5, 'Legend', 'Play 250 games');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (6, 'Unicorn', 'Play 500 games');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (7, 'Chocolate', 'Win 11–0');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (8, 'Bottle Job', 'Win 11–1');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (9, 'Clutch', 'Win 12–10');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (10, 'Marathon Madness', 'Win a game that goes to 15+ points');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (11, 'Heartbreaker', 'Lose 10–12');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (12, 'Streaky', 'Win 5 games in a row');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (13, 'Unstoppable', 'Win 10 games in a row');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (14, 'Probably Cheating', 'Win 15 games in a row');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (15, 'I Get Knocked Down', 'Lose 5 games in a row');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (16, 'Hat Trick', 'Beat the same opponent 3 times in a row in a single day');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (17, 'Brutal', 'Beat the same opponent 5 times in a row in a single day');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (18, 'Nemesis', 'Lose to the same opponent 15 times');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (19, 'Rivalry', 'Play the same opponent 25 times');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (20, 'Social Butterfly', 'Play 5 different people in the office');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (21, 'Daily Standup', 'Play 5 games in a single day');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (22, 'Do You Even Work Here?', 'Play 10 games in a single day');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (23, 'Go Home', 'Play before 9am or after 5pm');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (24, 'Dedicated', 'Play on 3 consecutive days');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (25, 'Addicted', 'Play on 5 consecutive days');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (26, 'Hostile Takeover', 'Beat someone 100+ ELO points above you');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (27, 'Rising Star', 'Reach an ELO of 1100');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (28, 'Big Shot', 'Reach an ELO of 1200');
INSERT INTO `table_tennis`.`achievement` (`id`, `title`, `description`) VALUES (29, 'Final Boss', 'Reach an ELO of 1300');

COMMIT;

