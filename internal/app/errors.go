package app

import "errors"

var (
	ErrAppInput       = errors.New("error reading input")
	ErrAppParsing     = errors.New("error parsing input")
	ErrAppCreatingMC  = errors.New("error creating mission control")
	ErrAppExecMission = errors.New("error executing mission")
)
