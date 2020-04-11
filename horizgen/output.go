package horizgen

import (
	"bytes"
	"os"
	"text/template"
)

const templateCommand = `package domain

import (
	eh "github.com/looplab/eventhorizon"
)

const (
	// {{.Event}}Type {{.Description}}
	{{.Event}}Type eh.CommandType = "{{.Event}}"
)
`

type EventData struct {
	Event       string
	Description string
}

type EventMember struct {
	MemberType string
	MemberName string
	MemberTag  string
}

func GenerateEvent(e EventData) string {
	const templateEvent = `package domain

import (
	eh "github.com/looplab/eventhorizon"
)

const (
	// {{.Event}}Type {{.Description}}
	{{.Event}}Type eh.EventType = "{{.Event}}"
)

func init() {
	eh.RegisterEventData({{.Event}}Type, func() eh.EventData {
		return &{{.Event}}{}
	})
}
`
	t := template.Must(template.New("event").Parse(templateEvent))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "event", e)
	return string(buff.Bytes())
}

func GenerateCommand(e EventData) string {
	t := template.Must(template.New("command").Parse(templateCommand))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "command", e)
	return string(buff.Bytes())
}

func GenerateRegisterCommand(events []EventData) string {
	const temp = `package domain

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
	err = handleError(err, h.SetHandler(c, {{.Event}}Type)){{end}}
	return err
}
`
	t := template.Must(template.New("register").Parse(temp))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "register", events)
	return string(buff.Bytes())
}

func GenerateHandleCommand(events []EventData) string {
	const temp = `package domain

import (
	"context"
	"fmt"

	eh "github.com/looplab/eventhorizon"
)

// HandleCommand todo
func (a *PlantAggregate) HandleCommand(ctx context.Context, cmd eh.Command) error {
	switch cmd := cmd.(type) { {{range .}}
	case *{{.Event}}:
		return a.handle{{.Event}}(ctx, cmd){{end}}
	}
	return fmt.Errorf("Command %s not handled, run generator again", cmd.CommandType())
}
`
	t := template.Must(template.New("handle").Parse(temp))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "handle", events)
	return string(buff.Bytes())
}
func GenerateCommandHandlers(events []EventData) string {
	const temp = `package domain

import (
	"context"

	eh "github.com/looplab/eventhorizon"
)
{{range .}}
func (a *PlantAggregate)handle{{.Event}}(ctx context.Context, cmd *{{.Event}}) error {
	return nil
}
{{end}}`
	t := template.Must(template.New("handle").Parse(temp))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "handle", events)
	return string(buff.Bytes())
}

func Write(outputFolder, data string) error {
	f, err := os.Create(outputFolder)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(data)
	return err
}

func GenerateApplyEvent(events []EventData) string {
	const temp = `package domain

import (
	"context"
	"fmt"

	eh "github.com/looplab/eventhorizon"
)

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *PlantAggregate) ApplyEvent(ctx context.Context, event eh.Event) error {
	switch event.EventType() { {{range .}}
	case {{.Event}}Type:
		return a.apply{{.Event}}(ctx, event){{end}}
	}
	return fmt.Errorf("Event %s not handled, run generator again", event.EventType())
}
`
	t := template.Must(template.New("handle").Parse(temp))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "handle", events)
	return string(buff.Bytes())
}

func GenerateEventAppliers(events []EventData) string {
	const temp = `package domain

import (
	"context"

	eh "github.com/looplab/eventhorizon"
)
{{range .}}
func (a *PlantAggregate) apply{{.Event}}(ctx context.Context, event eh.Event) error {
	_, ok := event.Data().({{.Event}})
	if !ok {
		panic("INVALID apply{{.Event}} is trying to convert invalid data")
	}
	return nil
}
{{end}}`
	t := template.Must(template.New("handle").Parse(temp))
	var buff bytes.Buffer
	t.ExecuteTemplate(&buff, "handle", events)
	return string(buff.Bytes())
}
