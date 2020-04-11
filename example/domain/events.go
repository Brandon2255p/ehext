package domain

//go:generate ehext -aggregate SpecialThingAggregate -search ./$GOFILE

type ThingMadeEvent struct {
}

type ThingDeletedEvent struct {
}
