package models

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
