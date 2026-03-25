package eventbus

import (
	"log"

	"starter/internal/domain/todo"
	domainuser "starter/internal/domain/user"

	// The asaskevich/eventbus package name is "EventBus".
	// We alias it to "evbus" to avoid a clash with our own package name.
	evbus "github.com/asaskevich/eventbus"
)

// Bus is a thin wrapper around asaskevich/eventbus that provides
// typed subscription and publishing for domain events.
//
// Internally it uses reflection: subscriber functions must have the
// concrete event type as their single argument, e.g.:
//
//	bus.Subscribe("todo_list.created", func(e todo.TodoListCreated) { … })
type Bus struct {
	eb evbus.Bus
}

// New returns a synchronous in-process event bus.
// Events are dispatched in the calling goroutine, in the order they were published.
func New() *Bus {
	return &Bus{eb: evbus.New()}
}

// -----------------------------------------------------------------
// Subscription
// -----------------------------------------------------------------

// Subscribe registers a synchronous handler for the given topic.
// The handler must be a function whose single parameter matches the event type.
func (b *Bus) Subscribe(topic string, fn interface{}) error {
	return b.eb.Subscribe(topic, fn)
}

// SubscribeAsync registers a handler that runs in its own goroutine.
// When transactional is true only one handler goroutine runs at a time.
func (b *Bus) SubscribeAsync(topic string, fn interface{}, transactional bool) error {
	return b.eb.SubscribeAsync(topic, fn, transactional)
}

// WaitAsync blocks until all asynchronous handlers have finished.
func (b *Bus) WaitAsync() {
	b.eb.WaitAsync()
}

// -----------------------------------------------------------------
// Publishing
// -----------------------------------------------------------------

// Publish dispatches an arbitrary payload to all subscribers of topic.
func (b *Bus) Publish(topic string, payload interface{}) {
	b.eb.Publish(topic, payload)
}

// PublishAll is a generic helper that dispatches a slice of domain events.
// Any event type that exposes EventName() string can be used:
//
//	eventbus.PublishAll(bus, todoList.PullEvents())
//	eventbus.PublishAll(bus, user.PullEvents())
func PublishAll[E interface{ EventName() string }](b *Bus, events []E) {
	for _, e := range events {
		b.Publish(e.EventName(), e)
	}
}

// -----------------------------------------------------------------
// Default logging handlers (replace with real side-effects in production)
// -----------------------------------------------------------------

// RegisterDefaultHandlers wires up built-in logging handlers so that
// every domain event at minimum produces a log line.
func (b *Bus) RegisterDefaultHandlers() {
	// -- TodoList events --
	must(b.Subscribe("todo_list.created", func(e todo.TodoListCreated) {
		log.Printf("[EVENT] %s – list=%s title=%q", e.EventName(), e.TodoListID, e.Title)
	}))
	must(b.Subscribe("todo_list.line_added", func(e todo.TodoListLineAdded) {
		log.Printf("[EVENT] %s – list=%s line=%s desc=%q", e.EventName(), e.TodoListID, e.LineID, e.Desc)
	}))
	must(b.Subscribe("todo_list.line_completed", func(e todo.TodoListLineCompleted) {
		log.Printf("[EVENT] %s – list=%s line=%s", e.EventName(), e.TodoListID, e.LineID)
	}))
	must(b.Subscribe("todo_list.completed", func(e todo.TodoListCompleted) {
		log.Printf("[EVENT] %s – list=%s 🎉", e.EventName(), e.TodoListID)
	}))

	// -- User events --
	must(b.Subscribe("user.created", func(e domainuser.UserCreated) {
		log.Printf("[EVENT] %s – user=%d name=%q", e.EventName(), e.UserID, e.Name)
	}))
}

func must(err error) {
	if err != nil {
		panic("eventbus: failed to register handler: " + err.Error())
	}
}
