package api

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/MaminirinaEdwino/turbostackapi/src/entity"
)

type SystemService struct{}

func (s *SystemService) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]string{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (s *SystemService) LoadLibrairie(w http.ResponseWriter, r *http.Request) {
	librairie := entity.LoadLibrairie()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(librairie)
}
