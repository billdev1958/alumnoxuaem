package models

import "time"

type Alumno struct {
	ID              int
	Name            string
	Lastname1       string
	Lastname2       string
	CourseID        int
	CurrentCourseID int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type RegisterAlumnRequest struct {
	Name            string
	Lastname1       string
	Lastname2       string
	CourseID        int
	CurrentCourseID int
	Subjects        []SubjectID
}

type SubjectID struct {
	ID int
}
