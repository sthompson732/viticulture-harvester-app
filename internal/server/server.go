/*
 * File: server.go
 * Description: Initializes and configures the HTTP server. This module sets up the server with the necessary
 *              middleware, routes, and starts the server. It uses the `gin` framework for routing and middleware
 *              handling.
 * Usage:
 *   - Create an instance of the server using NewServer and pass the configured router.
 *   - Start the server on a specified port using the Start method.
 * Dependencies:
 *   - Gin-Gonic for routing.
 *   - Custom middleware for logging and error handling.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
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

func (s *Server) Start(port string) {
	log.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(":"+port, s.Router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
