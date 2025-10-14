package rover

import (
	"errors"
	"fmt"
	"log"
)

type Direction int

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

const (
	UnknownDirection Direction = iota
	N                          // North
	E                          // East
	S                          // South
	W                          // West
)

const (
	CmdMove  Command = 'M' // Move
	CmdLeft  Command = 'L' // Left
	CmdRight Command = 'R' // Right
)

var (
	minPlateauX = 2
	minPlateauY = 2
)

type Coordinates struct {
	x int
	y int
}

type Position struct {
	coordinates Coordinates
	direction   Direction
}

type Rover struct {
	id       int
	position *Position
}

type Plateau struct {
	maxX int
	maxY int
}

type RoverInstruction struct {
	InitialPosition *Position
	Commands        string
}

type MissionControlInput struct {
	Instructions []RoverInstruction
}

type MissionControl struct {
	plateau         *Plateau
	occupiedSquares map[Coordinates]int // contains the position of an existing (not moving) rover as the key. Value is the rover ID
}

type MissionControlFactory interface {
	Create(plateau *Plateau) (*MissionControl, error)
}

type defaultMissionControlFactory struct{}

func NewMissionControlFactory() *defaultMissionControlFactory {
	return &defaultMissionControlFactory{}
}

func (f *defaultMissionControlFactory) Create(plateau *Plateau) (*MissionControl, error) {
	return NewMissionControl(plateau)
}

// NewCoordinates takes a pair of int x, y coordinates and returns a coordinates struct and performs no validation
func NewCoordinates(x, y int) Coordinates {
	return Coordinates{
		x: x,
		y: y,
	}
}

type Command rune

func (d Direction) validate() error {
	if d == UnknownDirection || d < N || d > W {
		return ErrDirectionUnknown
	}
	return nil
}

// NewPosition takes a plateau, coordinates and a direction and returns a pointer to a position or an error if the given arguments don't pass validation (plateau area too small or if coordinates are out of bounds)
func NewPosition(p *Plateau, c Coordinates, d Direction) (*Position, error) {
	pos := &Position{
		coordinates: c,
		direction:   d,
	}

	if err := pos.validate(p); err != nil {
		return nil, err
	}

	if err := d.validate(); err != nil {
		return nil, err
	}

	return pos, nil
}

// String implements the Stringer interface
func (p *Position) String() string {
	return fmt.Sprintf("%d %d %s", p.coordinates.x, p.coordinates.y, p.direction)
}

// NewRover takes an id and a Position and returns or a pointer to a Rover object or error if the given position is nil
func NewRover(id int, p *Position) (*Rover, error) {
	if p == nil {
		return nil, ErrRoverPositionIsNil
	}

	return &Rover{
		id:       id,
		position: p,
	}, nil
}

// move returns the resulting Position of applying movement to the Rover in the direction it's currently facing
func (r *Rover) move() Position {
	// we do not mutate the original
	nextPosition := *r.position
	switch r.position.direction {
	case N:
		nextPosition.coordinates.y++
	case E:
		nextPosition.coordinates.x++
	case S:
		nextPosition.coordinates.y--
	case W:
		nextPosition.coordinates.x--
	}
	return nextPosition
}

// turnLeft causes the Rover to rotate 90 degrees on itself to the left
func (r *Rover) turnLeft() {
	if r.position.direction == N {
		r.position.direction = W
		return
	}
	r.position.direction--
}

// turnRight causes the Rover to rotate 90 degrees on itself to the right
func (r *Rover) turnRight() {
	if r.position.direction == W {
		r.position.direction = N
		return
	}
	r.position.direction++
}

// validateBoundaries is a helper used through the rest of the code. It takes a pointer to a Position and a pointer to a Plateau returning an error should the position be out of bounds for the given plateau
func validateBoundaries(pos *Position, plateau *Plateau) error {
	if pos.coordinates.x < 0 || pos.coordinates.x > plateau.maxX || pos.coordinates.y < 0 || pos.coordinates.y > plateau.maxY {
		return ErrPositionOutOfBounds
	}
	return nil
}

