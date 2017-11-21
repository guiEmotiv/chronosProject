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

	worker1 := core.NewWorker(base, 2, .1)
	worker1.Recorder.InternalStatus = true
	worker2 := core.NewWorker(base, 2, .1)
	worker2.Recorder.InternalStatus = true
	//worker3 := core.NewWorker(base, 2, .1)
	//worker3.Recorder.InternalStatus = true


	hyperSpace := core.NewDiscreteSpace(1*time.Millisecond)
	hyperSpace.AddWorker(worker1)
	hyperSpace.AddWorker(worker2)
	//hyperSpace.AddWorker(worker3)
	hyperSpace.AddPlaces(clients...)

	algorithm := core.NewGreedySearch(base, core.GetTasksFromPlaces(clients))

	manager := core.NewManager(hyperSpace, algorithm)
	
	manager.FillTaskListOrdinaryTasks(worker1)



	hyperSpace.MainLoop(func(c *core.DiscreteSpace) {
		manager.DistributeTasks()
		//pp.Println(manager.MapDriversTasks)
		manager.RefreshTasksStatus()
		pp.Println(len(manager.TasksList))
	}, 2000*time.Millisecond)

	plotters.PlotSimulation(palette.Plan9, worker1.Recorder, worker2.Recorder)


}




/*
recorderChannel := make(chan *core.WorkerRecorder)

	go func() {
		for {
			recorderChannel <- worker1.Recorder
		}
	}()


	func () {
		for value := range recorderChannel {
			drawNewPoint(value.XPositions[len(value.XPositions)-1], value.YPositions[len(value.YPositions)-1])
		}
	}()
*/