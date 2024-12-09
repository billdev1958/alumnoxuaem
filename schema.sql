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

CREATE TABLE IF NOT EXISTS semester_grades (
    id SERIAL PRIMARY KEY,
    alumn_id INTEGER NOT NULL,
    semester_id INTEGER NOT NULL,
    final_semester_grade DOUBLE PRECISION, -- Promedio de las calificaciones finales de materias
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (alumn_id, semester_id),
    FOREIGN KEY (alumn_id) REFERENCES alumn(id) ON DELETE CASCADE,
    FOREIGN KEY (semester_id) REFERENCES cat_semesters(id) ON DELETE CASCADE
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

CREATE OR REPLACE FUNCTION update_final_grade()
RETURNS TRIGGER AS $$
BEGIN
    IF (SELECT COUNT(*) FROM partial_grades
        WHERE semester_course_id = NEW.semester_course_id) = 2 THEN
        
        UPDATE semester_course
        SET final_grade = (
            SELECT AVG(grade)
            FROM partial_grades
            WHERE semester_course_id = NEW.semester_course_id
        )
        WHERE id = NEW.semester_course_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER calculate_final_grade
AFTER INSERT OR UPDATE ON partial_grades
FOR EACH ROW
EXECUTE FUNCTION update_final_grade();

CREATE OR REPLACE FUNCTION update_final_semester_grade()
RETURNS TRIGGER AS $$
BEGIN
    -- Verifica si todas las materias del semestre tienen una `final_grade`
    IF (SELECT COUNT(*) 
        FROM semester_course
        WHERE semester_id = NEW.semester_id 
          AND alumn_id = NEW.alumn_id 
          AND final_grade IS NULL) = 0 THEN
        
        -- Calcula el promedio de `final_grade` de todas las materias del semestre
        UPDATE semester_grades
        SET final_semester_grade = (
            SELECT AVG(final_grade)
            FROM semester_course
            WHERE semester_id = NEW.semester_id 
              AND alumn_id = NEW.alumn_id
        ),
        updated_at = CURRENT_TIMESTAMP
        WHERE semester_id = NEW.semester_id 
          AND alumn_id = NEW.alumn_id;

        -- Si no existe un registro en `semester_grades`, lo inserta
        IF NOT FOUND THEN
            INSERT INTO semester_grades (alumn_id, semester_id, final_semester_grade, created_at, updated_at)
            VALUES (
                NEW.alumn_id,
                NEW.semester_id,
                (SELECT AVG(final_grade) 
                 FROM semester_course 
                 WHERE semester_id = NEW.semester_id 
                   AND alumn_id = NEW.alumn_id),
                CURRENT_TIMESTAMP,
                CURRENT_TIMESTAMP
            );
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER calculate_final_semester_grade
AFTER UPDATE OF final_grade ON semester_course
FOR EACH ROW
EXECUTE FUNCTION update_final_semester_grade();
