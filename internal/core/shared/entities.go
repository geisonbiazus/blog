package shared

import "time"

type Event struct {
	Type       string
	OccurredOn time.Time
	Payload    map[string]interface{}
}
