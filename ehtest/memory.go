package ehtest

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
	cb "github.com/looplab/eventhorizon/commandhandler/bus"
	eventbus "github.com/looplab/eventhorizon/eventbus/local"
	eventstore "github.com/looplab/eventhorizon/eventstore/memory"
	"github.com/stretchr/testify/assert"
)

var testCtx = eh.NewContextWithNamespace(context.Background(), "test")

var aggregateStore *events.AggregateStore
var estore *eventstore.EventStore
var aggType eh.AggregateType

var startTime time.Time

// Eveestorent is used to load startup events
type Event struct {
	Type eh.EventType
	Data eh.EventData
}

// NewMemoryFixture r
func NewMemoryFixture(aggregateType eh.AggregateType, fixtureTime time.Time) (ctx context.Context,
	bus *eventbus.EventBus,
	store *eventstore.EventStore,
	commandBus *cb.CommandHandler) {
	aggType = aggregateType
	startTime = fixtureTime

	group := eventbus.NewGroup()
	bus = eventbus.NewEventBus(group)
	store = eventstore.NewEventStore()
	estore = store
	commandBus = cb.NewCommandHandler()

	ctx = testCtx
	aggregateStore, _ = events.NewAggregateStore(store, bus)
	return
}

// AssertEvents f
func AssertEvents(t *testing.T, aggregateID uuid.UUID, expected []Event) {
	actual, err := estore.Load(testCtx, aggregateID)
	assert.NoError(t, err)
	assert.Len(t, actual, len(expected))
	for i, event := range actual {
		assert.EqualValues(t, expected[i].Data, event.Data())
		assert.EqualValues(t, expected[i].Type, event.EventType())
		assert.EqualValues(t, i+1, event.Version())
	}
}

func LoadStartupEvents(t *testing.T, aggregateID uuid.UUID, expected []Event) {
	evt := make([]eh.Event, len(expected))
	for i, e := range expected {
		aggEvent := eh.NewEventForAggregate(
			e.Type,
			e.Data, startTime, aggType, aggregateID, i+1)
		evt[i] = aggEvent
	}
	estore.Save(testCtx, evt, 0)
	agg, _ := aggregateStore.Load(testCtx, aggType, aggregateID)
	aggregateStore.Save(testCtx, agg)
}
