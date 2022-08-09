package cte

import (
	"reflect"
)

func isComplete(e Engine, planValue reflect.Value) error {
	planName := extractFullNameFromType(extractUnderlyingType(planValue))

	sd := newStructDisassembler()
	sd.extractAvailableMethods(planValue.Type())

	var verifyFn func(planName string, curPlanValue reflect.Value) error
	verifyFn = func(planName string, curPlanValue reflect.Value) error {
		ap := e.findAnalyzedPlan(planName, curPlanValue)

		for _, h := range ap.preHooks {
			expectedInout, ok := h.metadata.getInoutInterface()
			if !ok {
				return ErrInoutMetaMissing.Err(reflect.TypeOf(h.hook))
			}

			err := isInterfaceSatisfied(sd, expectedInout)
			if err != nil {
				return ErrPlanNotMeetingInoutRequirements.Err(planValue.Type(), expectedInout, err.Error())
			}
		}

		for _, component := range ap.components {
			if c, ok := e.computers[component.id]; ok {
				expectedInout, ok := c.metadata.getInoutInterface()
				if !ok {
					return ErrInoutMetaMissing.Err(component.id)
				}

				err := isInterfaceSatisfied(sd, expectedInout)
				if err != nil {
					return ErrPlanNotMeetingInoutRequirements.Err(planValue.Type(), expectedInout, err.Error())
				}
			}

			if _, ok := e.plans[component.id]; ok {
				nestedPlanValue := func() reflect.Value {
					if curPlanValue.Kind() == reflect.Pointer {
						return curPlanValue.Elem().Field(component.fieldIdx)
					}

					return curPlanValue.Field(component.fieldIdx)
				}()

				if err := verifyFn(component.id, nestedPlanValue); err != nil {
					return err
				}
			}
		}

		for _, h := range ap.postHooks {
			expectedInout, ok := h.metadata.getInoutInterface()
			if !ok {
				return ErrInoutMetaMissing.Err(reflect.TypeOf(h.hook))
			}

			err := isInterfaceSatisfied(sd, expectedInout)
			if err != nil {
				return ErrPlanNotMeetingInoutRequirements.Err(planValue.Type(), expectedInout, err.Error())
			}
		}

		return nil
	}

	return verifyFn(planName, planValue)
}

func isInterfaceSatisfied(sd structDisassembler, expectedInterface reflect.Type) error {
	for i := 0; i < expectedInterface.NumMethod(); i++ {
		rm := expectedInterface.Method(i)

		requiredMethod := extractMethodDetails(rm, false)

		methodSet, ok := sd.availableMethods[requiredMethod.name]
		if !ok {
			return ErrPlanMissingMethod.Err(requiredMethod)
		}

		if methodSet.Count() > 1 {
			return ErrPlanHavingAmbiguousMethods.Err(requiredMethod, methodSet)
		}

		foundMethod := methodSet.Items()[0]

		if !foundMethod.hasSameSignature(requiredMethod) {
			return ErrPlanHavingMethodButSignatureMismatched.Err(requiredMethod, foundMethod)
		}

		if sd.isAvailableMoreThanOnce(foundMethod) {
			return ErrPlanHavingSameMethodRegisteredMoreThanOnce.Err(foundMethod)
		}
	}

	return nil
}
