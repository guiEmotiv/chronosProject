package core

const (
	CREATED = 0

	ASSIGNED  = 1
	COMPLETED = 2
	CANCELLED = 3
)

type Task interface {
	Address() *Place
	Status() int
	SetStatus(int)
}

type OrdinaryTask struct {
	ToWhere *Place
	State   int
}
type EmergencyTask struct {
	ToWhere *Place
	State   int
}

func (task *OrdinaryTask) Address() *Place {
	return task.ToWhere
}

func (task *OrdinaryTask) Status() int {
	return task.State
}

func (task *OrdinaryTask) SetStatus(state int) {
	task.State = state
}


func (task *EmergencyTask) Address() *Place {
	return task.ToWhere
}

func (task *EmergencyTask) Status() int {
	return task.State
}

func (task *EmergencyTask) SetStatus(state int) {
	task.State = state
}


func GetTasksFromPlaces(places []*Place) []*Task {
	tasks := make([]*Task, 0)
	for _, place := range places {
		t := &OrdinaryTask{
			ToWhere: place,
			State: CREATED,
		}
		genericTask := Task(t)
		tasks = append(tasks, &genericTask)
	}
	return tasks
}