package api

import (
	"alumnos/models"
	"alumnos/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type API struct {
	Repo *repository.PgxStorage
}

func NewAPI(repo *repository.PgxStorage) *API {
	return &API{Repo: repo}
}

func (api *API) RegistrarAlumno(w http.ResponseWriter, r *http.Request) {
	// Leer y registrar el cuerpo recibido
	body, _ := io.ReadAll(r.Body)
	fmt.Println("Cuerpo recibido:", string(body))

	// Decodificar el JSON
	var request models.RegisterAlumnRequest
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Error al decodificar la solicitud: %v", err), http.StatusBadRequest)
		return
	}
	// Log del contenido decodificado
	fmt.Printf("Solicitud decodificada: %+v\n", request)

	// Validar campos requeridos
	if request.Name == "" || request.Lastname1 == "" || request.CourseID == 0 || request.CurrentCourseID == 0 {
		http.Error(w, "Faltan campos requeridos (name, lastname1, course_id, current_course_id)", http.StatusBadRequest)
		return
	}

	// Registrar alumno
	alumnoID, err := api.Repo.RegisterAlumn(r.Context(), request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al registrar alumno: %v", err), http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Alumno registrado exitosamente",
		"alumno_id": alumnoID,
	}); err != nil {
		http.Error(w, fmt.Sprintf("Error al codificar la respuesta: %v", err), http.StatusInternalServerError)
	}
}

func (api *API) RegistrarEnSemestre(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AlumnoID   int   `json:"alumno_id"`
		SemesterID int   `json:"semester_id"`
		SubjectIDs []int `json:"subject_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Error al decodificar la solicitud: %v", err), http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if input.AlumnoID == 0 || input.SemesterID == 0 || len(input.SubjectIDs) == 0 {
		http.Error(w, "Alumno ID, Semester ID y Subject IDs son obligatorios", http.StatusBadRequest)
		return
	}

	// Llama al método del repositorio
	err := api.Repo.RegistrarEnSemestreConMaterias(r.Context(), input.AlumnoID, input.SemesterID, input.SubjectIDs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al registrar en semestre: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Alumno registrado en el semestre exitosamente",
	})
}

func (api *API) RegistrarCalificacionParcial(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SemesterCourseID int     `json:"semester_course_id"`
		PartialNumber    int     `json:"partial_number"`
		Grade            float64 `json:"grade"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Error al decodificar la solicitud: %v", err), http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if input.SemesterCourseID == 0 || input.PartialNumber == 0 || input.Grade < 0 {
		http.Error(w, "SemesterCourseID, PartialNumber y Grade son obligatorios y deben ser válidos", http.StatusBadRequest)
		return
	}

	// Llama al método del repositorio
	err := api.Repo.RegistrarCalificacionParcial(r.Context(), input.SemesterCourseID, input.PartialNumber, input.Grade)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al registrar calificación parcial: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Calificación parcial registrada exitosamente",
	})
}

func (api *API) GenerarCalificacionesAgrupadas(w http.ResponseWriter, r *http.Request) {
	// Estructura para decodificar el cuerpo de la solicitud
	var input struct {
		AlumnoID int `json:"alumno_id"`
	}

	// Decodificar el cuerpo JSON
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Error al decodificar la solicitud: %v", err), http.StatusBadRequest)
		return
	}

	// Validar que el ID del alumno sea válido
	if input.AlumnoID == 0 {
		http.Error(w, "El parámetro 'alumno_id' es obligatorio y debe ser válido", http.StatusBadRequest)
		return
	}

	// Llama al método del repositorio
	semestres, promedioFinal, err := api.Repo.GenerarCalificacionesAgrupadasPorSemestre(r.Context(), input.AlumnoID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al generar calificaciones: %v", err), http.StatusInternalServerError)
		return
	}

	// Respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"promedio_final": promedioFinal,
		"semestres":      semestres,
	})
}

func (api *API) GetCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := api.Repo.GetCourses(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener cursos: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(courses)
}

func (api *API) GetSubjectsByCourse(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la solicitud
	var input struct {
		CourseID int `json:"course_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Error al decodificar la solicitud: %v", err), http.StatusBadRequest)
		return
	}

	// Validar que el CourseID sea válido
	if input.CourseID <= 0 {
		http.Error(w, "El campo 'course_id' debe ser un número positivo", http.StatusBadRequest)
		return
	}

	// Obtener materias por CourseID desde el repositorio
	subjects, err := api.Repo.GetSubjectsByCourse(r.Context(), input.CourseID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener materias: %v", err), http.StatusInternalServerError)
		return
	}

	// Responder con JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subjects)
}

func (api *API) GetStudents(w http.ResponseWriter, r *http.Request) {
	// Obtener alumnos desde el repositorio
	alumnos, err := api.Repo.GetStudents(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener alumnos: %v", err), http.StatusInternalServerError)
		return
	}

	// Responder con JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(alumnos)
}
