package common

import "management/core/users"

var urgentMap = map[string]float64{
	"difficult": 1,
	"timely":    4,
	"score":     1,
}
var nonUrgentMap = map[string]float64{
	"difficult": 1.5,
	"timely":    1.5,
	"score":     3,
}

func AllocationAlgorithm(employeeScore users.EmployeeScore) (noUrgentScore, urgentScore float64) {
	noUrgentScore = (float64(employeeScore.DifficultScore)*nonUrgentMap["difficult"] +
		float64(employeeScore.TimelyScore)*nonUrgentMap["timely"] +
		float64(employeeScore.OrderScore)*nonUrgentMap["score"]) / float64(employeeScore.OrderCount)
	urgentScore = (float64(employeeScore.DifficultScore)*urgentMap["difficult"] +
		float64(employeeScore.TimelyScore)*urgentMap["timely"] +
		float64(employeeScore.OrderScore)*urgentMap["score"]) / float64(employeeScore.OrderCount)
	return
}
