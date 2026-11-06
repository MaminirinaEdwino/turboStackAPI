package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MaminirinaEdwino/turbostackapi/src/activation"
	"github.com/MaminirinaEdwino/turbostackapi/src/api"
	"github.com/MaminirinaEdwino/turbostackapi/src/chatbot"
	"github.com/MaminirinaEdwino/turbostackapi/src/config"
)

func main() {
	// 1. Initialisation de la base de données PostgreSQL
	config.InitDB()
	activation.CheckActivationToken()
	chatbot.Reseaux, chatbot.MotsUniques, chatbot.TagsIntents, _ = chatbot.ChargerModele()
	mux := http.NewServeMux()

	userSvc := &api.UserService{}
	projSvc := &api.ProjectService{}
	// sysSvc := &api.SystemService{}

	// Routes Auth & User
	mux.HandleFunc("POST /api/auth/register", userSvc.Register)
	mux.HandleFunc("POST /api/auth/login", userSvc.Login)
	mux.HandleFunc("POST /api/user/token", userSvc.SaveToken)
	mux.HandleFunc("GET /api/user/token/check", userSvc.CheckToken)
	mux.HandleFunc("POST /api/auth/admin", userSvc.LoginAdmin)

	// Routes Projets (Liées à l'utilisateur via JWT)
	mux.HandleFunc("POST /api/projects", projSvc.CreateProject)
	mux.HandleFunc("GET /api/projects", projSvc.FetchProjects)
	mux.HandleFunc("GET /api/projects/{name}", projSvc.FetchProjectByName)
	mux.HandleFunc("PUT /api/projects", projSvc.SaveProject)
	mux.HandleFunc("DELETE /api/projects/{name}", projSvc.DeleteProject)
	mux.HandleFunc("POST /api/projects/assets", projSvc.UploadAsset)
	mux.HandleFunc("GET /api/librairie/page", projSvc.LoadPageLib)
	mux.HandleFunc("GET /api/librairie/component", projSvc.LoadComponentLib)
	mux.HandleFunc("GET /api/librairie/style", projSvc.LoadStyleLib)
	mux.HandleFunc("GET /api/all/projects", projSvc.FetchProjectsAll)
	mux.HandleFunc("GET /api/project/{id}", projSvc.FetchProjectsById)
	mux.HandleFunc("GET /api/users/{id}", projSvc.FetchUsersById)
	mux.HandleFunc("GET /api/users", projSvc.FetchUsers)
	mux.HandleFunc("POST /api/chat", projSvc.Chat)
	mux.HandleFunc("POST /api/fetchprojectfile", projSvc.FetchProjectFiles)
	mux.HandleFunc("POST /api/getfilecontent", projSvc.GetFileContent)
	mux.HandleFunc("POST /api/export", projSvc.ExportProject)

	// Routes Système & Librairies

	// Middleware CORS global optionnel
	handler := RequestLoggerMiddleware(corsMiddleware(mux))

	port := ":8080"
	fmt.Printf("🚀 TurboStack API Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Récupération de l'adresse IP du client
		// (gère également les cas où l'application est derrière un reverse proxy comme Nginx via X-Forwarded-For)
		clientIP := r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}

		method := r.Method
		uri := r.RequestURI

		// Affichage formaté dans la console du serveur TurboStack
		log.Printf("[TURBOSTACK-LOG] IP: %s | Méthode: %s | URI: %s", clientIP, method, uri)

		// Transmet la requête au prochain handler de la chaîne
		next.ServeHTTP(w, r)

		// Optionnel : Mesurer le temps de traitement de la requête
		duration := time.Since(start)
		_ = duration
	})
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
