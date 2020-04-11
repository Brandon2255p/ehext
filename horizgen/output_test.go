package horizgen_test

import (
	"testing"

	"github.com/Brandon2255p/ehext/horizgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateEventType(t *testing.T) {
	input := horizgen.EventData{
		Name:        "ThingyCreatedEvent",
		Description: "This is how thingies are made",
	}
	expected := `// Code generated .* DO NOT EDIT\.
package domain

import (
	eh "github.com/looplab/eventhorizon"
)

const (
	// ThingyCreatedEventType This is how thingies are made
	ThingyCreatedEventType eh.EventType = "ThingyCreatedEvent"
)

func init() {
	eh.RegisterEventData(ThingyCreatedEventType, func() eh.EventData {
		return &ThingyCreatedEvent{}
	})
}
`
	output := horizgen.GenerateEvent(input)
	assert.Equal(t, expected, output)
}

func TestGenerateCommandType(t *testing.T) {
	input := horizgen.EventData{
		Name:        "CreateCommand",
		Description: "This is how thingies are made",
	}
	expected := `// Code generated .* DO NOT EDIT\.
package domain

import (
	eh "github.com/looplab/eventhorizon"
)

const (
	// CreateCommandType This is how thingies are made
	CreateCommandType eh.CommandType = "CreateCommand"
)
`
	output := horizgen.GenerateCommand(input)
	assert.Equal(t, expected, output)
}

func TestGenerateCommandRegisterType(t *testing.T) {
	input := []horizgen.EventData{{
		Name:        "TransferOwnershipCommand",
		Description: "This is how thingies are made",
	}, {
		Name:        "StartOwnershipCommand",
		Description: "This is how thingies are made",
	}}
	expected := `// Code generated .* DO NOT EDIT\.
package domain

import (
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/commandhandler/bus"
	"github.com/pkg/errors"
)

func handleError(err error, new error) error {
	if new == nil {
		return err
	}
	if err != nil {
		errors.Wrap(err, new.Error())
	}
	return err
}

// registerCommands register all commands
func registerCommands(h *bus.CommandHandler, c eh.CommandHandler) error {
	var err error

	err = handleError(err, h.SetHandler(c, TransferOwnershipCommandType))
	err = handleError(err, h.SetHandler(c, StartOwnershipCommandType))
	return err
}
`

	output := horizgen.GenerateRegisterCommand(input)
	assert.Equal(t, expected, output)
}

func TestGenerateHandleCommand(t *testing.T) {
	input := []horizgen.EventData{{
		Name:        "TransferOwnershipCommand",
		Description: "This is how thingies are made",
	}, {
		Name:        "StartOwnershipCommand",
		Description: "This is how thingies are made",
	}}
	expected := `// Code generated .* DO NOT EDIT\.
package domain

import (
	"context"
	"fmt"

	eh "github.com/looplab/eventhorizon"
)

// HandleCommand todo
func (a *ThingAggregate) HandleCommand(ctx context.Context, cmd eh.Command) error {
	switch cmd := cmd.(type) { 
	case *TransferOwnershipCommand:
		return a.handleTransferOwnershipCommand(ctx, cmd)
	case *StartOwnershipCommand:
		return a.handleStartOwnershipCommand(ctx, cmd)
	}
	return fmt.Errorf("Command %s not handled, run generator again", cmd.CommandType())
}
`

	output := horizgen.GenerateHandleCommand("ThingAggregate", input)
	assert.Equal(t, expected, output)
}

func TestGenerateApplyEvent(t *testing.T) {
	input := []horizgen.EventData{{
		Name:        "ThingyCreatedEvent",
		Description: "This is how thingies are made",
	}, {
		Name:        "Thingy2CreatedEvent",
		Description: "This is how thingies are made",
	}}
	expected := `// Code generated .* DO NOT EDIT\.
package domain

import (
	"context"
	"fmt"

	eh "github.com/looplab/eventhorizon"
)

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *ThingAggregate) ApplyEvent(ctx context.Context, event eh.Event) error {
	switch event.EventType() { 
	case ThingyCreatedEventType:
		return a.applyThingyCreatedEvent(ctx, event)
	case Thingy2CreatedEventType:
		return a.applyThingy2CreatedEvent(ctx, event)
	}
	return fmt.Errorf("Event %s not handled, run generator again", event.EventType())
}
`

	output := horizgen.GenerateApplyEvent("ThingAggregate", input)
	assert.Equal(t, expected, output)
}

func TestGenerateCommandHandlers(t *testing.T) {
	input := []horizgen.EventData{{
		Name:        "TransferOwnershipCommand",
		Description: "This is how thingies are made",
	}, {
		Name:        "StartOwnershipCommand",
		Description: "This is how thingies are made",
	}}
	expected := `package domain

import (
	"context"

	eh "github.com/looplab/eventhorizon"
)

func (a *ThingAggregate)handleTransferOwnershipCommand(ctx context.Context, cmd *TransferOwnershipCommand) error {
	return nil
}

func (a *ThingAggregate)handleStartOwnershipCommand(ctx context.Context, cmd *StartOwnershipCommand) error {
	return nil
}
`

	output := horizgen.GenerateCommandHandlers("ThingAggregate", input)
	assert.Equal(t, expected, output)
}

func TestGenerateEventAppliers(t *testing.T) {
	input := []horizgen.EventData{{
		Name:        "ThingyCreatedEvent",
		Description: "This is how thingies are made",
	}, {
		Name:        "Thingy2CreatedEvent",
		Description: "This is how thingies are made",
	}}
	expected := `package domain

import (
	"context"

	eh "github.com/looplab/eventhorizon"
)

func (a *ThingAggregate) applyThingyCreatedEvent(ctx context.Context, event eh.Event) error {
	_, ok := event.Data().(ThingyCreatedEvent)
	if !ok {
		panic("INVALID applyThingyCreatedEvent is trying to convert invalid data")
	}
	return nil
}

func (a *ThingAggregate) applyThingy2CreatedEvent(ctx context.Context, event eh.Event) error {
	_, ok := event.Data().(Thingy2CreatedEvent)
	if !ok {
		panic("INVALID applyThingy2CreatedEvent is trying to convert invalid data")
	}
	return nil
}
`

	output := horizgen.GenerateEventAppliers("ThingAggregate", input)
	assert.Equal(t, expected, output)
}

func TestWrite(t *testing.T) {
	input := horizgen.EventData{
		Name:        "ThingyCreatedEvent",
		Description: "This is how thingies are made",
	}
	err := horizgen.Write("./_generate.go", horizgen.GenerateEvent(input))
	require.NoError(t, err)
	require.FileExists(t, ("./_generate.go"))
}

// a.applyThingy2CreatedEvent(ctx context.Context, event Thingy2CreatedEvent) error {
// 	a.started = true
// }
