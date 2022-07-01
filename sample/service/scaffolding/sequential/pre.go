package sequential

import "github.com/jamestrandung/go-die/sample/config"

type pre interface {
	GetTravelCost() float64
	SetTotalCost(float64)
}

type preHook struct{}

func (preHook) PreExecute(p any) error {
	config.Print("Before executing sequential plan")
	casted := p.(pre)

	casted.SetTotalCost(casted.GetTravelCost())

	return nil
}
