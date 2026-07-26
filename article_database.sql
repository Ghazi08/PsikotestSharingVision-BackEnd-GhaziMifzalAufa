-- Database dump for Sharing Vision 2023 - Test Backend (Post Article Use Case)
-- Database: article
-- Table: posts

CREATE DATABASE IF NOT EXISTS `article`;
USE `article`;

DROP TABLE IF EXISTS `posts`;

CREATE TABLE `posts` (
  `Id` INT NOT NULL AUTO_INCREMENT,
  `Title` VARCHAR(200) NOT NULL,
  `Content` TEXT NOT NULL,
  `Category` VARCHAR(100) NOT NULL,
  `Created_date` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `Updated_date` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `Status` VARCHAR(100) NOT NULL COMMENT 'Publish | Draft | Thrash',
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Sample Data
INSERT INTO `posts` (`Title`, `Content`, `Category`, `Status`) VALUES
('Belajar Backend Development', 'Konten artikel tentang dasar-dasar backend development dengan Golang atau Python.', 'Technology', 'Publish'),
('Tips Database MySQL', 'Tips dan trik optimasi query dan struktur tabel MySQL.', 'Database', 'Draft');
