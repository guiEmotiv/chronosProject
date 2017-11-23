package core

import (
	"time"
	"../util"
)

type DiscreteSpace struct {
	Places   map[string]*Place  `json:"places"`
	Workers  map[string]*Worker `json:"workers"`

	Status int `json:"status"`

	TickTime                  time.Duration `json:"tick_time"`
	ElapsedTime               time.Duration
	CurrentTime               time.Time
	SimulationInitializedTime time.Time
}

func NewDiscreteSpace(stepTime time.Duration) *DiscreteSpace {
	ds := DiscreteSpace{
		TickTime: stepTime,
		Status: 0,
		Places: make(map[string]*Place),
		Workers: make(map[string]*Worker),
	}
	return &ds
}

func (space *DiscreteSpace) tick() {
	elapsedTime := space.CurrentTime.Sub(space.SimulationInitializedTime)
	space.ElapsedTime = elapsedTime
}

func (space *DiscreteSpace) AddPlaces(places ...*Place) {
	for _, place  := range places {
		id := util.NextId(5)
		space.Places[id] = place
	}
}

func (space *DiscreteSpace) AddWorker(w *Worker) {
	id := util.NextId(5)
	space.Workers[id] = w
}






func (space *DiscreteSpace) MainLoop(extraHandler func(context *DiscreteSpace), maxDuration ...time.Duration) {
	if len(maxDuration) > 0 && len(maxDuration) < 2 {
		space.SimulationInitializedTime = time.Now()
		finalTime := space.SimulationInitializedTime.Add(maxDuration[0])
		for finalTime.Sub(space.CurrentTime) > 0 {

			extraHandler(space)
			space.Step()
			space.CurrentTime = time.Now()
			time.Sleep(space.TickTime)
		}
	} else {
		space.SimulationInitializedTime = time.Now()
		for {

			extraHandler(space)
			space.Step()
			space.CurrentTime = time.Now()
			time.Sleep(space.TickTime)

		}
	}

}

func (space *DiscreteSpace) Step() {
	space.tick()

	for _, worker := range space.Workers {
		worker.NextPosition(space)
		worker.RefreshStatus(space)
	}

}
