package domain

//go:generate ehext -aggregate SpecialThingAggregate -search ./$GOFILE -handlers
// Optionally add -handlers to the above command to generate the commandhandlers.go file

// StartThingCommand todo
type StartThingCommand struct {
}

// EndThingCommand todo
type EndThingCommand struct {
}
