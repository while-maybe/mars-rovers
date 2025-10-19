package config

import "errors"

var (
	ErrParserFlagsIncompatible = errors.New("cannot use -file and -webapi flags at the same time")
	ErrParserPlateauDimensions = errors.New("plateau dimensions must be positive")
	ErrParserServerAddr        = errors.New("server address required for WebAPI mode, leave empty for default address")
	ErrParserInvalidValue      = errors.New("invalid values given to parser")
	ErrParserNilConfig         = errors.New("config must not be nil")
	ErrParserModeUnknown       = errors.New("operating mode must be valid")
)
