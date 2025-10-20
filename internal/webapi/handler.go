package webapi

import (
	"errors"
	"fmt"
	"log"
	"mars/internal/app"
	"mars/internal/config"
	"mars/internal/rover"
	"net/http"
)

// Server is a struct that holds the dependencies for the web api
type Server struct {
	cfg     *config.Config
	parser  app.Parser
	factory rover.MissionControlFactory
}

const maxRequestSize = 1024 * 1024 // 1MB

var (
	ErrServerStart       = errors.New("error starting http server")
	ErrMissionProcessing = errors.New("mission processing failed")
)

// NewServer is the constructor for a new web api server
func NewServer(cfg *config.Config, p app.Parser, mcf rover.MissionControlFactory) *Server {
	return &Server{
		cfg:     cfg,
		parser:  p,
		factory: mcf,
	}
}

// Handler creates and returns a router
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /mcontrol", s.handleMission) // register POST endpoint only

	return mux
}

// Starts begins listening for requests
func (s *Server) Start() error {

	log.Printf("starting Mars rover HTTP server on %s\n", s.cfg.SrvAddr)

	if err := http.ListenAndServe(s.cfg.SrvAddr, s.Handler()); err != nil {
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

		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &maxBytesError):
			http.Error(w, "Request body is too large.", http.StatusRequestEntityTooLarge)

		case errors.Is(err, app.ErrAppParsing):
			http.Error(w, fmt.Sprintf("Bad request: %v", err), http.StatusBadRequest)

		default:
			http.Error(w, "An internal server error occurred.", http.StatusInternalServerError)
		}
		return
	}
}
