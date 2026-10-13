package api

import (
	"github.com/MaminirinaEdwino/turbostackapi/src/activation"
	webview "github.com/webview/webview_go"
)


func (s *UserService) SayHello(name string) string {
	return "Bonjour " + name
}

func (s *UserService) SaveToken(token string) string {
	activation.SaveTokenToFile(token)
	return "Token Saved"
}

func (s *UserService) CheckToken() map[string]interface{} {
	return activation.CheckActivationToken()
}

// On implémente l'interface Binder
func (s *UserService) Bind(w webview.WebView) {
	w.Bind("sayHello", s.SayHello)
	w.Bind("saveToken", s.SaveToken)
	w.Bind("checkToken", s.CheckToken)
}
