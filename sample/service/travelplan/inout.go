package travelplan

import (
	"github.com/jamestrandung/go-die/die"
	"github.com/jamestrandung/go-die/sample/dependencies/mapservice"
)

type plan interface {
	input
	output
}

type Dependencies interface {
	GetMapService() mapservice.Service
}

type input interface {
	Dependencies
	GetPointA() string
	GetPointB() string
}

type output interface {
	SetTravelPlan(die.AsyncResult)
}

type TravelPlan die.AsyncResult

func (p TravelPlan) GetTravelDistance() float64 {
	result := die.Outcome[mapservice.Route](p.Task)
	return result.Distance
}

func (p TravelPlan) GetTravelDuration() float64 {
	result := die.Outcome[mapservice.Route](p.Task)
	return result.Duration
}
