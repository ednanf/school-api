-- =============================================================================
-- SCHOOL DATABASE SCHEMA & INITIAL SEED DATA
-- Compatible with MariaDB / MySQL
-- =============================================================================

CREATE DATABASE IF NOT EXISTS school_db;
USE school_db;

-- =============================================================================
-- 1. ACADEMIC STRUCTURE (SCHEMA & SEED)
-- =============================================================================

-- `classes`: Defines individual student cohorts (e.g., Grade 1-A, Grade 5-B).
CREATE TABLE IF NOT EXISTS classes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    grade TINYINT NOT NULL CHECK (grade BETWEEN 1 AND 9),
    letter CHAR(1) NOT NULL CHECK (letter IN ('A', 'B', 'C', 'D')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    UNIQUE KEY uq_class_grade_letter (grade, letter)
) AUTO_INCREMENT=1;

-- `subjects`: Catalog of academic courses offered across grades.
CREATE TABLE IF NOT EXISTS subjects (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
) AUTO_INCREMENT=1;

-- Seed default curriculum subjects
INSERT IGNORE INTO subjects (name, created_at, updated_at) VALUES
('English Language Arts', NOW(), NOW()),
('Literature', NOW(), NOW()),
('Mathematics', NOW(), NOW()),
('Algebra', NOW(), NOW()),
('Geometry', NOW(), NOW()),
('Calculus', NOW(), NOW()),
('Statistics', NOW(), NOW()),
('General Science', NOW(), NOW()),
('Biology', NOW(), NOW()),
('Chemistry', NOW(), NOW()),
('Physics', NOW(), NOW()),
('Environmental Science', NOW(), NOW()),
('Social Studies', NOW(), NOW()),
('World History', NOW(), NOW()),
('US History', NOW(), NOW()),
('Geography', NOW(), NOW()),
('Civics & Government', NOW(), NOW()),
('Economics', NOW(), NOW()),
('Physical Education', NOW(), NOW()),
('Health', NOW(), NOW()),
('Spanish', NOW(), NOW()),
('French', NOW(), NOW()),
('Art', NOW(), NOW()),
('Music', NOW(), NOW()),
('Computer Science', NOW(), NOW());

-- Seed default class cohorts (Grades 1-8, sections A & B)
INSERT IGNORE INTO classes (grade, letter, created_at, updated_at) VALUES
(1, 'A', NOW(), NOW()), (1, 'B', NOW(), NOW()),
(2, 'A', NOW(), NOW()), (2, 'B', NOW(), NOW()),
(3, 'A', NOW(), NOW()), (3, 'B', NOW(), NOW()),
(4, 'A', NOW(), NOW()), (4, 'B', NOW(), NOW()),
(5, 'A', NOW(), NOW()), (5, 'B', NOW(), NOW()),
(6, 'A', NOW(), NOW()), (6, 'B', NOW(), NOW()),
(7, 'A', NOW(), NOW()), (7, 'B', NOW(), NOW()),
(8, 'A', NOW(), NOW()), (8, 'B', NOW(), NOW());


-- =============================================================================
-- 2. PEOPLE & INSTRUCTIONAL ASSIGNMENTS (SCHEMA & SEED)
-- =============================================================================

-- `students`: Enrolled students mapped to their primary class cohort.
CREATE TABLE IF NOT EXISTS students (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    class_id INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    CONSTRAINT fk_students_class FOREIGN KEY (class_id) REFERENCES classes(id)
) AUTO_INCREMENT=100;

CREATE INDEX idx_students_class_active ON students (class_id, is_active);

-- `teachers`: Instructional staff dedicated to subject teaching.
CREATE TABLE IF NOT EXISTS teachers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
) AUTO_INCREMENT=100;

