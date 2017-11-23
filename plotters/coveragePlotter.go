package plotters

import (
	"../core"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"fmt"
	"gonum.org/v1/plot/vg"
	"image/color"
)

func PLotCoverage(recorder *core.CoverageRecorder) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	//pts1 := make(plotter.XYs, len(recorder.GlobalCoverage))
	//
	//for i := 0; i < len(recorder.GlobalCoverage); i++ {
	//	pts1[i].X = float64(i)
	//	pts1[i].Y = recorder.GlobalCoverage[i]
	//}
	//
	//absoluteCoverage, err := plotter.NewLine(pts1)
	//if err != nil {
	//	panic(err.Error())
	//}
	//absoluteCoverage.Color = color.RGBA{95, 102, 234, 255}
	//
	//p.Add(absoluteCoverage)

	pts2 := make(plotter.XYs, len(recorder.RelativeCoverage))

	for i := 0; i < len(recorder.RelativeCoverage); i++ {
		pts2[i].X = float64(i)
		pts2[i].Y = recorder.RelativeCoverage[i]
	}

	relativeCoverage, err := plotter.NewLine(pts2)
	if err != nil {
		panic(err.Error())
	}
	relativeCoverage.Color = color.RGBA{95, 102, 234, 255}

	p.Add(relativeCoverage)



	p.Legend.Add(fmt.Sprintf("Percent coverage"), relativeCoverage)
	if err != nil {
		panic(err.Error())
	}

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "coverage.png"); err != nil {
		panic(err)
	}
}
