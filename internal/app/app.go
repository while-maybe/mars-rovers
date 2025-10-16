package app

import (
	"errors"
	"fmt"
	"io"
	"mars/internal/config"
	"mars/internal/rover"
)

type Parser interface {
	Parse(input string) (*rover.Plateau, []rover.RoverInstruction, error)
}

type App struct {
	parser Parser
	mcf    rover.MissionControlFactory
	input  io.Reader
	output io.Writer
	cfg    *config.Config
}

var (
	ErrAppInput       = errors.New("error reading input")
	ErrAppParsing     = errors.New("error parsing input")
	ErrAppCreatingMC  = errors.New("error creating mission control")
	ErrAppExecMission = errors.New("error executing mission")
)

// NewApp takes injects a parser, mission control factory, an io.Reader and an io.Writer returning a new application struct with all its dependencies
func NewApp(p Parser, mcf rover.MissionControlFactory, i io.Reader, o io.Writer, cfg *config.Config) *App {
	return &App{
		parser: p,
		mcf:    mcf,
		input:  i,
		output: o,
		cfg:    cfg,
	}
}

// Run starts the application
func (a *App) Run() error {
	inputBytes, err := io.ReadAll(a.input)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrAppInput, err)
	}

	plateau, instructions, err := a.parser.Parse(string(inputBytes))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrAppParsing, err)
	}

	mc, err := a.mcf.Create(plateau)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrAppCreatingMC, err)
	}

	missionControlInput := &rover.MissionControlInput{
		Instructions: instructions,
	}

	output, err := mc.Execute(missionControlInput)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrAppExecMission, err)
	}

	fmt.Fprintf(a.output, "info: Mission complete. Final rover positions:")
	fmt.Println()
	for _, singleRoverOutput := range output {
		fmt.Fprintln(a.output, singleRoverOutput)
	}
	return nil
}
