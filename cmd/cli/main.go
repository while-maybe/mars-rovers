package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"mars/internal/app"
	"mars/internal/config"
	"mars/internal/parser"
	"mars/internal/rover"
	"mars/internal/webapi"
	"os"
)

func main() {
	// parse cmd line flags
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// dispatch to the right runner function
	switch cfg.OpMode {
	case config.ModeCLI:
		if err := runCLI(cfg); err != nil {
			log.Fatalf("FATAL: CLI mode failed: %v", err)
		}

	case config.ModeWebAPI:
		if err := runWebAPI(cfg); err != nil {
			log.Fatalf("FATAL: Web API mode failed: %v", err)
		}

	default:
		log.Fatalf("FATAL: Unknown operating mode configured.")
	}
}

func runCLI(cfg *config.Config) error {
	inputReader, cleanup, err := getInputReader(cfg)
	if err != nil {
		return fmt.Errorf("failed to get input: %w", err)
	}
	defer cleanup()

	// if user accidentally specifies a very large file a buffered reader should handle this
	bufferedReader := bufio.NewReader(inputReader)

	p := parser.New()
	mcf := rover.NewMissionControlFactory()

	app := app.NewApp(p, mcf, bufferedReader, os.Stdout, cfg)
	return app.Run()
}

func runWebAPI(cfg *config.Config) error {
	server := webapi.NewServer(cfg)

	return server.Start()
}

func getInputReader(cfg *config.Config) (io.Reader, func(), error) {
	noOpCleanup := func() {}

	// check to see if a FilePath has been provided
	if cfg.FilePath != "" {

		file, err := os.Open(cfg.FilePath)

		if err != nil {
			return nil, noOpCleanup, fmt.Errorf("could not open file: %w", err)
		}
		return file, func() { file.Close() }, nil
	}

	// check if there is stdin data
	stat, err := os.Stdin.Stat()

	if err != nil {
		return nil, noOpCleanup, fmt.Errorf("could not stat stdin: %w", err)
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "Usage: program -file <path> OR echo 'data' | program")
		fmt.Fprintln(os.Stderr, "\nNo input provided. Use -file flag or pipe data to stdin.")
		return nil, noOpCleanup, fmt.Errorf("no input source provided")
	}

	return os.Stdin, noOpCleanup, nil
}
