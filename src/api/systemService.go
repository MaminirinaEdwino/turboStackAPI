package api

import (
	"runtime"

	"github.com/MaminirinaEdwino/turbostackapi/src/entity"
	webview "github.com/webview/webview_go"
)

type SystemService struct{}

func (s *SystemService) GetStats() map[string]string {
	return map[string]string{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	}
}

func (s *SystemService) LoadLibrairie() entity.LibrairieList {
	librairie := entity.LoadLibrairie()
	return librairie
}
func (s *SystemService) Bind(w webview.WebView) {
	w.Bind("getStats", s.GetStats)
	w.Bind("loadLibrairie", s.LoadLibrairie)
}
