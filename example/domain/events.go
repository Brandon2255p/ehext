package domain

//go:generate ehext -aggregate SpecialThingAggregate -search ./$GOFILE -appliers
// Optionally add -appliers to the above command to generate the eventAppliers.go file

// ThingMadeEvent todo
type ThingMadeEvent struct {
}

// ThingDeletedEvent todo
type ThingDeletedEvent struct {
}
