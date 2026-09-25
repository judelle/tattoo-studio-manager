package appointment

import "time"

type Status string

const (
	StatusPlanned   Status = "planned"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
)

type Appointment struct {
	ID           string
	ClientID     string
	StartAt      time.Time
	Description  string
	PlannedPrice int
	FinalPrice   int
	Deposit      int
	Status       Status
	SketchPaths  []string
}
