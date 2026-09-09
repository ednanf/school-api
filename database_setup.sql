-- Create database
CREATE DATABASE IF NOT EXISTS school_db;

USE school_db;

-- Create `classes` table
CREATE TABLE IF NOT EXISTS classes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    grade TINYINT NOT NULL CHECK (grade BETWEEN 1 AND 9),
    letter CHAR(1) NOT NULL CHECK (letter IN ('A', 'B', 'C', 'D')),
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

-- Insert base subjects into the table
INSERT IGNORE INTO subjects (name) VALUES
('English Language Arts'),
('Literature'),
('Mathematics'),
('Algebra'),
('Geometry'),
('Calculus'),
('Statistics'),
('General Science'),
('Biology'),
('Chemistry'),
('Physics'),
('Environmental Science'),
('Social Studies'),
('World History'),
('US History'),
('Geography'),
('Civics & Government'),
('Economics'),
('Physical Education'),
('Health'),
('Spanish'),
('French'),
('Art'),
('Music'),
('Computer Science');

-- Insert base classes into the table
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

-- Create `teachers` table
CREATE TABLE IF NOT EXISTS teachers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
) AUTO_INCREMENT=100;

-- Create `teacher_assignments` table (The Junction)
CREATE TABLE IF NOT EXISTS teacher_assignments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    teacher_id INT NOT NULL,
    class_id INT NOT NULL,
    subject_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    -- Foreign Keys
    CONSTRAINT fk_assignment_teacher FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE,
    CONSTRAINT fk_assignment_class FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE,
    CONSTRAINT fk_assignment_subject FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,

    -- Prevents assigning two different teachers to teach the exact same subject to the exact same class
    -- (Remove this UNIQUE KEY if the school allows co-teaching)
    UNIQUE KEY uq_class_subject (class_id, subject_id)
);
