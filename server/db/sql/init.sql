-- Database setup for the smart-home system
CREATE DATABASE smarthome CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE smarthome;

SET GLOBAL wait_timeout=3;
SET GLOBAL interactive_timeout=3;

CREATE TABLE IF NOT EXISTS event_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    lamp VARCHAR(255),
    date DATETIME,
    state BOOLEAN
);
