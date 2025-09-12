-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema table_tennis
-- -----------------------------------------------------

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


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
