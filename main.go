package main

import (
	"./core"
	"./parser"
	
	"time"
	"./plotters"
	"image/color/palette"
	"github.com/k0kubun/pp"
)


func main() {
	places := parser.GetPlacesFromJson("./data/dataPlaces.json")
	
	base := places[0]
	clients := places[1:]

	worker1 := core.NewWorker(base, 7.5, 1)
	worker1.Recorder.Record = true
	worker2 := core.NewWorker(base, 7.5, 1)
	worker2.Recorder.Record = true


	hyperSpace := core.NewDiscreteSpace(1*time.Millisecond)
	hyperSpace.AddWorker(worker1)
	hyperSpace.AddWorker(worker2)
	hyperSpace.AddPlaces(clients...)

	algorithm := core.NewGreedySearch(base, core.GetTasksFromPlaces(clients))

	manager := core.NewManager(hyperSpace, algorithm)
	
	manager.FillTaskListOrdinaryTasks(worker1)
	//pp.Println(manager.TasksList)


	coverageTape := core.NewCoverageRecorder()

	pp.Println(manager)

	hyperSpace.MainLoop(func(c *core.DiscreteSpace) {
		manager.RefreshTasksStatus()

		manager.DistributeTasks()

		manager.RecordAllCoverages(coverageTape)

	}, 200*time.Millisecond)

	plotters.PlotSimulation(palette.Plan9, hyperSpace)

	plotters.PLotCoverage(coverageTape)
}