-- Seed Teachers
INSERT IGNORE INTO teachers (first_name, last_name, email, is_active, created_at, updated_at) VALUES
('John', 'Doe', 'joh_doe@icloud.com', TRUE, NOW(), NOW()),
('Sarah', 'Jenkins', 's.jenkins@school.edu', TRUE, NOW(), NOW()),
('Robert', 'Chen', 'r.chen@school.edu', TRUE, NOW(), NOW()),
('Maria', 'Garcia', 'm.garcia@school.edu', TRUE, NOW(), NOW()),
('David', 'Kowalski', 'd.kowalski@school.edu', TRUE, NOW(), NOW());

-- `teacher_assignments`: Junction table linking Teachers to specific Subjects and Classes.
CREATE TABLE IF NOT EXISTS teacher_assignments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    teacher_id INT NOT NULL,
    class_id INT NOT NULL,
    subject_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    CONSTRAINT fk_assignment_teacher FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE,
    CONSTRAINT fk_assignment_class FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE,
    CONSTRAINT fk_assignment_subject FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,

    UNIQUE KEY uq_class_subject (class_id, subject_id)
);


-- =============================================================================
-- 3. ORGANIZATIONAL STRUCTURE (SCHEMA & SEED)
-- =============================================================================

-- `departments`: Top-level administrative units.
CREATE TABLE IF NOT EXISTS departments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Seed Departments
INSERT IGNORE INTO departments (name, description, is_active, created_at, updated_at) VALUES
('Executive Direction', 'Overall institutional leadership, policy enforcement, and strategic decision-making.', TRUE, NOW(), NOW()),
('Academic Affairs', 'Curriculum design, teaching quality control, pedagogical coordination, and academic scheduling.', TRUE, NOW(), NOW()),
('Student Services & Counseling', 'Student welfare, mental health support, discipline, attendance tracking, and career guidance.', TRUE, NOW(), NOW()),
('Administration & Secretarial', 'Student records, enrollment, official transcripts, parent communications, and office support.', TRUE, NOW(), NOW()),
('Finance & Human Resources', 'Tuition billing, payroll, staff hiring, budgeting, and general accounting.', TRUE, NOW(), NOW()),
('IT & Technical Services', 'Network maintenance, school management software, hardware support, and database administration.', TRUE, NOW(), NOW()),
('Facilities & Operations', 'Building maintenance, campus safety, janitorial services, and physical infrastructure.', TRUE, NOW(), NOW());

-- `staff_positions`: Standardized job titles tied directly to a department.
CREATE TABLE IF NOT EXISTS staff_positions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    department_id INT NOT NULL,
    title VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    FOREIGN KEY (department_id) REFERENCES departments(id)
);

CREATE INDEX idx_staff_positions_dept_id ON staff_positions(department_id);

-- Seed Staff Positions
INSERT IGNORE INTO staff_positions (department_id, title, description, is_active, created_at, updated_at) VALUES
((SELECT id FROM departments WHERE name = 'Executive Direction'), 'School Principal', 'Main administrator responsible for overall school operations, staff oversight, and regulatory compliance.', TRUE, NOW(), NOW()),
((SELECT id FROM departments WHERE name = 'Executive Direction'), 'Vice Principal', 'Assists the principal in daily operations, handles high-level discipline issues, and acts as acting director when needed.', TRUE, NOW(), NOW()),
((SELECT id FROM departments WHERE name = 'Academic Affairs'), 'Academic Coordinator', 'Supervises teachers, oversees syllabus execution, and manages exam schedules across grade levels.', TRUE, NOW(), NOW()),
((SELECT id FROM departments WHERE name = 'Academic Affairs'), 'Pedagogical Advisor', 'Supports teachers with instructional techniques, learning methodologies, and special needs integration.', TRUE, NOW(), NOW()),
((SELECT id FROM departments WHERE name = 'Student Services & Counseling'), 'School Guidance Counselor', 'Provides academic guidance, emotional support, and conflict resolution services for students.', TRUE, NOW(), NOW()),
((SELECT id FROM departments WHERE name = 'Student Services & Counseling'), 'Dean of Students', 'Manages student behavior policies, attendance enforcement, and extracurricular activities.', TRUE, NOW(), NOW()),
((SELECT id FROM departments WHERE name = 'Administration & Secretarial'), 'School Registrar', 'Maintains official student records, grade archives, enrollment files, and official transcript issuance.', TRUE, NOW(), NOW()),
((SELECT id FROM departments WHERE name = 'Administration & Secretarial'), 'Administrative Secretary', 'Handles reception, incoming inquiries, document filing, and direct communication with parents.', TRUE, NOW(), NOW()),
((SELECT id FROM departments WHERE name = 'Finance & Human Resources'), 'Bursar / Financial Administrator', 'Manages student tuition collections, vendor invoicing, expense tracking, and petty cash.', TRUE, NOW(), NOW()),
((SELECT id FROM departments WHERE name = 'Finance & Human Resources'), 'HR Generalist', 'Oversees employee hiring, staff attendance tracking, benefits management, and payroll processing.', TRUE, NOW(), NOW()),
((SELECT id FROM departments WHERE name = 'IT & Technical Services'), 'Database & Systems Administrator', 'Maintains the school database system, manages user authentication accounts, and handles system backups.', TRUE, NOW(), NOW()),
((SELECT id FROM departments WHERE name = 'IT & Technical Services'), 'IT Support Technician', 'Provides technical assistance for classroom computers, projector setups, and local network issues.', TRUE, NOW(), NOW()),
((SELECT id FROM departments WHERE name = 'Facilities & Operations'), 'Facilities Manager', 'Oversees campus physical safety, routine maintenance, repairs, and vendor contracts.', TRUE, NOW(), NOW());

