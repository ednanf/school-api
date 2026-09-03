-- Create database
CREATE DATABASE IF NOT EXISTS school_db;

USE school_db;

-- Create `classes` table
CREATE TABLE IF NOT EXISTS classes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    grade TINYINT NOT NULL CHECK (grade BETWEEN 1 AND 8),
    letter CHAR(1) NOT NULL CHECK (letter IN ('A', 'B')),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE KEY uq_class_grade_letter (grade, letter)
) AUTO_INCREMENT=1;

-- Create `subjects` table
CREATE TABLE IF NOT EXISTS subjects (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
) AUTO_INCREMENT=1;

-- Insert base classes to the table
INSERT IGNORE INTO classes (grade, letter) VALUES
(1, 'A'), (1, 'B'), (2, 'A'), (2, 'B'),
(3, 'A'), (3, 'B'), (4, 'A'), (4, 'B'),
(5, 'A'), (5, 'B'), (6, 'A'), (6, 'B'),
(7, 'A'), (7, 'B'), (8, 'A'), (8, 'B');

-- Create `students` table
CREATE TABLE IF NOT EXISTS students (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    class_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    CONSTRAINT fk_students_class FOREIGN KEY (class_id) REFERENCES classes(id)
) AUTO_INCREMENT=100;
