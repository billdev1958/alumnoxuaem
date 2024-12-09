package models

import "time"

type CalificacionParcial struct {
	PartialNumber int     `json:"partial_number"`
	Grade         float64 `json:"grade"`
}

type MateriaCalificaciones struct {
	SubjectID int                   `json:"subject_id"`
	Parciales []CalificacionParcial `json:"parciales"`
	Promedio  float64               `json:"promedio"`
}

type SemestreCalificaciones struct {
	SemesterID int                     `json:"semester_id"`
	Materias   []MateriaCalificaciones `json:"materias"`
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
	ID         int       `json:"id"`
	AlumnID    int       `json:"alumn_id"`
	SemesterID int       `json:"semester_id"`
	SubjectID  int       `json:"subject_id"`
	FinalGrade float64   `json:"final_grade"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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
