package rover

import "errors"

var (
	ErrPositionOutOfBounds = errors.New("position must be more than 0 and within boundaries")
	ErrDirectionUnknown    = errors.New("direction must be one of N, E, S, W")
	ErrRoverPositionIsNil  = errors.New("rover must not be nil")
	ErrRoverCollision      = errors.New("path is blocked by another rover")
	ErrRoverInstructions   = errors.New("rover error executing instruction")
	ErrRoverCreating       = errors.New("rover could not be created")
	ErrPlateauTooSmall     = errors.New("plateau must be at least 2 * 2")
	ErrPlateauIsNil        = errors.New("plateau must not be nil")
)
