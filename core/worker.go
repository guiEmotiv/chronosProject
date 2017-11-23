package core

import (
	"math"
	"time"
)

const (
	STOPPED   = -1
	TRAVELING = 0
	ARRIVED   = 1
	WAITING   = 2
	OPERATING = 3
)

type Worker struct {
	XPos        float64 `json:"pos_x"`
	YPos        float64 `json:"pos_y"`
	From        *Place  `json:"from"`
	To          *Place  `json:"to"`
	Status      int     `json:"status"`
	ActionRadio float64 `json:"action_radio"`
	Velocity    float64 `json:"velocity"`

	BaseActionRadio float64

	LastTaskFinished bool

	Recorder *WorkerRecorder `json:"recorder"`

	TimeBase time.Duration

	InitOperatingTime time.Duration
}

type WorkerRecorder struct {
	XPositions []float64
	YPositions []float64
	States     []map[int]float64

	Record bool
}

func NewWorker(initialPlace *Place, radio float64, velocity float64) *Worker {
	worker := new(Worker)
	recorder := new(WorkerRecorder)

	worker.XPos = initialPlace.XPos
	worker.YPos = initialPlace.YPos

	worker.From = initialPlace
	worker.To = initialPlace

	worker.Status = STOPPED
	worker.ActionRadio = radio
	worker.BaseActionRadio = radio

	worker.Velocity = velocity

	worker.TimeBase = time.Millisecond

	worker.Recorder = recorder
	worker.Recorder.States = make([]map[int]float64, 1)

	worker.LastTaskFinished = false

	return worker
}

func (worker *Worker) NextPosition(c *DiscreteSpace) {
	if worker.Recorder.Record {
		worker.Recorder.XPositions = append(worker.Recorder.XPositions, worker.XPos)
		worker.Recorder.YPositions = append(worker.Recorder.YPositions, worker.YPos)
		worker.Recorder.States = append(worker.Recorder.States, map[int]float64{worker.Status: c.ElapsedTime.Seconds()})
	}

	// fromPlace := worker.From
	toPlace := worker.To
	x1, y1 := worker.XPos, worker.YPos
	x2, y2 := toPlace.XPos, toPlace.YPos

	//m = (x2-x1)/(y2-y1)
	//
	//y - y1 = m(x - x1)

	if worker.Status == TRAVELING {
		x3 := x2 - x1
		y3 := y2 - y1

		mod3 := math.Sqrt(math.Pow(x3, 2) + math.Pow(y3, 2))

		vx := x3 / mod3
		vx = vx * worker.Velocity

		vy := y3 / mod3
		vy = vy * worker.Velocity

		worker.XPos = worker.XPos + vx
		worker.YPos = worker.YPos + vy
	}

}

func (worker *Worker) checkIfArrived(c *DiscreteSpace) bool {

	offset := 0.5

	if worker.Status == TRAVELING {
		// log.Println("Changing from traveling to arrived")
		if worker.XPos < worker.To.XPos+offset && worker.XPos > worker.To.XPos-offset {
			if worker.YPos < worker.To.YPos+offset && worker.YPos > worker.To.YPos-offset {

				worker.XPos = worker.To.XPos
				worker.YPos = worker.To.YPos

				//worker.Status = ARRIVED
				worker.Status = STOPPED

				return true
			}

		}

	}
	return false
}

func (worker *Worker) checkIfIsWaiting(c *DiscreteSpace) bool {
	if worker.Status == ARRIVED || worker.Status == WAITING {
		// log.Println("Changing from arrived  to waiting")

		expectedDuration := time.Duration(int64(worker.To.ExpectedArriveTime))
		trueExpectedTime := expectedDuration * worker.TimeBase

		if c.ElapsedTime < trueExpectedTime {
			worker.Status = WAITING
			return true
		}
	}
	return false
}

func (worker *Worker) checkIfIsOperating(c *DiscreteSpace) bool {
	if worker.Status == ARRIVED || worker.Status == WAITING {
		// log.Println("Changing from waiting  to operating")
		expectedDuration := time.Duration(int64(worker.To.ExpectedArriveTime))
		trueExpectedTime := expectedDuration * worker.TimeBase
		if c.ElapsedTime > trueExpectedTime {
			worker.Status = OPERATING
			worker.ActionRadio = 0
			worker.InitOperatingTime = time.Duration(c.ElapsedTime)
			return true

		}

	}

	if worker.Status == OPERATING {
		worker.Status = OPERATING
		dR := worker.BaseActionRadio/worker.To.OperationTime
		worker.ActionRadio += dR
		// worker.To.OperationTime
	}

	return false
}

func (worker *Worker) checkForNextTravel(c *DiscreteSpace) bool {
	if worker.Status == OPERATING {
		expectedOperatingTime := time.Duration(int64(worker.To.OperationTime)) * worker.TimeBase

		// middle := finalOperatingTime.
		// log.Println("Init at: ", worker.InitOperatingTime)

		if c.ElapsedTime-worker.InitOperatingTime > expectedOperatingTime {

			worker.Status = STOPPED
			worker.LastTaskFinished = true
			worker.ActionRadio = worker.BaseActionRadio

			return true
		}
	}

	return false

}

func (worker *Worker) RefreshStatus(c *DiscreteSpace) {
	worker.checkIfArrived(c)
	worker.checkIfIsWaiting(c)
	worker.checkIfIsOperating(c)
	worker.checkForNextTravel(c)

}
