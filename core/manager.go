package core

import (
	"strconv"
	"../util"
)

//type WorkerDriver struct {
//	Core *Worker
//}

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

func (manager *Manager) FillTaskListOrdinaryTasks(refWorker *Worker) {
	tasks := manager.Heuristic.GetSortedTasks(refWorker)
	manager.TasksList = tasks
}


func (manager *Manager) DistributeTasks() {
	for i:=0;i<len(manager.TasksList);i++ {
		idTask := strconv.Itoa(i)
		task := *manager.TasksList[idTask]

		if task.Status() == CREATED {
			for workerId, worker := range manager.Drivers {
				if worker.Status == STOPPED {
					worker.To = task.Address()
					worker.Status = TRAVELING

					manager.MapDriversTasks[workerId] = idTask
					task.SetStatus(ASSIGNED)
					break
				}
			}

		}
	}

}

//
// driver -> task (e.g. 12j4k -> i34u4)

func (manager *Manager) RefreshTasksStatus() {

	for driverId, taskId := range manager.MapDriversTasks {
		if manager.Drivers[driverId].Status == STOPPED {
			task := *(manager.TasksList[taskId])
			task.SetStatus(COMPLETED)
		}
	}
}

func (manager *Manager) CleanTaskList() {

	finalTasksList := make(map[string]*Task)

	for id, task := range manager.TasksList {
		dRefTask := *task

		if dRefTask.Status() != COMPLETED {
			finalTasksList[id] = task
		}
	}

	manager.TasksList = finalTasksList
}

func (manager *Manager) addEmergency(task *EmergencyTask) {

	idTask := util.NextId(5)


	manager.TasksList = map[string]*Task{} // Removed


	manager.Heuristic.GetTasks()


	for workerId, worker := range manager.Drivers {
		if worker.Status == TRAVELING || worker.Status == WAITING {
			idForCanceledTask := manager.MapDriversTasks[workerId]
			oldTask := *manager.TasksList[idForCanceledTask]
			oldTask.SetStatus(CANCELLED)

			worker.To = task.Address()
			worker.Status = TRAVELING

			manager.MapDriversTasks[workerId] = idTask
			task.SetStatus(ASSIGNED)

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