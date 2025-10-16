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
	"os"
)

func main() {
	cfg := config.ParseFlags()

	inputReader, cleanup, err := getInputReader(cfg)
	if err != nil {
		log.Fatalf("FATAL: %v", err)
	}
	defer cleanup()

	bufferedReader := bufio.NewReader(inputReader)

	if err := run(bufferedReader, cfg); err != nil {
		log.Fatalf("FATAL: Application failed: %v", err)
	}
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
		fmt.Fprintln(os.Stderr, "Options:")
		fmt.Fprintln(os.Stderr, "  -min-size-x int  Minimum size X for plateau")
		fmt.Fprintln(os.Stderr, "  -min-size-y int  Minimum size Y for plateau")

		return nil, noOpCleanup, fmt.Errorf("no input source provided")
	}

	return os.Stdin, noOpCleanup, nil
}

func run(reader io.Reader, cfg *config.Config) error {
	p := parser.New()
	mcf := rover.NewMissionControlFactory()

	app := app.NewApp(p, mcf, reader, os.Stdout, cfg)

	return app.Run()
}
