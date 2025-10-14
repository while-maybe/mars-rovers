package app

import (
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

func NewApp(p Parser, mcf rover.MissionControlFactory, i io.Reader, o io.Writer) *App {
	return &App{
		parser: p,
		mcf:    mcf,
		input:  i,
		output: o,
	}
}

func (a *App) Run() error {
	inputBytes, err := io.ReadAll(a.input)
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	plateau, instructions, err := a.parser.Parse(string(inputBytes))
	if err != nil {
		return fmt.Errorf("error parsing input: %w", err)
	}

	mc, err := a.mcf.Create(plateau)
	if err != nil {
		return fmt.Errorf("error creating mission control: %w", err)
	}

	missionControlInput := &rover.MissionControlInput{
		Instructions: instructions,
	}

	output, err := mc.Execute(missionControlInput)
	if err != nil {
		return fmt.Errorf("error executing mission: %w", err)
	}

	fmt.Fprintf(a.output, "info: Mission complete. Final rover positions:")
	fmt.Println()
	for _, singleRoverOutput := range output {
		fmt.Fprintln(a.output, singleRoverOutput)
	}
	return nil
}
