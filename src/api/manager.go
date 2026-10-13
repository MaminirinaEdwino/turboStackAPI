package api

import webview "github.com/webview/webview_go"

// Binder est l'interface pour tout module qui expose des fonctions au JS
type Binder interface {
	Bind(w webview.WebView)
}

// Manager centralise tous les binders
type Manager struct {
	binders []Binder
}

func NewManager() *Manager {
	return &Manager{}
}
func (m *Manager) Add(b Binder) {
	m.binders = append(m.binders, b)
}
func (m *Manager) RegisterAll(w webview.WebView) {
	for _, b := range m.binders {
		b.Bind(w)
	}
}

