/*
 * server.go: Initializes and manages the HTTP server.
 * Sets up routing and starts listening for requests.
 * Usage: Provides the runtime environment for the web interface.
 * Author(s): Shannon Thompson
 * Created on: 04/11/2024
 */

package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func NewServer(router *mux.Router) *Server {
	return &Server{
		Router: router,
	}
}

func (s *Server) Start(port string) error {
	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(":"+port, s.Router); err != nil {
		log.Printf("Server failed to start: %v", err)
		return err
	}
	return nil
}
