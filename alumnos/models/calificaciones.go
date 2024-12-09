package models

import "time"

type CalificacionParcial struct {
	PartialNumber int     `json:"partial_number"`
	Grade         float64 `json:"grade"`
}

type MateriaCalificaciones struct {
	SubjectID   int                   `json:"subject_id"`
	SubjectName string                `json:"subject_name"`
	Parciales   []CalificacionParcial `json:"parciales"`
	Promedio    float64               `json:"promedio"` // Promedio de la materia
}

type SemestreCalificaciones struct {
	SemesterID   int                     `json:"semester_id"`
	SemesterName string                  `json:"semester_name"`
	Materias     []MateriaCalificaciones `json:"materias"`
	Promedio     float64                 `json:"promedio"` // Promedio del semestre
}

type Course struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Subject struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Name  string `json:"name"`
	Coins int    `json:"coins"`
}

type PendingGrade struct {
	SubjectID     int    `json:"subject_id"`
	SubjectName   string `json:"subject_name"`
	SemesterID    int    `json:"semester_id"`
	PartialNumber int    `json:"partial_number"`
}

type SemesterCourse struct {
	ID           int       `json:"id"`
	AlumnID      int       `json:"alumn_id"`
	SemesterID   int       `json:"semester_id"`
	SemesterName string    `json:"semester_name"`
	SubjectID    int       `json:"subject_id"`
	SubjectName  string    `json:"subject_name"`
	FinalGrade   *float64  `json:"final_grade,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	// Si quieres incluir las calificaciones parciales dentro del objeto:
	PartialGrades []PartialGrade `json:"partial_grades,omitempty"`
}

type PartialGrade struct {
	ID               int       `json:"id"`
	SemesterCourseID int       `json:"semester_course_id"`
	PartialNumber    int       `json:"partial_number"`
	Grade            float64   `json:"grade"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CatSemester struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SemesterGrades struct {
	ID                 int       `json:"id"`                   // ID único del registro
	AlumnID            int       `json:"alumn_id"`             // ID del alumno
	SemesterID         int       `json:"semester_id"`          // ID del semestre
	FinalSemesterGrade float64   `json:"final_semester_grade"` // Promedio final del semestre
	CreatedAt          time.Time `json:"created_at"`           // Fecha de creación del registro
	UpdatedAt          time.Time `json:"updated_at"`           // Última fecha de actualización
}
