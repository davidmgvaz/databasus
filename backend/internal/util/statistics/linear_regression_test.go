package statistics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FitLine_WithPointsOnALine_ReturnsThatLine(t *testing.T) {
	line, isFitted := FitLine([]Point{{X: 0, Y: 1}, {X: 1, Y: 3}, {X: 2, Y: 5}, {X: 3, Y: 7}})

	require.True(t, isFitted)
	assert.InDelta(t, 2.0, line.Slope, 1e-9)
	assert.InDelta(t, 1.0, line.Intercept, 1e-9)
}

func Test_FitLine_WithScatteredPoints_ReturnsLeastSquaresLine(t *testing.T) {
	line, isFitted := FitLine([]Point{{X: 1, Y: 2}, {X: 2, Y: 3}, {X: 3, Y: 5}, {X: 4, Y: 4}})

	require.True(t, isFitted)
	assert.InDelta(t, 0.8, line.Slope, 1e-9)
	assert.InDelta(t, 1.5, line.Intercept, 1e-9)
}

func Test_FitLine_WithConstantY_ReturnsFlatLine(t *testing.T) {
	line, isFitted := FitLine([]Point{{X: 0, Y: 0.4}, {X: 1, Y: 0.4}, {X: 2, Y: 0.4}})

	require.True(t, isFitted)
	assert.InDelta(t, 0.0, line.Slope, 1e-9)
	assert.InDelta(t, 0.4, line.Intercept, 1e-9)
}

func Test_FitLine_WhenEveryPointSharesX_ReportsNoFit(t *testing.T) {
	_, isFitted := FitLine([]Point{{X: 5, Y: 1}, {X: 5, Y: 2}})

	assert.False(t, isFitted)
}

func Test_FitLine_WithOnePoint_ReportsNoFit(t *testing.T) {
	_, isFitted := FitLine([]Point{{X: 1, Y: 1}})

	assert.False(t, isFitted)
}
