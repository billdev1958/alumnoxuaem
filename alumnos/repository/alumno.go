package repository

import (
	"alumnos/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxStorage struct {
	DbPool *pgxpool.Pool
}

func NewPgxStorage(dbPool *pgxpool.Pool) *PgxStorage {
	return &PgxStorage{DbPool: dbPool}
}

func (s *PgxStorage) RegisterAlumn(ctx context.Context, request models.RegisterAlumnRequest) (int, error) {
	tx, err := s.DbPool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insertar al alumno en la tabla `alumn`
	var alumnoID int
	insertAlumnQuery := `
		INSERT INTO alumn (name, lastname1, lastname2, course_id, current_semester)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	err = tx.QueryRow(ctx, insertAlumnQuery,
		request.Name, request.Lastname1, request.Lastname2,
		request.CourseID, request.CurrentCourseID).Scan(&alumnoID)
	if err != nil {
		return 0, fmt.Errorf("error al registrar alumno: %w", err)
	}

	// Asignar materias del curso al alumno en el semestre actual
	insertSubjectQuery := `
		INSERT INTO semester_course (alumn_id, semester_id, subject_id)
		VALUES ($1, $2, $3);
	`

	for _, subject := range request.Subjects {
		_, err = tx.Exec(ctx, insertSubjectQuery, alumnoID, request.CurrentCourseID, subject.ID)
		if err != nil {
			return 0, fmt.Errorf("error al asignar materia (ID: %d): %w", subject.ID, err)
		}
	}

	// Confirmar la transacción
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("error al confirmar transacción: %w", err)
	}

	return alumnoID, nil
}

func (s *PgxStorage) RegistrarEnSemestreConMaterias(ctx context.Context, alumnoID, semesterID int, subjectIDs []int) error {
	tx, err := s.DbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO semester_course (alumn_id, semester_id, subject_id)
		VALUES ($1, $2, $3);
	`

	for _, subjectID := range subjectIDs {
		_, err = tx.Exec(ctx, query, alumnoID, semesterID, subjectID)
		if err != nil {
			return fmt.Errorf("error al registrar materia %d: %w", subjectID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error al confirmar transacción: %w", err)
	}

	return nil
}

func (s *PgxStorage) RegistrarCalificacionParcial(ctx context.Context, semesterCourseID, partialNumber int, grade float64) error {
	query := `
		INSERT INTO partial_grades (semester_course_id, partial_number, grade)
		VALUES ($1, $2, $3)
		ON CONFLICT (semester_course_id, partial_number)
		DO UPDATE SET grade = $3;
	`

	_, err := s.DbPool.Exec(ctx, query, semesterCourseID, partialNumber, grade)
	if err != nil {
		return fmt.Errorf("error al registrar o actualizar calificación del parcial %d: %w", partialNumber, err)
	}

	return nil
}

func (s *PgxStorage) GenerarCalificacionesAgrupadasPorSemestre(ctx context.Context, alumnoID int) ([]models.SemestreCalificaciones, float64, error) {
	query := `
		SELECT sc.semester_id, sc.subject_id, pg.partial_number, pg.grade
		FROM semester_course sc
		JOIN partial_grades pg ON sc.id = pg.semester_course_id
		WHERE sc.alumn_id = $1
		ORDER BY sc.semester_id, sc.subject_id, pg.partial_number;
	`

	rows, err := s.DbPool.Query(ctx, query, alumnoID)
	if err != nil {
		return nil, 0, fmt.Errorf("error al obtener calificaciones: %w", err)
	}
	defer rows.Close()

	var semestres []models.SemestreCalificaciones
	calificacionTotal := 0.0
	numCalificaciones := 0

	// Mapa para agrupar calificaciones por semestre
	calificacionesPorSemestre := make(map[int]*models.SemestreCalificaciones)

	for rows.Next() {
		var semesterID, subjectID, partialNumber int
		var grade float64

		if err := rows.Scan(&semesterID, &subjectID, &partialNumber, &grade); err != nil {
			return nil, 0, fmt.Errorf("error al procesar filas: %w", err)
		}

		// Verificar si el semestre ya fue agregado
		if _, exists := calificacionesPorSemestre[semesterID]; !exists {
			calificacionesPorSemestre[semesterID] = &models.SemestreCalificaciones{
				SemesterID: semesterID,
				Materias:   []models.MateriaCalificaciones{},
			}
		}

		// Verificar si la materia ya fue agregada al semestre
		materias := calificacionesPorSemestre[semesterID].Materias
		var materia *models.MateriaCalificaciones
		for i := range materias {
			if materias[i].SubjectID == subjectID {
				materia = &materias[i]
				break
			}
		}
		if materia == nil {
			// Agregar nueva materia si no existe
			materia = &models.MateriaCalificaciones{
				SubjectID: subjectID,
				Parciales: []models.CalificacionParcial{},
				Promedio:  0,
			}
			calificacionesPorSemestre[semesterID].Materias = append(calificacionesPorSemestre[semesterID].Materias, *materia)
		}

		// Agregar el parcial a la materia
		materia.Parciales = append(materia.Parciales, models.CalificacionParcial{
			PartialNumber: partialNumber,
			Grade:         grade,
		})

		// Actualizar totales
		calificacionTotal += grade
		numCalificaciones++
	}

	// Calcular promedios por materia y semestre
	for _, semestre := range calificacionesPorSemestre {
		for i, materia := range semestre.Materias {
			totalMateria := 0.0
			for _, parcial := range materia.Parciales {
				totalMateria += parcial.Grade
			}
			materia.Promedio = totalMateria / float64(len(materia.Parciales))
			semestre.Materias[i] = materia
		}
		semestres = append(semestres, *semestre)
	}

	// Calcular promedio general
	promedioFinal := calificacionTotal / float64(numCalificaciones)

	return semestres, promedioFinal, nil
}

func (s *PgxStorage) SeedCatCourses(ctx context.Context) error {
	query := `
		INSERT INTO cat_courses (name) 
		VALUES ($1)
		ON CONFLICT DO NOTHING;
	`

	_, err := s.DbPool.Exec(ctx, query, "INGENIERIA EN COMPUTACION")
	if err != nil {
		return fmt.Errorf("error al hacer seed de cat_courses: %w", err)
	}

	return nil
}

func (s *PgxStorage) SeedAcademycHistory(ctx context.Context) error {
	query := `
		INSERT INTO academyc_history (course_id, key, name, coins)
		VALUES
		(1, 'LINC01', 'ALGEBRA LINEAL', 7),
		(1, 'LINC02', 'ALGEBRA SUPERIOR', 7),
		(1, 'LINC03', 'CALCULO I', 7),
		(1, 'LINC04', 'CALCULO II', 7),
		(1, 'LINC05', 'CALCULO III', 7),
		(1, 'LINC06', 'COMUNICACION Y RELACIONES HUMANAS', 7),
		(1, 'LINC07', 'ECUACIONES DIFERENCIALES', 7),
		(1, 'LINC08', 'EL INGENIERO Y SU ENTORNO SOCIOECONOMICO', 7),
		(1, 'LINC09', 'ELECTROMAGNETISMO', 7),
		(1, 'LINC10', 'EPISTEMOLOGIA', 7),
		(1, 'LINC11', 'FISICA', 7),
		(1, 'LINC12', 'GEOMETRIA ANALITICA', 7),
		(1, 'LINC13', 'MATEMATICAS DISCRETAS', 7),
		(1, 'LINC14', 'PROBABILIDAD Y ESTADISTICA', 7),
		(1, 'LINC15', 'PROGRAMACION I', 7),
		(1, 'LINC29', 'QUIMICA', 7),
		(1, 'LMU209', 'INGLES 5', 6),
		(1, 'LMU306', 'INGLES 6', 6),
		(1, 'LMU404', 'INGLES 7', 6),
		(1, 'LMU505', 'INGLES 8', 6),
		(1, 'LINC16', 'ADMINISTRACION DE PROYECTOS INFORMATICOS', 7),
		(1, 'LINC17', 'ADMINISTRACION DE RECURSOS INFORMATICOS', 7),
		(1, 'LINC18', 'ARQUITECTURA DE COMPUTADORAS', 7),
		(1, 'LINC19', 'ARQUITECTURA DE REDES', 5),
		(1, 'LINC20', 'BASES DE DATOS I', 7),
		(1, 'LINC21', 'BASES DE DATOS II', 5),
		(1, 'LINC22', 'CIRCUITOS ELECTRICOS Y ELECTRONICOS', 10),
		(1, 'LINC23', 'COMPILADORES', 7),
		(1, 'LINC24', 'ENSAMBLADORES', 7),
		(1, 'LINC25', 'GRAFICACION COMPUTACIONAL', 5),
		(1, 'LINC26', 'INGENIERIA DE SOFTWARE I', 7),
		(1, 'LINC27', 'INGENIERIA DE SOFTWARE II', 7),
		(1, 'LINC28', 'INTELIGENCIA ARTIFICIAL', 7),
		(1, 'LINC30', 'METODOS ESTADISTICOS', 7),
		(1, 'LINC31', 'METODOS NUMERICOS', 5),
		(1, 'LINC32', 'PARADIGMAS DE PROGRAMACION I', 5),
		(1, 'LINC33', 'PARADIGMAS DE PROGRAMACION II', 5),
		(1, 'LINC34', 'PROCESAMIENTO DE IMAGENES DIGITALES', 7),
		(1, 'LINC35', 'PROGRAMACION II', 7),
		(1, 'LINC36', 'PROTOCOLOS DE COMUNICACION DE DATOS', 7),
		(1, 'LINC37', 'ROBOTICA', 7),
		(1, 'LINC38', 'SEGURIDAD DE LA INFORMACION', 7),
		(1, 'LINC39', 'SISTEMAS ANALOGICOS', 7),
		(1, 'LINC40', 'SISTEMAS DIGITALES', 7),
		(1, 'LINC41', 'SISTEMAS OPERATIVOS', 7),
		(1, 'LINC42', 'TRANSMISION DE DATOS', 7),
		(1, 'L41004', 'INVESTIGACION DE OPERACIONES', 7),
		(1, 'LINC43', 'CIENCIA DE LOS DATOS', 5),
		(1, 'LINC44', 'ETICA PROFESIONAL Y SUSTENTABILIDAD', 6),
		(1, 'LINC45', 'GESTION DE PROYECTOS DE INVESTIGACION', 4),
		(1, 'LINC46', 'PROYECTO INTEGRAL DE COMUNICACION DE DATOS', 5),
		(1, 'LINC47', 'PROYECTO INTEGRAL DE INGENIERIA DE SOFTWARE', 5),
		(1, 'LINC48', 'SISTEMAS EMBEBIDOS', 6),
		(1, 'LINC49', 'TECNOLOGIAS COMPUTACIONALES I', 5),
		(1, 'LINC50', 'TECNOLOGIAS COMPUTACIONALES II', 5),
		(1, 'LINC51', 'INTEGRATIVA PROFESIONAL', 8),
		(1, 'LINC52', 'PRACTICA PROFESIONAL', 30),
		(1, 'LINC53', 'ANALISIS Y DISEÑO DE REDES', 5),
		(1, 'LINC54', 'COMPUTING IN INDUSTRY', 5),
		(1, 'LINC55', 'GESTION DE REDES', 5),
		(1, 'LINC56', 'INTERACCION HOMBRE-MAQUINA', 5),
		(1, 'LINC57', 'RECONOCIMIENTO DE PATRONES', 5),
		(1, 'LINC58', 'SISTEMAS INTERACTIVOS', 5),
		(1, 'LINC59', 'TECNOLOGIAS EMERGENTES', 5),
		(1, 'LINC60', 'TOPICOS DE TECNOLOGIAS DE DATOS', 5),
		(1, 'LINC61', 'VISION ARTIFICIAL', 5)
		ON CONFLICT DO NOTHING;
	`

	_, err := s.DbPool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error al insertar materias en academyc_history: %w", err)
	}

	return nil
}

func (s *PgxStorage) GetCourses(ctx context.Context) ([]models.Course, error) {
	query := `SELECT id, name FROM cat_courses`
	rows, err := s.DbPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener cursos: %w", err)
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name); err != nil {
			return nil, fmt.Errorf("error al escanear cursos: %w", err)
		}
		courses = append(courses, course)
	}

	return courses, nil
}

