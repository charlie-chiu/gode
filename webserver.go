package go_ws

import (
	"html/template"
	"log"
	"net/http"
)

type WSServer struct {
	http.Handler
}

func NewWSServer() *WSServer {
	server := new(WSServer)

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const demoTemplatePath = "demo.html"
		tmpl, err := template.ParseFiles(demoTemplatePath)
		if err != nil {
			log.Fatalf("problem opening %s %v", demoTemplatePath, err)
		}

		const welcomeMsg = "a simple API demo page"
		tmpl.Execute(w, struct{ WelcomeMsg string }{welcomeMsg})
	}))

	server.Handler = router

	return server
}
