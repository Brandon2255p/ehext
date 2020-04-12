package domain

import "github.com/google/uuid"

//go:generate ehext -aggregate SpecialThingAggregate -search ./$GOFILE -handlers
// Optionally add -handlers to the above command to generate the commandhandlers.go file

// StartThingCommand todo
type StartThingCommand struct {
	ID uuid.UUID
}

// EndThingCommand todo
type EndThingCommand struct {
	ID uuid.UUID
}
