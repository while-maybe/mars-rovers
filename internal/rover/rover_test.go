package rover

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCoordinates(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		x               int
		y               int
		wantCoordinates Coordinates
	}{
		"zero coordinates": {
			x: 0, y: 0,
			wantCoordinates: Coordinates{x: 0, y: 0},
		},
		"negative coordinates": {
			x: -8, y: -3,
			wantCoordinates: Coordinates{x: -8, y: -3},
		},
		"positive coordinates": {
			x: 5, y: 5,
			wantCoordinates: Coordinates{x: 5, y: 5},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {

			newCoords := NewCoordinates(testCase.x, testCase.y)

			assert.Equal(t, newCoords, testCase.wantCoordinates)
		})
	}
}

func TestNewPlateau(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		maxX        int
		maxY        int
		wantPlateau *Plateau
		wantErr     error
	}{
		"ok - dimensions more than min X and min Y": {
			maxX: 10, maxY: 10,
			wantPlateau: &Plateau{maxX: 10, maxY: 10},
			wantErr:     nil,
		},
		"err - less than min x dimension": {
			maxX: 1, maxY: 10,
			wantPlateau: nil,
			wantErr:     ErrPlateauTooSmall,
		},
		"err - less than min y dimension": {
			maxX: 10, maxY: 1,
			wantPlateau: nil,
			wantErr:     ErrPlateauTooSmall,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {

			newPlateau, err := NewPlateau(testCase.maxX, testCase.maxY)
			assert.Equal(t, newPlateau, testCase.wantPlateau)

			if testCase.wantErr != nil {
				require.ErrorIs(t, err, testCase.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestNewPosition(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		plateau      *Plateau
		coordinates  Coordinates
		direction    Direction
		wantPosition *Position
		wantErr      error
	}{
		"ok - nominal": {
			plateau:      &Plateau{maxX: 10, maxY: 10},
			coordinates:  NewCoordinates(5, 5),
			direction:    N,
			wantPosition: &Position{coordinates: Coordinates{x: 5, y: 5}, direction: N},
			wantErr:      nil,
		},
		"err - ErrPositionOutOfBounds": {
			plateau:      &Plateau{maxX: 4, maxY: 4},
			coordinates:  NewCoordinates(5, 5),
			direction:    N,
			wantPosition: nil,
			wantErr:      ErrPositionOutOfBounds,
		},
		"err - ErrDirectionUnknown": {
			plateau:      &Plateau{maxX: 10, maxY: 10},
			coordinates:  NewCoordinates(5, 5),
			direction:    UnknownDirection,
			wantPosition: nil,
			wantErr:      ErrDirectionUnknown,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			newPosition, err := NewPosition(tc.plateau, tc.coordinates, tc.direction)
			assert.Equal(t, newPosition, tc.wantPosition)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestNewRover(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		id        int
		pos       *Position
		wantRover *Rover
		wantErr   error
	}{
		"ok - nominal": {
			id:  1,
			pos: &Position{coordinates: Coordinates{x: 5, y: 5}, direction: N},
			wantRover: &Rover{
				id:       1,
				position: &Position{coordinates: Coordinates{x: 5, y: 5}, direction: N},
			},
			wantErr: nil,
		},
		"err - ErrRoverPositionIsNil": {
			id:        1,
			pos:       nil,
			wantRover: nil,
			wantErr:   ErrRoverPositionIsNil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			newRover, err := NewRover(tc.id, tc.pos)
			assert.Equal(t, newRover, tc.wantRover)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestNewMissionControl(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		plateau *Plateau
		wantMC  *MissionControl
		wantErr error
	}{
		"ok - nominal": {
			plateau: &Plateau{maxX: 10, maxY: 10},
			wantMC: &MissionControl{
				plateau:         &Plateau{maxX: 10, maxY: 10},
				occupiedSquares: map[Coordinates]int{},
			},
			wantErr: nil,
		},
		"err - ErrPlateauIsNil": {
			plateau: nil,
			wantMC:  nil,
			wantErr: ErrPlateauIsNil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			newMissionControl, err := NewMissionControl(tc.plateau)
			assert.Equal(t, newMissionControl, tc.wantMC)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestRunRover(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		mc         *MissionControl
		rover      *Rover
		commands   string
		wantString string
		wantErr    error
	}{
		"ok - nominal": {
			mc: &MissionControl{
				plateau:         &Plateau{maxX: 10, maxY: 10},
				occupiedSquares: map[Coordinates]int{},
			},
			rover: &Rover{
				id:       1,
				position: &Position{coordinates: Coordinates{x: 5, y: 5}, direction: N},
			},
			commands:   "MMLLLLRRRR",
			wantString: "5 7 N",
			wantErr:    nil,
		},
		"ok - empty command string": {
			mc: &MissionControl{
				plateau:         &Plateau{maxX: 10, maxY: 10},
				occupiedSquares: map[Coordinates]int{},
			},
			rover: &Rover{
				id:       1,
				position: &Position{coordinates: Coordinates{x: 5, y: 5}, direction: N},
			},
			commands:   "",
			wantString: "5 5 N",
			wantErr:    nil,
		},
		"ok - invalid command characters ignored": {
			mc: &MissionControl{
				plateau:         &Plateau{maxX: 10, maxY: 10},
				occupiedSquares: map[Coordinates]int{},
			},
			rover: &Rover{
				id:       1,
				position: &Position{coordinates: Coordinates{x: 5, y: 5}, direction: N},
			},
			commands:   "MXYZM", // X, Y, Z should be ignored
			wantString: "5 7 N",
			wantErr:    nil,
		},
		"ok - multiple rovers on plateau": {
			mc: &MissionControl{
				plateau: &Plateau{maxX: 10, maxY: 10},
				occupiedSquares: map[Coordinates]int{
					{3, 3}: 1,
					{7, 7}: 2,
				},
			},
			rover: &Rover{
				id:       3,
				position: &Position{coordinates: Coordinates{x: 5, y: 5}, direction: N},
			},
			commands:   "M",
			wantString: "5 6 N",
			wantErr:    nil,
		},
		"err - ErrCollisionDetected - placing new rover": {
			mc: &MissionControl{
				plateau:         &Plateau{maxX: 10, maxY: 10},
				occupiedSquares: map[Coordinates]int{{5, 7}: 0},
			},
			rover: &Rover{
				id:       1,
				position: &Position{coordinates: Coordinates{x: 5, y: 7}, direction: N},
			},
			commands:   "MM",
			wantString: "",
			wantErr:    ErrRoverCollision,
		},
		// this is not a mistake as the rover stops if it tries to move to an occupied location
		"ok - ErrCollisionDetected - rover en route": {
			mc: &MissionControl{
				plateau:         &Plateau{maxX: 10, maxY: 10},
				occupiedSquares: map[Coordinates]int{{5, 7}: 0},
			},
			rover: &Rover{
				id:       1,
				position: &Position{coordinates: Coordinates{x: 5, y: 5}, direction: N},
			},
			commands:   "MMLL",
			wantString: "5 6 S",
			wantErr:    nil,
		},
		"ok - stops before out of bounds": {
			mc: &MissionControl{
				plateau:         &Plateau{maxX: 10, maxY: 10},
				occupiedSquares: map[Coordinates]int{{5, 7}: 0},
			},
			rover: &Rover{
				id:       1,
				position: &Position{coordinates: Coordinates{x: 10, y: 10}, direction: N},
			},
			commands:   "MMR",
			wantString: "10 10 E",
			wantErr:    nil,
		},
		"ok - turn in all directions": {
			mc: &MissionControl{
				plateau:         &Plateau{maxX: 10, maxY: 10},
				occupiedSquares: map[Coordinates]int{{5, 7}: 0},
			},
			rover: &Rover{
				id:       1,
				position: &Position{coordinates: Coordinates{x: 5, y: 5}, direction: N},
			},
			commands:   "MRMRMRM",
			wantString: "5 5 W",
			wantErr:    nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			result, err := tc.mc.RunRover(tc.rover, tc.commands)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.Equal(t, result, tc.wantString)

			require.NoError(t, err)
		})
	}
}

func TestMove(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		rover        *Rover
		wantPosition *Position
	}{
		"ok - N": {
			rover: &Rover{
				id:       1,
				position: &Position{coordinates: Coordinates{x: 10, y: 10}, direction: N},
			},
			wantPosition: &Position{coordinates: Coordinates{x: 10, y: 11}, direction: N},
		},
		"ok - E": {
			rover: &Rover{
				id:       1,
				position: &Position{coordinates: Coordinates{x: 10, y: 10}, direction: E},
			},
			wantPosition: &Position{coordinates: Coordinates{x: 11, y: 10}, direction: E},
		},
		"ok - S": {
			rover: &Rover{
				id:       1,
				position: &Position{coordinates: Coordinates{x: 10, y: 10}, direction: S},
			},
			wantPosition: &Position{coordinates: Coordinates{x: 10, y: 9}, direction: S},
		},
		"ok - W": {
			rover: &Rover{
				id:       1,
				position: &Position{coordinates: Coordinates{x: 10, y: 10}, direction: W},
			},
			wantPosition: &Position{coordinates: Coordinates{x: 9, y: 10}, direction: W},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.rover.move(), *tc.wantPosition)
		})
	}
}

// turnLeft and turnRight are implicitly tested

func createTestPlateau(t *testing.T, x, y int) *Plateau {
	t.Helper()

	testPlateau, _ := NewPlateau(x, y)
	return testPlateau
}

func createTestRoverPosition(t *testing.T, p *Plateau, x, y int, dir Direction) *Position {
	t.Helper()

	newCoordinates := NewCoordinates(x, y)
	newPosition, _ := NewPosition(p, newCoordinates, dir)
	return newPosition
}

func createTestSingleRoverInstruction(t *testing.T, p *Plateau, x, y int, dir Direction, commands string) *RoverInstruction {
	t.Helper()

	return &RoverInstruction{
		InitialPosition: createTestRoverPosition(t, p, x, y, dir),
		Commands:        commands,
	}
}

func TestMissionControlExecute(t *testing.T) {
	t.Parallel()
	testPlateau := createTestPlateau(t, 5, 5)

	testCases := map[string]struct {
		mc         *MissionControl
		mcInput    *MissionControlInput
		wantOutput []string
		wantErr    error
	}{
		"ok - multiple rovers": {
			mc: &MissionControl{
				plateau:         testPlateau,
				occupiedSquares: map[Coordinates]int{},
			},
			mcInput: &MissionControlInput{
				Instructions: []RoverInstruction{
					*createTestSingleRoverInstruction(t, testPlateau, 1, 2, N, "LMLMLMLMM"),
					*createTestSingleRoverInstruction(t, testPlateau, 3, 3, E, "MMRMMRMRRM"),
				},
			},
			wantOutput: []string{"1 3 N", "5 1 E"},
			wantErr:    nil,
		},
		"err - ErrRoverCreating": {
			mc: &MissionControl{
				plateau:         testPlateau,
				occupiedSquares: map[Coordinates]int{},
			},
			mcInput: &MissionControlInput{
				Instructions: []RoverInstruction{
					*createTestSingleRoverInstruction(t, testPlateau, 10, 2, N, "LMLMLMLMM"),
					*createTestSingleRoverInstruction(t, testPlateau, 30, 3, E, "MMRMMRMRRM"),
				},
			},
			wantOutput: []string{},
			wantErr:    ErrRoverCreating,
		},
		"err - ErrRoverInstructions": {
			mc: &MissionControl{
				plateau:         testPlateau,
				occupiedSquares: map[Coordinates]int{},
			},
			mcInput: &MissionControlInput{
				Instructions: []RoverInstruction{
					*createTestSingleRoverInstruction(t, testPlateau, 1, 2, N, "LMLMLMLMM"),
					// force a collision with another existing rover
					*createTestSingleRoverInstruction(t, testPlateau, 1, 3, N, "MMRMMRMRRM"),
				},
			},
			wantOutput: []string{},
			wantErr:    ErrRoverInstructions,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var output []string

			output, err := tc.mc.Execute(tc.mcInput)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				assert.Len(t, output, 0)
				return
			}

			assert.Equal(t, output, tc.wantOutput)
			require.NoError(t, err)
		})
	}
}
