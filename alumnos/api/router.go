package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, apiInstance *API) {
	// Rutas para alumnos
	mux.Handle("POST /v1/alumnos", http.HandlerFunc(apiInstance.RegistrarAlumno))

	// Rutas para semestres y materias
	mux.Handle("POST /v1/semestres", http.HandlerFunc(apiInstance.RegistrarEnSemestre))

	mux.Handle("POST /v1/alumnos/pending-grades", http.HandlerFunc(apiInstance.GetPendingGradesHandler))

	// Rutas para calificaciones parciales
	mux.Handle("POST /v1/calificaciones/parcial", http.HandlerFunc(apiInstance.RegistrarCalificacionParcial))

	// Rutas para generar calificaciones agrupadas
	mux.Handle("POST /v1/calificaciones/agrupadas", http.HandlerFunc(apiInstance.GenerarCalificacionesAgrupadas))

	mux.Handle("POST /v1/courses/subjects", http.HandlerFunc(apiInstance.GetSubjectsByCourse))

	mux.Handle("GET /v1/courses", http.HandlerFunc(apiInstance.GetCourses))

	mux.Handle("GET /v1/students", http.HandlerFunc(apiInstance.GetStudents))

	mux.Handle("POST /v1/semester-courses", http.HandlerFunc(apiInstance.GetSemesterCoursesByAlumnId))

	mux.Handle("GET /v1/semesters", http.HandlerFunc(apiInstance.GetCatSemesters))

}
