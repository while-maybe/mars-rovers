package parser

import "errors"

var (
	ErrParseInvalidFormat    = errors.New("must have a plateau line first and pairs of rover lines")
	ErrParsePlateauFormat    = errors.New("wrong plateau element count, must be X Y")
	ErrParsePositionFormat   = errors.New("wrong rover position element count, must be x y direction")
	ErrParsePlateauX         = errors.New("invalid plateau width")
	ErrParsePlateauY         = errors.New("invalid plateau height")
	ErrParsePositionX        = errors.New("invalid position given for X coordinate")
	ErrParsePositionY        = errors.New("invalid position given for Y coordinate")
	ErrParseInvalidDirection = errors.New("invalid direction given, must be N, E, S, W")
	ErrParseInvalidCommand   = errors.New("invalid command character given, must be L, R, M")
)
