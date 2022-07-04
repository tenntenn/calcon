package calcon

import (
	"net/url"
	"time"
)

type Event struct {
	ID          string
	Title       string
	Description string
	StartAt     time.Time
	EndAt       time.Time
	Location    string
	URL         url.URL
	Attendees   []string
	TimeZone    *time.Location
}
