package models

import "time"

type Alumno struct {
	ID        int
	Name      string
	Lastname1 string
	Lastname2 string
	CourseID  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
