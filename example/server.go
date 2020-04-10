package example

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	http.Handler
}

func NewServer() *Server {
	server := new(Server)

	router := http.NewServeMux()
	//this will match all not handled route
	router.Handle("/", http.HandlerFunc(server.rootHandler))
	router.Handle("/example", http.HandlerFunc(server.demoPageHandler))
	router.Handle("/ws/echo", http.HandlerFunc(server.wsEchoHandler))
	router.Handle("/ws/time", http.HandlerFunc(server.wsTimeHandler))

	server.Handler = router

	return server
}
func (s Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (s *Server) demoPageHandler(w http.ResponseWriter, r *http.Request) {
	const demoTemplatePath = "example/demo.html"
	tmpl, err := template.ParseFiles(demoTemplatePath)
	if err != nil {
		log.Fatalf("problem opening %s %v", demoTemplatePath, err)
	}

	const welcomeMsg = "a simple API demo page"
	tmpl.Execute(w, struct{ WelcomeMsg string }{welcomeMsg})
}

func (s *Server) wsEchoHandler(w http.ResponseWriter, r *http.Request) {
	const echoPrefix = "ECHO: "
	ws := newWSServer(w, r)

	for {
		messageType, bytes, err := ws.ReadMessage()
		if err != nil {
			log.Println("ReadMessage Error: ", err)
			break
		}

		msg := echoPrefix + string(bytes)
		err = ws.WriteMessage(messageType, []byte(msg))
		if err != nil {
			log.Println("Write Error: ", err)
			break
		}

	}

	//this well generate close code 1006 at client
	//ws.Close()
	//ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server closed"))
}

func (s *Server) wsTimeHandler(w http.ResponseWriter, r *http.Request) {
	const timeFormat = "15:04:05"

	ws := newWSServer(w, r)
	ws.write([]byte(time.Now().Format(timeFormat)))

	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server closed"))
}
