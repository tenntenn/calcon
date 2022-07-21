package google

import (
	"context"
	"fmt"
	"time"

	"github.com/tenntenn/calcon"
	"google.golang.org/api/calendar/v3"
)

type Option[T any] func(T) T

func TimeMin[T interface{ TimeMin(string) T }](t time.Time) Option[T] {
	return func(c T) T {
		return c.TimeMin(t.Format(time.RFC3339))
	}
}

func TimeMax[T interface{ TimeMax(string) T }](t time.Time) Option[T] {
	return func(c T) T {
		return c.TimeMax(t.Format(time.RFC3339))
	}
}

func UpdateMin[T interface{ UpdateMin(string) T }](t time.Time) Option[T] {
	return func(c T) T {
		return c.UpdateMin(t.Format(time.RFC3339))
	}
}

func Q[T interface{ Q(string) T }](q string) Option[T] {
	return func(c T) T {
		return c.Q(q)
	}
}

func Events(ctx context.Context, s *calendar.Service, id string, opts ...Option[*calendar.EventsListCall]) ([]*calcon.Event, error) {
	call := s.Events.List(id).Context(ctx)
	for _, opt := range opts {
		call = opt(call)
	}
	ges, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("calcon/google.Events(%q): %w", id, err)
	}

	es := make([]*calcon.Event, len(ges.Items))
	for i, ge := range ges.Items {
		startAt, err := toTime(ge.Start)
		if err != nil {
			return nil, err
		}

		endAt, err := toTime(ge.End)
		if err != nil {
			return nil, err
		}

		attendees := make([]string, len(ge.Attendees))
		for i := range ge.Attendees {
			attendees[i] = ge.Attendees[i].DisplayName
		}

		es[i] = &calcon.Event{
			ID:          ge.Id,
			Title:       ge.Summary,
			Description: ge.Description,
			StartAt:     startAt,
			EndAt:       endAt.In(startAt.Location()),
			Location:    ge.Location,
			Attendees:   attendees,
			TimeZone:    startAt.Location(),
		}
	}

	return es, nil
}

func toTime(edt *calendar.EventDateTime) (time.Time, error) {
	// TODO: all day event

	tm, err := time.Parse(time.RFC3339, edt.DateTime)
	if err != nil {
		return time.Time{}, err
	}

	loc, err := time.LoadLocation(edt.TimeZone)
	if err != nil {
		return time.Time{}, err
	}

	return tm.In(loc), nil
}
