package graph

import (
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var (
	XLabel string = "X"
	YLabel string = "Y"

	OutputFilepath string = "points.png"
)

func PointsFromValues(allValues [][]string) (plotter.XYs, error) {
	var points plotter.XYs
	for i, value := range allValues {
		xVal := float64(i)

		yVal, err := strconv.ParseFloat(value[1], 64)
		if err != nil {
			return points, err
		}

		points = append(points, plotter.XY{
			X: xVal,
			Y: yVal,
		})
	}
	return points, nil
}

func GraphPoints(points plotter.XYs) error {
	p := plot.New()

	p.Title.Text = "Price History"
	p.X.Label.Text = XLabel
	p.Y.Label.Text = YLabel

	err := plotutil.AddLinePoints(p, "First", points)
	if err != nil {
		return err
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, OutputFilepath); err != nil {
		return err
	}
	return nil
}
