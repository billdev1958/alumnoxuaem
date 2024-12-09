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
	Name            string      `json:"name"`
	Lastname1       string      `json:"lastname1"`
	Lastname2       string      `json:"lastname2"`
	CourseID        int         `json:"course_id"`
	CurrentCourseID int         `json:"current_course_id"`
	Subjects        []SubjectID `json:"subjects"`
}

type SubjectID struct {
	ID int `json:"id"`
}
