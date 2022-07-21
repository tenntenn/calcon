package google

import (
	"net/url"
	"time"

	"github.com/tenntenn/calcon"
)

const (
	urlPrefix = "https://www.google.com/calendar/render?"
	layout    = "20060102T150405Z0700"
)

type Event struct {
	UUID     string
	Text     string
	StartAt  time.Time
	EndAt    time.Time
	Details  string
	Location string
	URL      string
	TimeZone *time.Location
}

func (e *Event) Link() string {
	val := make(url.Values)

	val.Set("action", "TEMPLATE")

	if e.UUID != "" {
		val.Set("uuid", e.UUID)
	}

	if e.Text != "" {
		val.Set("text", e.Text)
	}

	if !e.StartAt.IsZero() {
		dates := e.StartAt.Format(layout)
		if !e.EndAt.IsZero() {
			dates += "/" + e.EndAt.Format(layout)
		}
		val.Set("dates", dates)
	}

	if e.Details != "" {
		val.Set("details", e.Details)
	}

	if e.Location != "" {
		val.Set("location", e.Location)
	}

	if e.URL != "" {
		val.Set("url", e.Location)
	}

	if e.TimeZone != nil {
		val.Set("url", e.Location)
	}

	return urlPrefix + val.Encode()
}

func New(e *calcon.Event) *Event {
	startAt := e.StartAt
	endAt := e.EndAt
	if e.TimeZone != nil {
		startAt = startAt.In(e.TimeZone)
		endAt = endAt.In(e.TimeZone)
	}

	return &Event{
		UUID:     e.ID,
		Text:     e.Title,
		StartAt:  startAt,
		EndAt:    endAt,
		Details:  e.Description,
		Location: e.Location,
		URL:      e.URL.String(),
		TimeZone: e.TimeZone,
	}
}
