package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MaminirinaEdwino/turbostackapi/src/activation"
	"github.com/MaminirinaEdwino/turbostackapi/src/api"
	"github.com/MaminirinaEdwino/turbostackapi/src/config"
)

func main() {
	// 1. Initialisation de la base de données PostgreSQL
	config.InitDB()
	activation.CheckActivationToken()

	mux := http.NewServeMux()

	userSvc := &api.UserService{}
	projSvc := &api.ProjectService{}
	sysSvc := &api.SystemService{}

	// Routes Auth & User
	mux.HandleFunc("POST /api/auth/register", userSvc.Register)
	mux.HandleFunc("POST /api/auth/login", userSvc.Login)
	// mux.HandleFunc("POST /api/user/token", userSvc.SaveToken)
	// mux.HandleFunc("GET /api/user/token/check", userSvc.CheckToken)

	// Routes Projets (Liées à l'utilisateur via JWT)
	mux.HandleFunc("POST /api/projects", projSvc.CreateProject)
	mux.HandleFunc("GET /api/projects", projSvc.FetchProjects)
	mux.HandleFunc("GET /api/projects/{name}", projSvc.FetchProjectByName)
	mux.HandleFunc("PUT /api/projects", projSvc.SaveProject)
	mux.HandleFunc("DELETE /api/projects/{name}", projSvc.DeleteProject)
	mux.HandleFunc("POST /api/projects/assets", projSvc.UploadAsset)

	// Routes Système & Librairies
	mux.HandleFunc("GET /api/system/librairie", sysSvc.LoadLibrairie)

	// Middleware CORS global optionnel
	handler := corsMiddleware(mux)

	port := ":8080"
	fmt.Printf("🚀 TurboStack API Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
