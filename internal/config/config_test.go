package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		args       []string
		wantConfig *Config
		wantErr    error
	}{
		"ok - with file": {
			args:       []string{"-file", "data.txt"},
			wantConfig: New(DefaultMinSizeX, DefaultMinSizeY, "data.txt", ModeCLI, DefaultServerAddr),
			wantErr:    nil,
		},
		"ok - webapi with custom port": {
			args:       []string{"-webapi", "-addr", ":5000"},
			wantConfig: New(DefaultMinSizeX, DefaultMinSizeY, "", ModeWebAPI, ":5000"),
			wantErr:    nil,
		},
		"ok - no flags (defaults)": {
			args:       []string{},
			wantConfig: New(DefaultMinSizeX, DefaultMinSizeY, "", ModeCLI, DefaultServerAddr),
			wantErr:    nil,
		},
		"ok - webapi default addr": {
			args:       []string{"-webapi"},
			wantConfig: New(DefaultMinSizeX, DefaultMinSizeY, "", ModeWebAPI, DefaultServerAddr),
			wantErr:    nil,
		},
		"ok - with minX and minY flags": {
			args:       []string{"-min-size-x", "5", "-min-size-y", "6"},
			wantConfig: New(5, 6, "", ModeCLI, DefaultServerAddr),
			wantErr:    nil,
		},
		"err - negative dimensions": {
			args:    []string{"-min-size-x", "-1", "-min-size-y", "5"},
			wantErr: ErrParserPlateauDimensions,
		},
		"err - incompatible flags": {
			args:    []string{"-webapi", "-file", "data.txt"},
			wantErr: ErrParserFlagsIncompatible,
		},
		"err - invalid minX and minY dimensions": {
			args:    []string{"-min-size-x", "0", "-min-size-y", "0"},
			wantErr: ErrParserPlateauDimensions,
		},
		"err - parsing error": {
			args:    []string{"-min-size-x", "invalid"},
			wantErr: ErrParserInvalidValue,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gotConfig, err := ParseFlags(tc.args)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, gotConfig, tc.wantConfig)
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	minX := 20
	minY := 25
	filepath := "something.txt"
	opMode := ModeWebAPI
	addr := ":7171"

	newConfig := New(minX, minY, filepath, opMode, addr)

	require.NotNil(t, newConfig)
	assert.Equal(t, minX, newConfig.MinPlateauX)
	assert.Equal(t, minY, newConfig.MinPlateauY)
	assert.Equal(t, filepath, newConfig.FilePath)
	assert.Equal(t, opMode, newConfig.OpMode)
	assert.Equal(t, addr, newConfig.SrvAddr)
}

func TestDefault(t *testing.T) {
	t.Parallel()

	const (
		defaultFilePath = ""
	)

	cfgDefault := Default()

	require.NotNil(t, cfgDefault)
	assert.Equal(t, DefaultMinSizeX, cfgDefault.MinPlateauX)
	assert.Equal(t, DefaultMinSizeY, cfgDefault.MinPlateauY)
	assert.Equal(t, defaultFilePath, cfgDefault.FilePath)
	assert.Equal(t, ModeCLI, cfgDefault.OpMode)
	assert.Equal(t, DefaultServerAddr, cfgDefault.SrvAddr)
}

func TestValidate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config  *Config
		wantErr error
	}{
		"ok - default": {
			config:  Default(),
			wantErr: nil,
		},
		"err - ErrParserNilConfig": {
			config:  nil,
			wantErr: ErrParserNilConfig,
		},
		"err - ErrParserPlateauDimensions": {
			config:  New(0, 0, "", ModeCLI, ""),
			wantErr: ErrParserPlateauDimensions,
		},
		"err - ErrParserModeUnknown": {
			config:  New(5, 5, "", ModeUnknown, ""),
			wantErr: ErrParserModeUnknown,
		},
		"err - ErrParserServerAddr": {
			config:  New(5, 5, "", ModeWebAPI, ""),
			wantErr: ErrParserServerAddr,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			err := tc.config.Validate()

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, tc.wantErr)
		})
	}
}
