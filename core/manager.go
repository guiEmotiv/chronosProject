package core

import (
	"../util"
	"sort"
	"strconv"
	"math"
)

//type WorkerDriver struct {
//	Core *Worker
//}

type CoverageRecorder struct {
	GlobalCoverage   []float64 // Accumulated
	RelativeCoverage   []float64 // Percent
}

func NewCoverageRecorder() *CoverageRecorder {
	coverageTape := CoverageRecorder{
		GlobalCoverage: make([]float64, 0),
		RelativeCoverage: make([]float64, 0),
	}
	return &coverageTape
}


type Manager struct {
	Space   *DiscreteSpace
	Drivers map[string]*Worker
	MapDriversTasks map[string]string
	Heuristic Heuristic

	TasksList map[string]*Task
}

func NewManager(space *DiscreteSpace, heuristic Heuristic) *Manager {
	drivers := space.Workers

	m := Manager{
		Space: space,
		Drivers: drivers,
		Heuristic: heuristic,
		TasksList: make(map[string]*Task),
		MapDriversTasks: make(map[string]string),
	}
	return &m
}

func (manager *Manager) getTasksCoverage() ([][]int, []string) {
	finalCoverageMatrix := make([][]int, 0)
	keys := make([]string, 0)


	for key, place := range manager.Space.Places {
		driverArray := make([]int, 0)

		for id := range manager.Drivers {
			dx := math.Pow(place.XPos-manager.Drivers[id].XPos, 2)
			dy := math.Pow(place.YPos-manager.Drivers[id].YPos, 2)

			dist := math.Sqrt(dx + dy)
			// log.Println("Driver Position:", manager.Drivers[id].XPos, manager.Drivers[id].YPos)
			if dist < manager.Drivers[id].ActionRadio {
				driverArray = append(driverArray, 1)
			} else {
				driverArray = append(driverArray, 0)
			}
		}
		finalCoverageMatrix = append(finalCoverageMatrix, driverArray)
		keys = append(keys, key)

	}

	return finalCoverageMatrix, keys

}

// [[0 1 0 1 1 0], [0 0 0 1 1 1]]
// [0 1 0 1 1 0] OR [0 0 0 1 1 1] -> 0 1 0 1 1 1

func (manager *Manager) GetAccumulateCoverage() float64 {

	const lambda = 2.5
	const e = 0.1
	const p = 0.1

	taskList := make(map[string]*Task)
	for key, place := range manager.Space.Places {
		task := Task(&OrdinaryTask{
			ToWhere: place,
			State: CREATED,
		})
		taskList[key] = &task
	}

	coverageMatrix, keys := manager.getTasksCoverage()

	mixedCoverage := make([]int, 0)

	for _, t := range coverageMatrix {
		orResult := 0
		for _, dr := range t {
			orResult |= dr
		}
		mixedCoverage = append(mixedCoverage, orResult)
	}

	finalCoverage := 0.0

	for i, q := range mixedCoverage {
		task := *taskList[keys[i]]
		r := math.Exp(-1*lambda*p) * (1 - math.Exp(-1*lambda*e))*float64(q)*task.Address().Priority
		finalCoverage += r
	}

	return finalCoverage
}

func (manager *Manager) GetRelativeCoverage() float64{
	coverageMatrix, _ := manager.getTasksCoverage()

	mixedCoverage := make([]int, 0)

	sumCoverage := 0

	for _, t := range coverageMatrix {
		orResult := 0
		for _, dr := range t {
			orResult |= dr
		}
		if orResult == 0 {
			continue
		}
		sumCoverage += 1
		mixedCoverage = append(mixedCoverage, orResult)
	}

	placesSize := len(manager.Space.Places)

	return float64(sumCoverage)/float64(placesSize)

}

func (manager *Manager) RecordAllCoverages(cov *CoverageRecorder) {
	aCov := manager.GetAccumulateCoverage()
	rCov := manager.GetRelativeCoverage()

	cov.GlobalCoverage = append(cov.GlobalCoverage, aCov)
	cov.RelativeCoverage = append(cov.RelativeCoverage, rCov)
}


func (manager *Manager) FillTaskListOrdinaryTasks(refWorker *Worker) {
	tasks := manager.Heuristic.GetSortedTasks(refWorker)
	manager.TasksList = tasks
}

func getKeysOfTaskList(manager *Manager, withExclude bool) []string {
	keys := make([]string, 0)

	for k := range manager.TasksList {
		if withExclude {
			if _, err := strconv.Atoi(k); err == nil {
				keys = append(keys, k)
			}
		}else{
			keys = append(keys, k)
		}

	}
	return keys
}

func (manager *Manager) DistributeTasks() {
	sortedKeys := getKeysOfTaskList(manager, true) //Necessary cause we need make a priority

	sort.Sort(SortString(sortedKeys))

	for i:=0;i<len(sortedKeys);i++ {
		task := *manager.TasksList[sortedKeys[i]]
		if task.Status() == CREATED || task.Status() == CANCELLED {
			for workerId, worker := range manager.Drivers {
				if worker.Status == STOPPED {
					worker.To = task.Address()
					worker.Status = TRAVELING
					manager.MapDriversTasks[workerId] = sortedKeys[i]
					task.SetStatus(ASSIGNED)
					break
				}
			}

		}
	}

}

//
// driver -> task (e.g. 12j4k -> i34u4)

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}


func (manager *Manager) RefreshTasksStatus() {
	for driverId, taskId := range manager.MapDriversTasks {
		idsOfTasks := getKeysOfTaskList(manager, false)
		if contains(idsOfTasks, taskId) {
			if manager.Drivers[driverId].Status == STOPPED || manager.Drivers[driverId].LastTaskFinished {
				task := *manager.TasksList[taskId]
				task.SetStatus(COMPLETED)
				manager.Drivers[driverId].LastTaskFinished = false


				delete(manager.TasksList, taskId) // DELETE IT


				manager.DistributeTasks()

			}
		}

	}

	// Debounce, I'll implement tomorrow
	for _, taskId := range manager.MapDriversTasks {
		if !contains(getKeysOfTaskList(manager, false), taskId) {

		}
	}

}


func (manager *Manager) AddEmergency(task *EmergencyTask) {

	idTask := util.NextId(5)


	for workerId, worker := range manager.Drivers {
		if worker.Status == TRAVELING || worker.Status == WAITING {
			idForCanceledTask := manager.MapDriversTasks[workerId]
			oldTask := *manager.TasksList[idForCanceledTask]
			oldTask.SetStatus(CANCELLED)
			worker.To = task.Address()
			worker.Status = TRAVELING

			worker.ActionRadio = 0

			manager.MapDriversTasks[workerId] = idTask
			task.SetStatus(ASSIGNED)

			manager.Heuristic.SetTasks(manager.TasksList)

			newTasks := manager.Heuristic.GetSortedTasks(worker)

			manager.TasksList = newTasks
			pTask := Task(task)
			manager.TasksList[idTask] = &pTask

			manager.DistributeTasks()
			break
		}
	}
}

func NewEmergencyTask(place *Place) *EmergencyTask {
	return &EmergencyTask{
		ToWhere: place,
		State: CREATED,
	}
}
