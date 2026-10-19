package appointment

import "time"

type CreateInput struct {
	ClientID     string
	StartAt      time.Time
	Description  string
	PlannedPrice int
	Deposit      int
}
