package user

import "time"

// User is the aggregate root for the User bounded context.
// ID is assigned by the database on first save (auto-increment).
// All mutations go through the exported methods on this type.
type User struct {
	id   uint
	name string

	// Uncommitted domain events – cleared after dispatch.
	events []DomainEvent
}

// -----------------------------------------------------------------
// Constructors
// -----------------------------------------------------------------

// NewUser creates a brand-new User and records a UserCreated event.
// The ID field is intentionally 0 – it will be set by the database
// after the first Save and reflected back via Reconstitute.
func NewUser(name string) (*User, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	u := &User{name: name}
	u.record(UserCreated{
		Name:       name,
		OccurredOn: time.Now(),
	})
	return u, nil
}

// Reconstitute rebuilds a User from persistence data without triggering events.
// Called by the repository when loading from the database.
func Reconstitute(id uint, name string) *User {
	return &User{id: id, name: name}
}

// SetID is called by the repository after the DB assigns an auto-increment ID.
func (u *User) SetID(id uint) {
	u.id = id
}

// -----------------------------------------------------------------
// Queries (getters)
// -----------------------------------------------------------------

func (u *User) ID() uint     { return u.id }
func (u *User) Name() string { return u.name }

// -----------------------------------------------------------------
// Commands
// -----------------------------------------------------------------

// Rename changes the user's name.
func (u *User) Rename(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	u.name = name
	return nil
}

// -----------------------------------------------------------------
// Domain event handling
// -----------------------------------------------------------------

// PullEvents returns all uncommitted events and clears the internal slice.
func (u *User) PullEvents() []DomainEvent {
	events := u.events
	u.events = nil
	return events
}

func (u *User) record(e DomainEvent) {
	u.events = append(u.events, e)
}