-- `staff`: Non-instructional / operational employees occupying a specific position.
CREATE TABLE IF NOT EXISTS staff (
    id INT AUTO_INCREMENT PRIMARY KEY,
    position_id INT NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    hire_date DATE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    FOREIGN KEY (position_id) REFERENCES staff_positions(id)
) AUTO_INCREMENT=100;

CREATE INDEX idx_staff_position_id ON staff(position_id);

-- Seed Staff
INSERT IGNORE INTO staff (position_id, first_name, last_name, email, hire_date, is_active, created_at, updated_at) VALUES
((SELECT id FROM staff_positions WHERE title = 'School Principal'), 'Margaret', 'Thatcher', 'm.thatcher@school.edu', '2018-08-01', TRUE, NOW(), NOW()),
((SELECT id FROM staff_positions WHERE title = 'Academic Coordinator'), 'Jane', 'Smith', 'j.smith@school.edu', '2020-01-15', TRUE, NOW(), NOW()),
((SELECT id FROM staff_positions WHERE title = 'School Guidance Counselor'), 'Emily', 'Watson', 'e.watson@school.edu', '2021-09-01', TRUE, NOW(), NOW()),
((SELECT id FROM staff_positions WHERE title = 'School Registrar'), 'Carlos', 'Mendoza', 'c.mendoza@school.edu', '2019-03-10', TRUE, NOW(), NOW()),
((SELECT id FROM staff_positions WHERE title = 'Database & Systems Administrator'), 'Alex', 'Turner', 'a.turner@school.edu', '2022-05-20', TRUE, NOW(), NOW()),
((SELECT id FROM staff_positions WHERE title = 'Bursar / Financial Administrator'), 'Patricia', 'Arquette', 'p.arquette@school.edu', '2017-11-12', TRUE, NOW(), NOW());


-- =============================================================================
-- 4. AUTHENTICATION & ACCESS CONTROL (SCHEMA ONLY)
-- =============================================================================

-- `users`: System login accounts linked optionally to a physical staff member.
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    staff_id INT UNIQUE NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    password_reset_token VARCHAR(255) NULL,
    password_token_expires DATETIME NULL,
    role VARCHAR(30) NOT NULL DEFAULT 'STAFF',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMP NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,

    FOREIGN KEY (staff_id) REFERENCES staff(id) ON DELETE SET NULL
) AUTO_INCREMENT=100;

CREATE INDEX idx_users_reset_token ON users(password_reset_token);
CREATE INDEX idx_users_is_active ON users(is_active);