// validate is a helper that wraps validateBoundaries
func (pos *Position) validate(plateau *Plateau) error {
	return validateBoundaries(pos, plateau)
}

// set updates the Position's coordinates and direction
func (p *Position) set(newPos Position) {
	p.coordinates.x = newPos.coordinates.x
	p.coordinates.y = newPos.coordinates.y
	p.direction = newPos.direction
}

// NewPlateau returns a pointer to a new Plateau, validating against minimum dimensions and returning an error accordingly
func NewPlateau(maxX, maxY int) (*Plateau, error) {
	if maxX < minPlateauX || maxY < minPlateauY {
		return nil, ErrPlateauTooSmall
	}

	return &Plateau{
		maxX: maxX,
		maxY: maxY,
	}, nil
}

// NewMissionControl takes a pointer to a Plateau struct and returns a pointer to a new MissionControl struct returning an error should the given Plateau be nil
func NewMissionControl(p *Plateau) (*MissionControl, error) {
	if p == nil {
		return nil, ErrPlateauIsNil
	}

	return &MissionControl{
		plateau:         p,
		occupiedSquares: make(map[Coordinates]int),
	}, nil
}

// validate takes a Position pointer and returns an error should the desired Position fail validation or if another rover is already at that Position
func (mc *MissionControl) validate(pos *Position) error {
	if err := validateBoundaries(pos, mc.plateau); err != nil {
		return err
	}

	if _, ok := mc.occupiedSquares[pos.coordinates]; ok {
		return ErrRoverCollision
	}

	return nil
}

// RunRover takes a Rover pointer and a command string, returning a feedback string and an error should the commands fail. It keeps track of previous placed Rover in the Plateau and processes the commands giving feedback  to the user
func (mc *MissionControl) RunRover(r *Rover, commands string) (string, error) {
	// check to see if mission control is attempting to place a rover on a location that's occupied
	if err := mc.validate(r.position); err != nil {
		// original error remains wrapped
		return "", fmt.Errorf("new rover with id %d cannot be placed at (%s): %w", r.id, r.position.String(), err)
	}

	// place an entry in the occupied map using x, y coordinates as key and rover id as the value
	mc.occupiedSquares[r.position.coordinates] = r.id

	// process commands
	for _, c := range commands {
		switch Command(c) {
		case CmdLeft:
			r.turnLeft()
		case CmdRight:
			r.turnRight()
		case CmdMove:
			// store current position before moving
			currentPosKey := r.position.coordinates

			nextPos := r.move()

			// handle invalid moves
			if err := mc.validate(&nextPos); err != nil {
				// this is an invalid move so it will be ignored and we carry on attempting remaining commands
				log.Printf("WARN: Rover %d ignored move to (%v): %s", r.id, nextPos.String(), err.Error())
				continue
			}

			r.position.set(nextPos)

			// delete existing state from the map after rover moves
			delete(mc.occupiedSquares, currentPosKey)

			// and update with new position here
			mc.occupiedSquares[nextPos.coordinates] = r.id
		}
	}

	return r.position.String(), nil
}

// implement stringer interface so we can print a friendly direction when using a Print function
func (d Direction) String() string {
	switch d {
	case N:
		return "N"
	case E:
		return "E"
	case S:
		return "S"
	case W:
		return "W"
	default:
		return "?" // should never happen
	}
}

func (mc *MissionControl) Execute(input *MissionControlInput) ([]string, error) {
	var output []string

	for i, instruction := range input.Instructions {
		roverID := i + 1

		currentRover, err := NewRover(roverID, instruction.InitialPosition)
		if err != nil {
			return nil, fmt.Errorf("%w %d: %v", ErrRoverCreating, roverID, err)
		}

		singleRoverOutput, err := mc.RunRover(currentRover, instruction.Commands)
		if err != nil {
			return nil, fmt.Errorf("%w %d: %v", ErrRoverInstructions, roverID, err)
		}

		output = append(output, singleRoverOutput)
	}

	return output, nil
}
