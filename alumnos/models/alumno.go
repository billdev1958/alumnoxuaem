package models

import "time"

type Alumno struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Lastname1       string    `json:"lastname1"`
	Lastname2       string    `json:"lastname2,omitempty"` // omitempty si puede ser nulo
	CourseID        int       `json:"course_id"`
	CurrentCourseID int       `json:"current_course_id"` // correlación con current_semester
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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
