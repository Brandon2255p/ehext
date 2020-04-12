package horizgen

import (
	"bytes"
	"os"
	"text/template"
)

func makeAggregateData(aggregateName string, events []EventData) []AggregateData {
	data := make([]AggregateData, len(events))
	for i, v := range events {
		data[i] = AggregateData{aggregateName, v}
	}
	return data
}

// AggregateData todo
type AggregateData struct {
	AggregateName string
	Event         EventData
}

// EventData todo
type EventData struct {
	Name        string
	Description string
}

// EventMember todo
type EventMember struct {
	MemberType string
	MemberName string
	MemberTag  string
}

// GenerateEvent todo
func GenerateEvent(e ...EventData) string {
	const templateEvent = `package domain

// Code generated .* DO NOT EDIT\.

import (
	eh "github.com/looplab/eventhorizon"
)

const (
{{range .}}	// {{.Name}}Type {{.Description}}
	{{.Name}}Type eh.EventType = "{{.Name}}"
{{end}}
)

func init() {
{{range .}}	eh.RegisterEventData({{.Name}}Type, func() eh.EventData {
		return &{{.Name}}{}
	})
{{end}}
}
`
	t := template.Must(template.New("event").Parse(templateEvent))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "event", e)
	return string(buff.Bytes())
}

// GenerateCommand todo
func GenerateCommand(aggregate string, e ...EventData) string {
	const templateCommand = `package domain

// Code generated .* DO NOT EDIT\.

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
)

const (
{{range .}}	// {{.Event.Name}}Type {{.Event.Description}}
	{{.Event.Name}}Type eh.CommandType = "{{.Event.Name}}"
{{end}})
{{range .}}
// AggregateID generate for {{.Event.Name}}
func (c *{{.Event.Name}}) AggregateID() uuid.UUID {
	return c.ID
}

// AggregateType generate for {{.Event.Name}}
func (c *{{.Event.Name}}) AggregateType() eh.AggregateType {
	return {{.AggregateName}}Type
}

// CommandType generate for {{.Event.Name}}
func (c *{{.Event.Name}}) CommandType() eh.CommandType {
	return {{.Event.Name}}Type
}
{{end}}`
	t := template.Must(template.New("command").Parse(templateCommand))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "command", makeAggregateData(aggregate, e))
	return string(buff.Bytes())
}

// GenerateRegisterCommand todo
func GenerateRegisterCommand(events []EventData) string {
	const temp = `package domain

// Code generated .* DO NOT EDIT\.


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
{{range .}}
	err = handleError(err, h.SetHandler(c, {{.Name}}Type)){{end}}
	return err
}
`
	t := template.Must(template.New("register").Parse(temp))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "register", events)
	return string(buff.Bytes())
}

// GenerateHandleCommand todo
func GenerateHandleCommand(aggregateName string, events []EventData) string {
	const temp = `package domain

// Code generated .* DO NOT EDIT\.


import (
	"context"
	"fmt"

	eh "github.com/looplab/eventhorizon"
)

// HandleCommand todo
func (a *{{(index . 0).AggregateName}}) HandleCommand(ctx context.Context, cmd eh.Command) error {
	switch cmd := cmd.(type) { {{range .}}
	case *{{.Event.Name}}:
		return a.handle{{.Event.Name}}(ctx, cmd){{end}}
	}
	return fmt.Errorf("Command %s not handled, run generator again", cmd.CommandType())
}
`
	t := template.Must(template.New("handle").Parse(temp))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "handle", makeAggregateData(aggregateName, events))
	return string(buff.Bytes())
}

// GenerateCommandHandlers todo
func GenerateCommandHandlers(aggregateName string, events []EventData) string {
	const temp = `package domain

import (
	"context"

	eh "github.com/looplab/eventhorizon"
)
{{range .}}
func (a *{{.AggregateName}})handle{{.Event.Name}}(ctx context.Context, cmd *{{.Event.Name}}) error {
	return nil
}
{{end}}`
	t := template.Must(template.New("handle").Parse(temp))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "handle", makeAggregateData(aggregateName, events))
	return string(buff.Bytes())
}

// Write todo
func Write(outputFolder, data string) error {
	f, err := os.Create(outputFolder)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(data)
	return err
}

// GenerateApplyEvent todo
func GenerateApplyEvent(aggregateName string, events []EventData) string {
	const temp = `package domain

// Code generated .* DO NOT EDIT\.


import (
	"context"
	"fmt"

	eh "github.com/looplab/eventhorizon"
)

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *{{(index . 0).AggregateName}}) ApplyEvent(ctx context.Context, event eh.Event) error {
	switch event.EventType() { {{range .}}
	case {{.Event.Name}}Type:
		return a.apply{{.Event.Name}}(ctx, event){{end}}
	}
	return fmt.Errorf("Event %s not handled, run generator again", event.EventType())
}
`
	t := template.Must(template.New("handle").Parse(temp))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "handle", makeAggregateData(aggregateName, events))
	return string(buff.Bytes())
}

// GenerateEventAppliers todo
func GenerateEventAppliers(aggregateName string, events []EventData) string {
	const temp = `package domain

import (
	"context"

	eh "github.com/looplab/eventhorizon"
)
{{range .}}
func (a *{{.AggregateName}}) apply{{.Event.Name}}(ctx context.Context, event eh.Event) error {
	_, ok := event.Data().({{.Event.Name}})
	if !ok {
		panic("INVALID apply{{.Event.Name}} is trying to convert invalid data")
	}
	return nil
}
{{end}}`

	data := makeAggregateData(aggregateName, events)
	t := template.Must(template.New("handle").Parse(temp))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "handle", data)
	return string(buff.Bytes())
}

// GenerateRegisterAggregate todo
func GenerateRegisterAggregate(aggregateName string) string {
	const temp = `package domain

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
)

func init() {
	eh.RegisterAggregate(func(id uuid.UUID) eh.Aggregate {
		return &{{ . }}{
			AggregateBase: events.NewAggregateBase({{ . }}Type, id),
		}
	})
}

// interface check
var _ = eh.Aggregate(&{{ . }}{})
`

	t := template.Must(template.New("handle").Parse(temp))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "handle", aggregateName)
	return string(buff.Bytes())
}
