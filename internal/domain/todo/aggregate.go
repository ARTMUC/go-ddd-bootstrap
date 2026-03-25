package todo

import (
	"time"

	"github.com/google/uuid"
)

// TodoList is the aggregate root. All mutations go through its methods.
// It accumulates domain events that are dispatched after the transaction commits.
type TodoList struct {
	id     uuid.UUID
	title  string
	status TodoStatus
	lines  []*TodoListLine

	// Uncommitted domain events – cleared after dispatch.
	events []DomainEvent
}

// -----------------------------------------------------------------
// Constructors
// -----------------------------------------------------------------

// NewTodoList creates a brand-new TodoList and records a TodoListCreated event.
func NewTodoList(title string) (*TodoList, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}
	tl := &TodoList{
		id:     uuid.New(),
		title:  title,
		status: StatusDraft,
	}
	tl.record(TodoListCreated{
		TodoListID: tl.id,
		Title:      title,
		OccurredOn: time.Now(),
	})
	return tl, nil
}

// Reconstitute rebuilds a TodoList from persistence without triggering events.
// Called by the repository when loading from the database.
func Reconstitute(
	id uuid.UUID,
	title string,
	status TodoStatus,
	lines []*TodoListLine,
) *TodoList {
	return &TodoList{
		id:     id,
		title:  title,
		status: status,
		lines:  lines,
	}
}

// -----------------------------------------------------------------
// Queries (getters)
// -----------------------------------------------------------------

func (tl *TodoList) ID() uuid.UUID          { return tl.id }
func (tl *TodoList) Title() string          { return tl.title }
func (tl *TodoList) Status() TodoStatus     { return tl.status }
func (tl *TodoList) Lines() []*TodoListLine { return tl.lines }

// -----------------------------------------------------------------
// Commands (business operations)
// -----------------------------------------------------------------

// AddLine appends a new item to the list.
func (tl *TodoList) AddLine(description string) (*TodoListLine, error) {
	if tl.status.IsCompleted() {
		return nil, ErrAlreadyCompleted
	}
	line, err := newTodoListLine(tl.id, description)
	if err != nil {
		return nil, err
	}
	tl.lines = append(tl.lines, line)

	tl.record(TodoListLineAdded{
		TodoListID: tl.id,
		LineID:     line.id,
		Desc:       description,
		OccurredOn: time.Now(),
	})
	return line, nil
}

// CompleteLine marks the given line as done.
// When every line is done the list itself transitions to Completed.
func (tl *TodoList) CompleteLine(lineID uuid.UUID) error {
	if tl.status.IsCompleted() {
		return ErrAlreadyCompleted
	}
	line, err := tl.findLine(lineID)
	if err != nil {
		return err
	}
	line.complete()

	tl.record(TodoListLineCompleted{
		TodoListID: tl.id,
		LineID:     lineID,
		OccurredOn: time.Now(),
	})

	if tl.allLinesDone() {
		tl.status = StatusCompleted
		tl.record(TodoListCompleted{
			TodoListID: tl.id,
			OccurredOn: time.Now(),
		})
	}
	return nil
}

// Activate moves the list from Draft to Active.
func (tl *TodoList) Activate() {
	if tl.status == StatusDraft {
		tl.status = StatusActive
	}
}

// -----------------------------------------------------------------
// Domain event handling
// -----------------------------------------------------------------

// PullEvents returns all uncommitted events and clears the internal slice.
// The application service calls this after the repository transaction
// succeeds and dispatches the events via the event bus.
func (tl *TodoList) PullEvents() []DomainEvent {
	events := tl.events
	tl.events = nil
	return events
}

func (tl *TodoList) record(e DomainEvent) {
	tl.events = append(tl.events, e)
}

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

func (tl *TodoList) findLine(id uuid.UUID) (*TodoListLine, error) {
	for _, l := range tl.lines {
		if l.id == id {
			return l, nil
		}
	}
	return nil, ErrLineNotFound
}

func (tl *TodoList) allLinesDone() bool {
	if len(tl.lines) == 0 {
		return false
	}
	for _, l := range tl.lines {
		if !l.done {
			return false
		}
	}
	return true
}
