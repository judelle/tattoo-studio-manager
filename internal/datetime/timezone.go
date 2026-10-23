package datetime

import (
	"time"
	_ "time/tzdata"
)

var Moscow = mustLoadLocation("Europe/Moscow")

func mustLoadLocation(name string) *time.Location {
	location, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}

	return location
}
