package parser

import (
	"errors"
	"fmt"
	"mars/internal/rover"

	"strconv"
	"strings"
)

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

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(input string) (*rover.Plateau, []rover.RoverInstruction, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	// reject inputs that are not one plateau line + n * pair of instruction lines (a pair per rover with a min of 1 pair)
	if len(lines) < 3 || (len(lines)-1)%2 != 0 {
		return nil, nil, ErrParseInvalidFormat
	}

	// parse plateau
	plateau, err := parsePlateauLine(lines[0])
	if err != nil {
		return nil, nil, err
	}

	// parse rover instructions
	instructions := make([]rover.RoverInstruction, 0, (len(lines)-1)/2)
	for i := 1; i < len(lines); i += 2 {
		positionLine := lines[i]
		commandsLine := lines[i+1]

		position, err := parsePositionLine(positionLine, plateau)
		if err != nil {
			return nil, nil, err
		}

		cmds, err := parseCommandsLine(commandsLine)
		if err != nil {
			return nil, nil, err
		}

		instruction := rover.RoverInstruction{
			InitialPosition: position,
			Commands:        cmds,
		}

		instructions = append(instructions, instruction)
	}

	return plateau, instructions, nil
}

// parsePlateauLine takes a string and returns a Plateau pointer or an error if the given data is not a line of pair of integers
func parsePlateauLine(line string) (*rover.Plateau, error) {
	parts := strings.Fields(strings.TrimSpace(line))

	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: want 2 elements, got %d", ErrParsePlateauFormat, len(parts))
	}

	maxX, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %v %v", ErrParsePlateauX, parts[0], err)
	}

	maxY, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w: %v %v", ErrParsePlateauY, parts[1], err)
	}

	return rover.NewPlateau(maxX, maxY)
}

// parsePositionLine
func parsePositionLine(line string, plateau *rover.Plateau) (*rover.Position, error) {
	parts := strings.Fields(line)

	if len(parts) != 3 {
		return nil, fmt.Errorf("%w: want 3 elements, got %d", ErrParsePositionFormat, len(parts))
	}

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParsePositionX, err)
	}

	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParsePositionY, err)
	}

	dir, err := parseDirection(parts[2])
	if err != nil {
		return nil, err
	}

	coords := rover.NewCoordinates(x, y)
	return rover.NewPosition(plateau, coords, dir)
}

// parseDirection
func parseDirection(dir string) (rover.Direction, error) {
	// make it case-insensitive as a convenience feature
	switch strings.ToUpper(strings.TrimSpace(dir)) {
	case "N":
		return rover.N, nil
	case "E":
		return rover.E, nil
	case "S":
		return rover.S, nil
	case "W":
		return rover.W, nil
	}
	return rover.UnknownDirection, fmt.Errorf("%w: given %s", ErrParseInvalidDirection, dir)
}

// parseCommandsLine
func parseCommandsLine(line string) (string, error) {
	// make it case-insensitive as a convenience feature
	upperLine := strings.ToUpper(strings.TrimSpace(line))

	for i, char := range upperLine {

		switch rover.Command(char) {
		case rover.CmdLeft, rover.CmdRight, rover.CmdMove:
			continue

		default:
			return "", fmt.Errorf("%w: character %v at position %d", ErrParseInvalidCommand, char, i)
		}
	}
	return upperLine, nil
}
