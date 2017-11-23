package core

type Heuristic interface {
	GetWeight(id string, worker *Worker) float64
	GetSortedTasks(worker *Worker) map[string]*Task
	GetTasks() map[string]*Task
	GetBase() *Place

	SetTasks(map[string]*Task)
	SetBase(place *Place)
}
