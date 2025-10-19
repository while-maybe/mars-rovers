package config

import (
	"errors"
	"flag"
	"fmt"
)

type OpMode int

const (
	DefaultServerAddr = ":8080"
	DefaultMinSizeX   = 2
	DefaultMinSizeY   = 2
)

const (
	ModeUnknown OpMode = iota
	ModeCLI
	ModeWebAPI
)

var (
	ErrFlagsIncompatible = errors.New("cannot use -file and -webapi flags at the same time")
	ErrPlateauDimensions = errors.New("plateau dimensions must be positive")
	ErrServerAddr        = errors.New("server address required for WebAPI mode, leave empty for default address")
)

type Config struct {
	FilePath    string
	MinPlateauX int
	MinPlateauY int
	OpMode      OpMode
	SrvAddr     string
}

// New returns a pointer to a new Config struct from a filePath, minPlateauX and minPlateauY
func New(minPlateauX, minPlateauY int, filePath string, opMode OpMode, srvAddr string) *Config {
	return &Config{
		FilePath:    filePath,
		MinPlateauX: minPlateauX,
		MinPlateauY: minPlateauY,
		OpMode:      opMode,
		SrvAddr:     srvAddr,
	}
}

// Default returns a pointer to a new Config struct with predefined sensible (ModeCLI) defaults
func Default() *Config {
	return &Config{
		MinPlateauX: DefaultMinSizeX,
		MinPlateauY: DefaultMinSizeY,
		OpMode:      ModeCLI,
	}
}

// ParseFlags returns a pointer to a new Config struct from user provided cli flags
func ParseFlags(args []string) (*Config, error) {
	cfg := &Config{}

	flags := flag.NewFlagSet("mars-rovers", flag.ContinueOnError)

	// flags for cli mode
	flags.StringVar(&cfg.FilePath, "file", "", "Input file. If not provided, reads from stdin.")
	flags.IntVar(&cfg.MinPlateauX, "min-size-x", DefaultMinSizeX, "Minimum size X for plateau (optional)")
	flags.IntVar(&cfg.MinPlateauY, "min-size-y", DefaultMinSizeY, "Minimum size Y for plateau (optional)")

	// flags for webapi mode
	webAPIFlag := flag.Bool("webapi", false, "run in webapi server mode")
	flags.StringVar(&cfg.SrvAddr, "addr", DefaultServerAddr, "port for webapi server")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	// assign operating mode based on -webapi flag being present
	if *webAPIFlag {
		if cfg.FilePath != "" {
			return nil, ErrFlagsIncompatible
		}
		cfg.OpMode = ModeWebAPI

	} else {
		cfg.OpMode = ModeCLI
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.MinPlateauX < 1 || c.MinPlateauY < 1 {
		return fmt.Errorf("%w: (got %dx%d)", ErrPlateauDimensions, c.MinPlateauX, c.MinPlateauY)
	}

	if c.OpMode == ModeWebAPI && c.SrvAddr == "" {
		return ErrServerAddr
	}

	return nil
}
