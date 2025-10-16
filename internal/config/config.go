package config

import "flag"

type Config struct {
	FilePath    string
	MinPlateauX int
	MinPlateauY int
}

// New returns a pointer to a new Config struct from a filePath, minPlateauX and minPlateauY
func New(minPlateauX, minPlateauY int, filePath string) *Config {
	return &Config{
		FilePath:    filePath,
		MinPlateauX: minPlateauX,
		MinPlateauY: minPlateauY,
	}
}

// Default returns a pointer to a new Config struct with predefined sensible defaults
func Default() *Config {
	return &Config{
		FilePath:    "",
		MinPlateauX: 2,
		MinPlateauY: 2,
	}
}

// ParseFlags returns a pointer to a new Config struct from user provided cli flags
func ParseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.FilePath, "file", "", "Input file. If not provided, reads from stdin.")
	flag.IntVar(&cfg.MinPlateauX, "min-size-x", 2, "Minimum size X for plateau (optional)")
	flag.IntVar(&cfg.MinPlateauY, "min-size-y", 2, "Minimum size Y for plateau (optional)")

	flag.Parse()

	return cfg
}
