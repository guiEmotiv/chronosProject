package core

type Place struct {
	XPos               float64 `json:"pos_x"`
	YPos               float64 `json:"pos_y"`
	Priority           float64 `json:"priority"`
	ExpectedArriveTime float64 `json:"expected_arrive_time"`
	OperationTime      float64 `json:"operation_time"`
	Arrived            float64 `json:"arrived"`
	Status             float64 `json:"status"`
}

func NewPlace(pX, pY float64, priority float64, expectedArriveTime float64,
	operationTime float64, arrived float64, L float64) *Place {

	place := new(Place)
	place.XPos = pX
	place.YPos = pY
	place.Priority = priority
	place.ExpectedArriveTime = expectedArriveTime
	place.OperationTime = operationTime
	place.Arrived = arrived
	place.Status = status
	return place
}
