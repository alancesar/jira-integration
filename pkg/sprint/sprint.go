package sprint

import "time"

const (
	Open   State = "open"
	Future State = "future"
	Closed State = "closed"
)

type (
	State string

	Sprint struct {
		ID          uint
		Name        string
		State       State
		Goal        string
		StartedAt   time.Time
		EndedAt     time.Time
		CompletedAt time.Time
	}
)
