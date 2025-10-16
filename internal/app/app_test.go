package app

import (
	"bytes"
	"errors"
	"io"

	"mars/internal/config"
	"mars/internal/rover"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockParser will implement the Parser interface
type MockParser struct {
	mock.Mock
}

func (m *MockParser) Parse(input string, cfg *config.Config) (*rover.Plateau, []rover.RoverInstruction, error) {
	args := m.Called(input)

	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*rover.Plateau), args.Get(1).([]rover.RoverInstruction), args.Error(2)
}

// MockMissionControlFactory will implement the rover.MissionControlFactory interface
type MockMissionControlFactory struct {
	mock.Mock
}

func (m *MockMissionControlFactory) Create(plateau *rover.Plateau) (*rover.MissionControl, error) {
	args := m.Called(plateau)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rover.MissionControl), args.Error(1)
}

// errReader is a custom reader that always returns an error
type errReader struct{}

func (e errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestApp_Run(t *testing.T) {
	t.Parallel()

	cfg := config.Default()

	testCases := map[string]struct {
		inputData   string
		inputReader io.Reader
		setupMocks  func(*MockParser, *MockMissionControlFactory)
		wantOutput  string
		wantErr     error
	}{
		"ok - complete mission": {
			inputData: "5 5\n1 2 N\nLMLMLMLMM\n3 3 E\nMMRMMRMRRM",

			setupMocks: func(mp *MockParser, mmcf *MockMissionControlFactory) {
				plateau, _ := rover.NewPlateau(5, 5, cfg.MinPlateauX, cfg.MinPlateauY)
				pos1, _ := rover.NewPosition(plateau, rover.NewCoordinates(1, 2), rover.N)
				pos2, _ := rover.NewPosition(plateau, rover.NewCoordinates(3, 3), rover.E)

				instructions := []rover.RoverInstruction{
					{InitialPosition: pos1, Commands: "LMLMLMLMM"},
					{InitialPosition: pos2, Commands: "MMRMMRMRRM"},
				}

				mp.On("Parse", "5 5\n1 2 N\nLMLMLMLMM\n3 3 E\nMMRMMRMRRM").Return(plateau, instructions, nil)

				mc, _ := rover.NewMissionControl(plateau)
				mmcf.On("Create", plateau).Return(mc, nil)
			},
			wantOutput: "1 3 N\n5 1 E\n",
			wantErr:    nil,
		},
		"ok - single rover": {
			inputData: "5 5\n1 2 N\nLMLMLMLMM",

			setupMocks: func(mp *MockParser, mmcf *MockMissionControlFactory) {
				plateau, _ := rover.NewPlateau(5, 5, cfg.MinPlateauX, cfg.MinPlateauY)
				pos1, _ := rover.NewPosition(plateau, rover.NewCoordinates(1, 2), rover.N)

				instructions := []rover.RoverInstruction{
					{InitialPosition: pos1, Commands: "LMLMLMLMM"},
				}

				mp.On("Parse", "5 5\n1 2 N\nLMLMLMLMM").Return(plateau, instructions, nil)

				mc, _ := rover.NewMissionControl(plateau)
				mmcf.On("Create", plateau).Return(mc, nil)
			},
			wantOutput: "1 3 N\n",
			wantErr:    nil,
		},
		"err - reading input fails": {
			inputReader: errReader{},

			setupMocks: func(mp *MockParser, mmcf *MockMissionControlFactory) {
				// No mock needed - error happens before parsing
			},
			wantErr: ErrAppInput,
		},
		"err - parsing fails": {
			inputData: "invalid input",

			setupMocks: func(mp *MockParser, mmcf *MockMissionControlFactory) {
				mp.On("Parse", "invalid input").Return(nil, nil, errors.New("parse error"))
			},
			wantErr: ErrAppParsing,
		},
		"err - mission control creation fails": {
			inputData: "5 5\n1 2 N\nLMLMLMLMM",

			setupMocks: func(mp *MockParser, mmcf *MockMissionControlFactory) {
				plateau, _ := rover.NewPlateau(5, 5, cfg.MinPlateauX, cfg.MinPlateauY)
				pos1, _ := rover.NewPosition(plateau, rover.NewCoordinates(1, 2), rover.N)
				instructions := []rover.RoverInstruction{
					{InitialPosition: pos1, Commands: "LMLMLMLMM"},
				}

				mp.On("Parse", mock.Anything).Return(plateau, instructions, nil)
				mmcf.On("Create", plateau).Return(nil, errors.New("mc creation failed"))
			},
			wantErr: ErrAppCreatingMC,
		},
		"err - execution fails with collision": {
			inputData: "5 5\n1 2 N\nLMLMLMLMM\n1 3 N\nM",

			setupMocks: func(mp *MockParser, mmcf *MockMissionControlFactory) {
				plateau, _ := rover.NewPlateau(5, 5, cfg.MinPlateauX, cfg.MinPlateauY)
				pos1, _ := rover.NewPosition(plateau, rover.NewCoordinates(1, 2), rover.N)
				pos2, _ := rover.NewPosition(plateau, rover.NewCoordinates(1, 3), rover.N)
				instructions := []rover.RoverInstruction{
					{InitialPosition: pos1, Commands: "LMLMLMLMM"},
					{InitialPosition: pos2, Commands: "M"},
				}

				mp.On("Parse", mock.Anything).Return(plateau, instructions, nil)

				mc, _ := rover.NewMissionControl(plateau)
				mmcf.On("Create", plateau).Return(mc, nil)
			},
			wantErr: ErrAppExecMission,
		},
	}

	for name, tc := range testCases {

		t.Run(name, func(t *testing.T) {
			// Setup mocks
			mockParser := new(MockParser)
			mockMCFactory := new(MockMissionControlFactory)

			tc.setupMocks(mockParser, mockMCFactory)

			// Setup input/output
			var input io.Reader
			if tc.inputReader != nil {
				input = tc.inputReader
			} else {
				input = strings.NewReader(tc.inputData)
			}
			output := &bytes.Buffer{}

			// Create app and run
			app := NewApp(mockParser, mockMCFactory, input, output, cfg)
			err := app.Run()

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Contains(t, output.String(), tc.wantOutput)

			// Verify mocks were called as expected
			mockParser.AssertExpectations(t)
			mockMCFactory.AssertExpectations(t)
		})
	}
}

func TestNewApp(t *testing.T) {
	t.Parallel()

	cfg := config.Default()

	mockParser := new(MockParser)
	mockMCFactory := new(MockMissionControlFactory)

	input := strings.NewReader("test input")
	output := &bytes.Buffer{}

	app := NewApp(mockParser, mockMCFactory, input, output, cfg)

	require.NotNil(t, app)
	assert.Equal(t, mockParser, app.parser)
	assert.Equal(t, mockMCFactory, app.mcf)
	assert.Equal(t, input, app.input)
	assert.Equal(t, output, app.output)
}