func (s *PgxStorage) GetSubjectsByCourse(ctx context.Context, courseID int) ([]models.Subject, error) {
	query := `SELECT id, key, name, coins FROM academyc_history WHERE course_id = $1`
	rows, err := s.DbPool.Query(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener materias: %w", err)
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		var subject models.Subject
		if err := rows.Scan(&subject.ID, &subject.Key, &subject.Name, &subject.Coins); err != nil {
			return nil, fmt.Errorf("error al escanear materias: %w", err)
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (s *PgxStorage) GetStudents(ctx context.Context) ([]models.Alumno, error) {
	query := `
		SELECT id, name, lastname1, lastname2, course_id
		FROM alumn
	`

	rows, err := s.DbPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener alumnos: %w", err)
	}
	defer rows.Close()

	var alumnos []models.Alumno
	for rows.Next() {
		var alumno models.Alumno
		if err := rows.Scan(&alumno.ID, &alumno.Name, &alumno.Lastname1, &alumno.Lastname2, &alumno.CourseID); err != nil {
			return nil, fmt.Errorf("error al escanear alumnos: %w", err)
		}
		alumnos = append(alumnos, alumno)
	}

	return alumnos, nil
}

func (s *PgxStorage) SeedCatSemesters(ctx context.Context) error {
	query := `
		INSERT INTO cat_semesters (id, name)
		VALUES 
		(1, 'Primer Semestre'),
		(2, 'Segundo Semestre'),
		(3, 'Tercer Semestre'),
		(4, 'Cuarto Semestre'),
		(5, 'Quinto Semestre'),
		(6, 'Sexto Semestre'),
		(7, 'Séptimo Semestre'),
		(8, 'Octavo Semestre')
		ON CONFLICT (id) DO NOTHING;
	`

	_, err := s.DbPool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error al insertar datos en cat_semesters: %w", err)
	}

	return nil
}

func (s *PgxStorage) GetPendingGradesForCurrentSemester(ctx context.Context, alumnID int) ([]models.PendingGrade, error) {
	query := `
		SELECT 
			sc.subject_id,
			ah.name AS subject_name,
			sc.semester_id,
			pg.partial_number
		FROM semester_course sc
		LEFT JOIN partial_grades pg ON sc.id = pg.semester_course_id
		JOIN academyc_history ah ON sc.subject_id = ah.id
		WHERE sc.alumn_id = $1
		  AND sc.semester_id = (
			  SELECT current_semester
			  FROM alumn
			  WHERE id = $1
		  )
		  AND pg.grade IS NULL
		ORDER BY sc.semester_id, sc.subject_id, pg.partial_number;
	`

	rows, err := s.DbPool.Query(ctx, query, alumnID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener calificaciones pendientes: %w", err)
	}
	defer rows.Close()

	var pendingGrades []models.PendingGrade
	for rows.Next() {
		var pendingGrade models.PendingGrade
		if err := rows.Scan(&pendingGrade.SubjectID, &pendingGrade.SubjectName, &pendingGrade.SemesterID, &pendingGrade.PartialNumber); err != nil {
			return nil, fmt.Errorf("error al procesar filas de calificaciones pendientes: %w", err)
		}
		pendingGrades = append(pendingGrades, pendingGrade)
	}

	return pendingGrades, nil
}

func (s *PgxStorage) GetAlumnIDBySemesterCourseID(ctx context.Context, semesterCourseID int, alumnID *int) error {
	query := `
		SELECT sc.alumn_id
		FROM semester_course sc
		WHERE sc.id = $1;
	`
	return s.DbPool.QueryRow(ctx, query, semesterCourseID).Scan(alumnID)
}
