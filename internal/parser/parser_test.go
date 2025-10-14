package parser

import (
	"mars/internal/rover"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestPlateau(t *testing.T, x, y int) *rover.Plateau {
	t.Helper()

	testPlateau, _ := rover.NewPlateau(x, y)
	return testPlateau
}

func createTestRoverPosition(t *testing.T, p *rover.Plateau, x, y int, dir rover.Direction) *rover.Position {
	t.Helper()

	newCoordinates := rover.NewCoordinates(x, y)
	newPosition, _ := rover.NewPosition(p, newCoordinates, dir)
	return newPosition
}

func createTestSingleRoverInstruction(t *testing.T, p *rover.Plateau, x, y int, dir rover.Direction, commands string) *rover.RoverInstruction {
	t.Helper()

	return &rover.RoverInstruction{
		InitialPosition: createTestRoverPosition(t, p, x, y, dir),
		Commands:        commands,
	}
}

func TestParseDirection(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		direction          string
		wantRoverDirection rover.Direction
		wantErr            error
	}{
		"ok - N": {
			direction:          "N",
			wantRoverDirection: rover.N,
			wantErr:            nil,
		},
		"ok - E": {
			direction:          "\tE",
			wantRoverDirection: rover.E,
			wantErr:            nil,
		},
		"ok - E - added whitespace": {
			direction:          "\tE   ",
			wantRoverDirection: rover.E,
			wantErr:            nil,
		},
		"ok - S": {
			direction:          "S",
			wantRoverDirection: rover.S,
			wantErr:            nil,
		},
		"ok - W": {
			direction:          "W",
			wantRoverDirection: rover.W,
			wantErr:            nil,
		},
		"ok - n": { // lower case
			direction:          "n",
			wantRoverDirection: rover.N,
			wantErr:            nil,
		},
		"err - unkonwn direction": {
			direction:          "XYZ",
			wantRoverDirection: rover.UnknownDirection,
			wantErr:            ErrParseInvalidDirection,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			roverDirection, err := parseDirection(tc.direction)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, roverDirection, tc.wantRoverDirection)
		})
	}
}

func TestParsePositionLine(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		line    string
		plateau *rover.Plateau
		wantErr error
	}{
		"ok - nominal": {
			line:    "5 5 N",
			plateau: createTestPlateau(t, 10, 10),
			wantErr: nil,
		},
		"err - ErrParsePositionFormat - > 3 elements": {
			line:    "5 5 N XYZ",
			plateau: createTestPlateau(t, 10, 10),
			wantErr: ErrParsePositionFormat,
		},
		"err - ErrParsePositionFormat - < 3 elements": {
			line:    "5 5",
			plateau: createTestPlateau(t, 10, 10),
			wantErr: ErrParsePositionFormat,
		},
		"err - ErrParsePositionX": {
			line:    "XYZ 5 N",
			plateau: createTestPlateau(t, 10, 10),
			wantErr: ErrParsePositionX,
		},
		"err - ErrParsePositionY": {
			line:    "5 XYZ N",
			plateau: createTestPlateau(t, 10, 10),
			wantErr: ErrParsePositionY,
		},
		"err - ErrParseInvalidDirection": {
			line:    "5 5 XYZ",
			plateau: createTestPlateau(t, 10, 10),
			// wrapped from parseDirection
			wantErr: ErrParseInvalidDirection,
		},
	}

	for name, tc := range testCases {

		t.Run(name, func(t *testing.T) {

			wantRoverPos := createTestRoverPosition(t, tc.plateau, 5, 5, rover.N)

			roverPosition, err := parsePositionLine(tc.line, tc.plateau)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, roverPosition, wantRoverPos)
		})
	}
}

func TestParseCommandsLines(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input   string
		wanted  string
		wantErr error
	}{
		"ok - nominal":                 {input: "M", wanted: "M", wantErr: nil},
		"ok - empty":                   {input: "", wanted: "", wantErr: nil},
		"ok - lower case":              {input: "m", wanted: "M", wantErr: nil},
		"err - ErrParseInvalidCommand": {input: "?", wanted: "", wantErr: ErrParseInvalidCommand},
	}

	for name, tc := range testCases {

		t.Run(name, func(t *testing.T) {

			parsingResult, err := parseCommandsLine(tc.input)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				assert.Len(t, parsingResult, 0)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, parsingResult, tc.wanted)
		})
	}
}

func TestParsePlateauLine(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input       string
		wantPlateau *rover.Plateau
		wantErr     error
	}{
		"ok - nominal": {input: "10 10", wantPlateau: createTestPlateau(t, 10, 10), wantErr: nil},
		"err - ErrParsePlateauFormat - < 1 element":  {input: "10", wantPlateau: nil, wantErr: ErrParsePlateauFormat},
		"err - ErrParsePlateauFormat - > 2 elements": {input: "10 10 10", wantPlateau: nil, wantErr: ErrParsePlateauFormat},
		"err - ErrParsePlateauX":                     {input: "XYZ 10", wantPlateau: nil, wantErr: ErrParsePlateauX},
		"err - ErrParsePlateauY":                     {input: "10 XYZ", wantPlateau: nil, wantErr: ErrParsePlateauY},
	}

	for name, tc := range testCases {

		t.Run(name, func(t *testing.T) {

			newPlateauLine, err := parsePlateauLine(tc.input)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, newPlateauLine, tc.wantPlateau)
		})
	}
}

func TestParse(t *testing.T) {
	testPlateau := createTestPlateau(t, 5, 5)

	testCases := map[string]struct {
		input            string
		wantPlateau      *rover.Plateau
		wantInstructions []rover.RoverInstruction
		wantErr          error
	}{"ok - nominal case for problem description": {
		input: `
5 5
1 2 N
LMLMLMLMM
3 3 E
MMRMMRMRRM`,
		wantPlateau: testPlateau,
		wantInstructions: []rover.RoverInstruction{
			*createTestSingleRoverInstruction(t, testPlateau, 1, 2, rover.N, "LMLMLMLMM"),
			*createTestSingleRoverInstruction(t, testPlateau, 3, 3, rover.E, "MMRMMRMRRM"),
		},
		wantErr: nil,
	},
		"error - invalid format (not enough lines)": {
			input: `
5 5
1 2 N`,
			wantPlateau:      nil,
			wantInstructions: nil,
			wantErr:          ErrParseInvalidFormat,
		},
		"error - invalid plateau line": {
			input: `
5 X
1 2 N
LMLM`,
			wantPlateau:      nil,
			wantInstructions: nil,
			wantErr:          ErrParsePlateauY,
		},
		"error - invalid position line (out of bounds)": {
			input: `
5 5
6 6 N
LMLM`,
			wantPlateau:      nil,
			wantInstructions: nil,
			// This error comes from rover.NewPosition, which is called by parsePositionLine
			wantErr: rover.ErrPositionOutOfBounds,
		},
		"error - invalid command line": {
			input: `
5 5
1 2 N
LMXLM`, // 'X' is an invalid command
			wantPlateau:      nil,
			wantInstructions: nil,
			wantErr:          ErrParseInvalidCommand,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			p := New() // Create a new parser instance

			gotPlateau, gotInstructions, err := p.Parse(tc.input)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantPlateau, gotPlateau)
			assert.Equal(t, tc.wantInstructions, gotInstructions)
		})
	}
}
