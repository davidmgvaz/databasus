package statistics

type Point struct {
	X float64
	Y float64
}

type Line struct {
	Intercept float64
	Slope     float64
}

// Ordinary least squares. The fit is undefined when every point shares one X, so it reports
// false instead of dividing by zero.
func FitLine(points []Point) (Line, bool) {
	if len(points) < 2 {
		return Line{}, false
	}

	var sumX, sumY float64
	for _, point := range points {
		sumX += point.X
		sumY += point.Y
	}

	meanX := sumX / float64(len(points))
	meanY := sumY / float64(len(points))

	var covariance, varianceX float64
	for _, point := range points {
		deviationX := point.X - meanX
		covariance += deviationX * (point.Y - meanY)
		varianceX += deviationX * deviationX
	}

	if varianceX == 0 {
		return Line{}, false
	}

	slope := covariance / varianceX

	return Line{Intercept: meanY - slope*meanX, Slope: slope}, true
}
