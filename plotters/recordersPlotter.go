package plotters

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"../core"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"fmt"
)

func PlotSimulation(palette color.Palette, recorders ...*core.WorkerRecorder) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.X.Label.Text = "X position [m]"
	p.Y.Label.Text = "Y position [m]"


	colors := palette

	for i, r := range recorders {

		pts := make(plotter.XYs, len(r.XPositions))
		for i := 0; i < len(r.XPositions); i++ {
			pts[i].X = r.XPositions[i]
			pts[i].Y = r.YPositions[i]
		}

		scatter, err := plotter.NewScatter(pts)
		if err != nil {
			panic(err.Error())
		}

		scatter.Shape = draw.CircleGlyph{}
		scatter.Radius = vg.Points(3)
		scatter.Color = colors[i*15]

		p.Add(scatter)



		p.Legend.Add(fmt.Sprintf("Worker %d", i), scatter)
		if err != nil {
			panic(err.Error())
		}
	}

	if err := p.Save(10*vg.Inch, 10*vg.Inch, "simulation.png"); err != nil {
		panic(err)
	}
}
