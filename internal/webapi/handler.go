package webapi

import (
	"errors"
	"fmt"
	"log"
	"mars/internal/app"
	"mars/internal/config"
	"mars/internal/parser"
	"mars/internal/rover"
	"net/http"
)

// Server is a struct that holds the dependencies for the web api
type Server struct {
	cfg     *config.Config
	parser  *parser.Parser
	factory rover.MissionControlFactory
}

const maxRequestSize = 1024 * 1024 // 1MB

var (
	ErrServerStart = errors.New("error starting http server")
)

// NewServer is the constructor for a new web api server
func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:     cfg,
		parser:  parser.New(),
		factory: rover.NewMissionControlFactory(),
	}
}

// Starts begins listening for requests
func (s *Server) Start() error {
	log.Printf("starting Mars rover HTTP server on %s\n", s.cfg.SrvAddr)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /mcontrol", s.handleMission) // register POST endpoint only

	if err := http.ListenAndServe(s.cfg.SrvAddr, mux); err != nil {
		return fmt.Errorf("%w: %v", ErrServerStart, err)
	}

	return nil
}

func (s *Server) handleMission(w http.ResponseWriter, r *http.Request) {

	// limit the size of what we accept
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	defer r.Body.Close()

	application := app.NewApp(s.parser, s.factory, r.Body, w, s.cfg)

	if err := application.Run(); err != nil {
		log.Printf("ERROR: mission failed: %v", err)

		http.Error(w, "mission processing failed", http.StatusInternalServerError)
		return
	}
}
