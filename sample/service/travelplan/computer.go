package travelplan

import (
	"context"
	"github.com/jamestrandung/go-die/sample/dependencies/mapservice"

	"github.com/jamestrandung/go-die/sample/config"
)

// Computers with external dependencies still has to register itself with the
// engine using init() so that we can perform validations on plans
func init() {
	// config.Print("travelplan")
	config.Engine.RegisterImpureComputer(TravelPlan{}, computer{})
	// config.Print(config.Engine)
}

type computer struct{}

func (c computer) Compute(ctx context.Context, p any) (any, error) {
	casted := p.(plan)

	route, err := casted.GetMapService().GetRoute(casted.GetPointA(), casted.GetPointB())
	if err != nil {
		return c.calculateStraightLineDistance(casted), nil
	}

	return route, nil
}

func (c computer) calculateStraightLineDistance(p plan) mapservice.Route {
	config.Printf("Building route from %s to %s using straight-line distance\n", p.GetPointA(), p.GetPointB())
	return mapservice.Route{
		Distance: 4,
		Duration: 5,
	}
}
