package ics

import (
	"io"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/tenntenn/calcon"
)

func Serialize(w io.Writer, events []*calcon.Event) error {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)

	for _, e := range events {
		event := cal.AddEvent(e.ID)
		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		if e.TimeZone != nil {
			event.SetStartAt(e.StartAt.In(e.TimeZone))
			event.SetEndAt(e.EndAt.In(e.TimeZone))
		} else {
			event.SetStartAt(e.StartAt)
			event.SetEndAt(e.EndAt)
		}
		event.SetSummary(e.Title)
		event.SetLocation(e.Location)
		event.SetDescription(e.Description)
		event.SetURL(e.URL.String())

		for _, a := range e.Attendees {
			event.AddAttendee(a)
		}
	}

	if err := cal.SerializeTo(w); err != nil {
		return err
	}

	return nil
}
