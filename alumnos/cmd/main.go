package main

import (
	"alumnos/api"
	"alumnos/repository"
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// DSN directo en el código
	dsn := "postgresql://root:root@db:5432/alumnos?sslmode=disable"

	// Conexión a la base de datos
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Printf("Error al conectar a la base de datos: %v\n", err)
		return
	}
	defer dbPool.Close()

	fmt.Println("Conexión a la base de datos exitosa")

	// Inicializar repositorio y API
	repo := repository.NewPgxStorage(dbPool)
	apiInstance := api.NewAPI(repo)

	// Configurar enrutador
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, apiInstance)

	// Ejecutar seed
	if err := repo.SeedCatCourses(context.Background()); err != nil {
		fmt.Printf("Error al ejecutar el seed: %v\n", err)
		return
	}
	fmt.Println("Seed ejecutado exitosamente")

	// Ejecutar el seed
	if err := repo.SeedAcademycHistory(context.Background()); err != nil {
		fmt.Printf("Error al ejecutar el seed: %v\n", err)
		return
	}
	fmt.Println("Seed de academyc_history ejecutado exitosamente.")

	if err := repo.SeedCatSemesters(context.Background()); err != nil {
		fmt.Printf("Error al ejecutar el seed: %v\n", err)
		return
	}
	fmt.Println("Seed de cat_semesters ejecutado exitosamente.")

	// Iniciar servidor
	port := "8080"
	fmt.Printf("Servidor escuchando en :%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
