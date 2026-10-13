package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MaminirinaEdwino/turbostack/src/config"
	turbojwt "github.com/MaminirinaEdwino/turbostackapi/src/config"
	"github.com/MaminirinaEdwino/turbostackapi/src/entity"
)

type ProjectService struct{}

// Middleware d'authentification pour extraire le user_id du JWT
func getUserIdFromRequest(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("missing token")
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 {
		return 0, fmt.Errorf("invalid token format")
	}
	payload, err := turbojwt.Verify(JwtSecret, tokenParts[1], 0)
	if err != nil {
		return 0, err
	}
	uidFloat, ok := payload["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid token payload")
	}
	return int(uidFloat), nil
}

func (s *ProjectService) CreateProject(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	pJson := entity.ProjectJSON{
		Nom:         body.Name,
		Description: body.Description,
		Type:        body.Type,
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
	}

	dataJSON, _ := json.Marshal(pJson)

	_, err = config.DB.Exec(
		"INSERT INTO projects (user_id, nom, type_project, description, data) VALUES ($1, $2, $3, $4, $5)",
		userId, body.Name, body.Type, body.Description, dataJSON,
	)
	if err != nil {
		http.Error(w, "Project already exists or database error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Success"})
}

func (s *ProjectService) FetchProjects(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	rows, err := config.DB.Query("SELECT data FROM projects WHERE user_id = $1 ORDER BY update_at DESC", userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var projects []entity.ProjectJSON
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err == nil {
			var p entity.ProjectJSON
			json.Unmarshal(data, &p)
			projects = append(projects, p)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (s *ProjectService) FetchProjectByName(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	projectName := r.PathValue("name")
	var data []byte
	err = config.DB.QueryRow("SELECT data FROM projects WHERE nom = $1 AND user_id = $2", projectName, userId).Scan(&data)
	if err == sql.ErrNoRows {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	var p entity.ProjectJSON
	json.Unmarshal(data, &p)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (s *ProjectService) SaveProject(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var pJson entity.ProjectJSON
	if err := json.NewDecoder(r.Body).Decode(&pJson); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	pJson.UpdateAt = time.Now()
	dataJSON, _ := json.Marshal(pJson)

	_, err = config.DB.Exec(
		"UPDATE projects SET data = $1, update_at = NOW() WHERE nom = $2 AND user_id = $3",
		dataJSON, pJson.Nom, userId,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Success"})
}

func (s *ProjectService) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	projectName := r.PathValue("name")
	_, err = config.DB.Exec("DELETE FROM projects WHERE nom = $1 AND user_id = $2", projectName, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Success"})
}

func (s *ProjectService) UploadAsset(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserIdFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var body struct {
		ProjectName string `json:"project_name"`
		FileName    string `json:"file_name"`
		Base64Image string `json:"base64_image"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	var data []byte
	err = config.DB.QueryRow("SELECT data FROM projects WHERE nom = $1 AND user_id = $2", body.ProjectName, userId).Scan(&data)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	var p entity.ProjectJSON
	json.Unmarshal(data, &p)

	p.Assets = append(p.Assets, entity.AssetJSON{
		FileName:    body.FileName,
		Base64Image: body.Base64Image,
	})

	updatedData, _ := json.Marshal(p)
	config.DB.Exec("UPDATE projects SET data = $1 WHERE nom = $2 AND user_id = $3", updatedData, body.ProjectName, userId)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}
