CREATE TABLE IF NOT EXISTS cat_courses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS academyc_history (
    id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL,
    key VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL, 
    coins INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cat_semesters (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS semester_course (
    id SERIAL PRIMARY KEY,
    alumn_id INTEGER NOT NULL,
    semester_id INTEGER NOT NULL,
    subject_id INTEGER NOT NULL,
    final_grade DOUBLE PRECISION, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS partial_grades (
    id SERIAL PRIMARY KEY,
    semester_course_id INTEGER NOT NULL,
    partial_number INTEGER NOT NULL,
    grade DOUBLE PRECISION NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (semester_course_id, partial_number)
);

CREATE TABLE IF NOT EXISTS alumn (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    lastname1 VARCHAR(255) NOT NULL, 
    lastname2 VARCHAR(255),
    course_id INTEGER NOT NULL,
    current_semester INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



CREATE TABLE IF NOT EXISTS teacher (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    lastname1 VARCHAR(255) NOT NULL,
    lastname2 VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE academyc_history
ADD CONSTRAINT fk_academyc_history_course_id
FOREIGN KEY (course_id) REFERENCES cat_courses(id) ON DELETE CASCADE;

ALTER TABLE semester_course
ADD CONSTRAINT fk_semester_course_semester_id
FOREIGN KEY (semester_id) REFERENCES cat_semesters(id) ON DELETE CASCADE;

ALTER TABLE semester_course
ADD CONSTRAINT fk_semester_course_subject_id
FOREIGN KEY (subject_id) REFERENCES academyc_history(id) ON DELETE CASCADE;

ALTER TABLE semester_course
ADD CONSTRAINT fk_semester_course_alumn_id
FOREIGN KEY (alumn_id) REFERENCES alumn(id) ON DELETE CASCADE;

ALTER TABLE partial_grades
ADD CONSTRAINT fk_partial_grades_semester_course
FOREIGN KEY (semester_course_id) REFERENCES semester_course(id) ON DELETE CASCADE;

ALTER TABLE alumn 
ADD CONSTRAINT fk_current_semester_alumn_id
FOREIGN KEY (current_semester) REFERENCES cat_semesters(id);