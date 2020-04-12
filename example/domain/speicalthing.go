package domain

import (
	eh "github.com/looplab/eventhorizon"
)

// SpecialThingAggregateType todo
const SpecialThingAggregateType = eh.AggregateType("SpecialThingAggregate")

type SpecialThingAggregate struct {
}
