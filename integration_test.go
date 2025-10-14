package main_test

import (
	"testing"

	"mars/internal/parser"
	"mars/internal/rover"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// this is the test proposed in the exercise
func TestMarsRoverRequirements(t *testing.T) {
	input := `5 5
1 2 N
LMLMLMLMM
3 3 E
MMRMMRMRRM`

	// Parse input
	p := parser.New()
	plateau, instructions, err := p.Parse(input)
	require.NoError(t, err)

	// Create mission control
	factory := rover.NewMissionControlFactory()
	mc, err := factory.Create(plateau)
	require.NoError(t, err)

	// Execute mission
	missionInput := &rover.MissionControlInput{
		Instructions: instructions,
	}
	outputs, err := mc.Execute(missionInput)
	require.NoError(t, err)

	// Verify expected outputs
	expectedOutputs := []string{
		"1 3 N",
		"5 1 E",
	}

	assert.Equal(t, expectedOutputs, outputs)
}

func TestRoverCollisionDetection(t *testing.T) {
	// Rover 1 occupies (3, 3)
	// Rover 2 tries to move into (3, 3) but is blocked
	input := `5 5
3 3 N
L
3 1 N
MM`

	p := parser.New()
	plateau, instructions, err := p.Parse(input)
	require.NoError(t, err)

	factory := rover.NewMissionControlFactory()
	mc, err := factory.Create(plateau)
	require.NoError(t, err)

	missionInput := &rover.MissionControlInput{
		Instructions: instructions,
	}
	outputs, err := mc.Execute(missionInput)
	require.NoError(t, err)

	// Rover 1: (3,3) N, turns left to face W, stays at (3,3) W
	assert.Equal(t, "3 3 W", outputs[0])

	// Rover 2: (3,1) N, moves M→(3,2), tries M→(3,3) BLOCKED, stays at (3,2) N
	assert.Equal(t, "3 2 N", outputs[1])
}

func TestMultipleRoversWithCollisionAvoidance(t *testing.T) {
	// Rover 1 moves to (5, 7) and stays there
	// Rover 2 tries to move through (5, 7) but should be blocked
	input := `10 10
5 5 N
MM
5 9 S
MMM`

	p := parser.New()
	plateau, instructions, err := p.Parse(input)
	require.NoError(t, err)

	factory := rover.NewMissionControlFactory()
	mc, err := factory.Create(plateau)
	require.NoError(t, err)

	missionInput := &rover.MissionControlInput{
		Instructions: instructions,
	}
	outputs, err := mc.Execute(missionInput)
	require.NoError(t, err)

	// Rover 1: (5,5) N → moves MM → ends at (5,7) N
	assert.Equal(t, "5 7 N", outputs[0])

	// Rover 2: (5,9) S → tries MMM south
	// (5,9) → (5,8) ✓ → tries (5,7) ✗ BLOCKED by Rover 1 → stays at (5,8)
	assert.Equal(t, "5 8 S", outputs[1])
}

func TestRoverBoundaryRespect(t *testing.T) {
	input := `5 5
0 0 S
MMM
5 5 N
MMM`

	p := parser.New()
	plateau, instructions, err := p.Parse(input)
	require.NoError(t, err)

	factory := rover.NewMissionControlFactory()
	mc, err := factory.Create(plateau)
	require.NoError(t, err)

	missionInput := &rover.MissionControlInput{
		Instructions: instructions,
	}
	outputs, err := mc.Execute(missionInput)
	require.NoError(t, err)

	// Rover 1: (0,0) S, can't go south (already at boundary)
	assert.Equal(t, "0 0 S", outputs[0])

	// Rover 2: (5,5) N, can't go north (already at boundary)
	assert.Equal(t, "5 5 N", outputs[1])
}

func TestComplexRoverNavigationPattern(t *testing.T) {
	// Test rover navigating in a specific pattern
	input := `10 10
2 2 N
MMMRMMMLMM`

	p := parser.New()
	plateau, instructions, err := p.Parse(input)
	require.NoError(t, err)

	factory := rover.NewMissionControlFactory()
	mc, err := factory.Create(plateau)
	require.NoError(t, err)

	missionInput := &rover.MissionControlInput{
		Instructions: instructions,
	}
	outputs, err := mc.Execute(missionInput)
	require.NoError(t, err)

	// Trace:
	// Start: (2,2) N
	// MMM:   (2,5) N    [3 north]
	// R:     (2,5) E    [turn right]
	// MMM:   (5,5) E    [3 east]
	// L:     (5,5) N    [turn left]
	// MM:    (5,7) N    [2 north]
	assert.Equal(t, "5 7 N", outputs[0])
}
