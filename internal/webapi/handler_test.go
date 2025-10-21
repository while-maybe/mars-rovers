package webapi

import (
	"errors"
	"mars/internal/app"
	"mars/internal/config"
	"mars/internal/rover"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockParser struct {
	mock.Mock
}

func (m *MockParser) Parse(input string, cfg *config.Config) (*rover.Plateau, []rover.RoverInstruction, error) {
	args := m.Called(input, cfg)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}

	return args.Get(0).(*rover.Plateau), args.Get(1).([]rover.RoverInstruction), args.Error(2)
}

type MockMCFactory struct {
	mock.Mock
}

func (m *MockMCFactory) Create(plateau *rover.Plateau) (*rover.MissionControl, error) {
	args := m.Called(plateau)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*rover.MissionControl), args.Error(1)
}

func TestHandleMission(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		httpMethod       string
		requestBody      string
		setupMocks       func(*MockParser, *MockMCFactory)
		wantStatusCode   int
		wantBodyContains string
	}{
		"ok - nominal": {
			httpMethod:  http.MethodPost,
			requestBody: "5 5\n1 2 N\nLMLMLMLMM",
			setupMocks: func(mp *MockParser, mcf *MockMCFactory) {
				plateau, _ := rover.NewPlateau(5, 5, 2, 2)
				pos, _ := rover.NewPosition(plateau, rover.NewCoordinates(1, 2), rover.N)
				instructions := []rover.RoverInstruction{{InitialPosition: pos, Commands: "LMLMLMLMM"}}

				mp.On("Parse", "5 5\n1 2 N\nLMLMLMLMM", mock.AnythingOfType("*config.Config")).Return(plateau, instructions, nil)

				mc, _ := rover.NewMissionControl(plateau)
				mcf.On("Create", plateau).Return(mc, nil)
			},
			wantStatusCode:   http.StatusOK,
			wantBodyContains: "1 3 N",
		},
		"err - parser fails": {
			httpMethod:  http.MethodPost,
			requestBody: "make parsing fail",
			setupMocks: func(mp *MockParser, mcf *MockMCFactory) {
				mp.On("Parse", "make parsing fail", mock.AnythingOfType("*config.Config")).Return(nil, nil, errors.New("parsing failed"))
			},
			wantStatusCode:   http.StatusBadRequest,
			wantBodyContains: app.ErrAppParsing.Error(),
		},
		"err - req body too large": {
			httpMethod:  http.MethodPost,
			requestBody: strings.Repeat("?", maxRequestSize+1),
			setupMocks: func(mp *MockParser, mm *MockMCFactory) {
				// No mocks needed
			},
			wantStatusCode:   http.StatusRequestEntityTooLarge,
			wantBodyContains: "Request body is too large",
		},
		"err - unhandled internal error": {
			httpMethod:  http.MethodPost,
			requestBody: "5 5\n1 2 N\nLMLMLMLMM",
			setupMocks: func(mp *MockParser, mcf *MockMCFactory) {
				plateau, _ := rover.NewPlateau(5, 5, 2, 2)
				pos, _ := rover.NewPosition(plateau, rover.NewCoordinates(1, 2), rover.N)
				instructions := []rover.RoverInstruction{{InitialPosition: pos, Commands: "LMLMLMLMM"}}

				mp.On("Parse", mock.Anything, mock.Anything).Return(plateau, instructions, nil)

				mcf.On("Create", plateau).Return(nil, errors.New("a different type of error goes here"))
			},
			wantStatusCode:   http.StatusInternalServerError,
			wantBodyContains: "An internal server error occurred",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			// dependencies
			mockParser := new(MockParser)
			mockMCFactory := new(MockMCFactory)
			tc.setupMocks(mockParser, mockMCFactory)

			// new server instance
			server := NewServer(config.Default(), mockParser, mockMCFactory)

			// fake request
			req := httptest.NewRequest(tc.httpMethod, "/mcontrol", strings.NewReader(tc.requestBody))
			// capture response
			rcap := httptest.NewRecorder()

			// call handler directly
			server.handleMission(rcap, req)

			assert.Equal(t, tc.wantStatusCode, rcap.Code)
			assert.Contains(t, rcap.Body.String(), tc.wantBodyContains)

			mockParser.AssertExpectations(t)
			mockMCFactory.AssertExpectations(t)
		})
	}
}
func TestServer_RoutingOnly(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		method         string
		endpoint       string
		wantStatusCode int
	}{
		"ok method, ok path": {
			method:         http.MethodPost,
			endpoint:       "/mcontrol",
			wantStatusCode: http.StatusOK,
		},
		"ok method, err path": {
			method:         http.MethodPost,
			endpoint:       "/doesnotexist",
			wantStatusCode: http.StatusNotFound,
		},
		"err method, ok path": {
			method:         http.MethodGet,
			endpoint:       "/mcontrol",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		"err method, err path": {
			method:         http.MethodPatch,
			endpoint:       "/doesnotexist",
			wantStatusCode: http.StatusNotFound,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockParser := new(MockParser)
			mockMCFactory := new(MockMCFactory)

			if tc.wantStatusCode == http.StatusOK {
				plateau, _ := rover.NewPlateau(5, 5, 2, 2)
				mockParser.On("Parse", mock.Anything, mock.Anything).Return(plateau, []rover.RoverInstruction{}, nil)

				mc, _ := rover.NewMissionControl(plateau)
				mockMCFactory.On("Create", plateau).Return(mc, nil)
			}

			server := NewServer(config.Default(), mockParser, mockMCFactory)

			router := server.Handler()

			// fake request
			req := httptest.NewRequest(tc.method, tc.endpoint, nil)
			// response capture (recorder)
			rcap := httptest.NewRecorder()

			router.ServeHTTP(rcap, req)

			assert.Equal(t, tc.wantStatusCode, rcap.Code)

			mockParser.AssertExpectations(t)
			mockMCFactory.AssertExpectations(t)
		})
	}
}
