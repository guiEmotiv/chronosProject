package plotters

import (
	"../core"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func PlotTravel(recorder *core.WorkerRecorder, name string) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.X.Label.Text = "X position [m]"
	p.Y.Label.Text = "Y position [m]"

	pts := make(plotter.XYs, len(recorder.XPositions))

	for i := 0; i < len(recorder.XPositions); i++ {
		pts[i].X = recorder.XPositions[i]
		pts[i].Y = recorder.YPositions[i]
	}

	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err.Error())
	}
	scatter.Shape = plotutil.Shape(5)
	err = plotutil.AddScatters(p, scatter)
	if err != nil {
		panic(err.Error())
	}

	if err := p.Save(4*vg.Inch, 4*vg.Inch, name); err != nil {
		panic(err)
	}
}
