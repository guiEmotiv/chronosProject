package core

import (
	"math"
	"sort"
	"strconv"
)

type GreedySearch struct {
	Base   *Place
	Tasks  map[string]*Task
}

func NewGreedySearch(base *Place, tasks []*Task) *GreedySearch {
	tasksMap := make(map[string]*Task)

	for i := 0; i < len(tasks); i++ {
		id := strconv.Itoa(i)
		tasksMap[id] = tasks[i]
	}

	gs := GreedySearch{
		Base:   base,
		Tasks:  tasksMap,
	}

	return &gs
}

func (algorithm *GreedySearch) GetBase() *Place {
	return algorithm.Base
}

func (algorithm *GreedySearch) GetTasks() map[string]*Task {
	return algorithm.Tasks
}

func (algorithm *GreedySearch) GetWeight(id string, worker *Worker) float64 {
	task := *algorithm.Tasks[id]
	goToPlace := task.Address()

	cXPos := worker.XPos
	cYPos := worker.YPos

	dstToPlace := math.Sqrt(math.Pow(goToPlace.XPos-cXPos, 2) + math.Pow(goToPlace.YPos-cYPos, 2))

	dstToBase := math.Sqrt(math.Pow(algorithm.Base.XPos-goToPlace.XPos, 2) + math.Pow(algorithm.Base.YPos-goToPlace.YPos, 2))

	arrivedTimeToPlace := dstToPlace * 1.0 // Because the velocity is always 1.0
	timeToGoToBase := dstToBase * 1.0      // Because same above reason

	totalTime := arrivedTimeToPlace + goToPlace.OperationTime + timeToGoToBase

	return math.Pow(goToPlace.Priority/totalTime, 3)
}

func (algorithm *GreedySearch) GetSortedTasks(worker *Worker) map[string]*Task {

	weights := make(map[float64]string)
	keys := make([]float64, 0)

	for id, _ := range algorithm.Tasks {
		w := algorithm.GetWeight(id, worker)
		weights[w] = id
		keys = append(keys, w)
	}

	sort.Float64s(keys)

	finalTasks := make(map[string]*Task, 0)

	nKeys := len(keys)
	for i := 0; i < nKeys; i++ {
		k := keys[nKeys-1-i]

		id := weights[k]

		finalTasks[strconv.Itoa(i)] = algorithm.Tasks[id]
	}

	return finalTasks

}

func (algorithm *GreedySearch) SetTasks(newTasks map[string]*Task) {
	algorithm.Tasks = newTasks
}
func (algorithm *GreedySearch) SetBase(place *Place) {
	algorithm.Base = place
}