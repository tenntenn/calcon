package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/k3a/html2text"
	"github.com/tenntenn/calcon"
	"github.com/tenntenn/calcon/google"
	"github.com/tenntenn/calcon/ics"
	"go.uber.org/multierr"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"google.golang.org/api/calendar/v3"
)

var idRegexp = regexp.MustCompile(`^\[(.+)\]`)

var (
	flagFormat string
	flagOutput string
)

func init() {
	flag.StringVar(&flagFormat, "format", "google", "output format[google-json, google-csv, ics]")
	flag.StringVar(&flagOutput, "output", "", "output file path, default value is empty (stdout)")
}

func main() {
	flag.Parse()
	ctx := context.Background()
	if err := run(ctx, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) (rerr error) {
	if len(args) == 0 {
		return errors.New("Calender ID must be specified")
	}
	id := args[0]

	s, err := calendar.NewService(ctx)
	if err != nil {
		return err
	}

	evts, err := google.Events(ctx, s, id)
	if err != nil {
		return err
	}

	switch flagFormat {
	case "google-csv", "google-json":
		if err := outputGoogle(evts); err != nil {
			return err
		}
	case "ics":
		if err := outputICSAll(evts); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected format: %q", flagFormat)
	}

	return nil
}

func outputGoogle(events []*calcon.Event) (rerr error) {
	if len(events) == 0 {
		return nil
	}

	gevents := make(map[string]string, len(events))
	for i := range events {
		id := idRegexp.FindString(events[i].Title)
		if len(id) < 3 {
			return fmt.Errorf("Title %q does not have id", events[i].Title)
		}
		events[i].Title = events[i].Title[len(id):]
		id = id[1 : len(id)-1]
		gevents[id] = google.New(events[i]).Link()
	}

	var w io.Writer = os.Stdout
	if flagOutput != "" {
		f, err := os.Create(flagOutput)
		if err != nil {
			return err
		}
		defer func() {
			rerr = multierr.Append(rerr, f.Close())
		}()
		w = f
	}

	switch flagFormat {
	case "google-csv":
		if err := outputGoogleCSV(w, gevents); err != nil {
			return err
		}
	case "google-json":
		if err := outputGoogleJSON(w, gevents); err != nil {
			return err
		}
	}

	return nil
}

func outputGoogleCSV(w io.Writer, links map[string]string) error {
	cw := csv.NewWriter(w)

	header := []string{"ID", "URL"}
	if err := cw.Write(header); err != nil {
		return err
	}

	ids := maps.Keys(links)
	slices.Sort(ids)
	for _, id := range ids {
		record := []string{id, links[id]}
		if err := cw.Write(record); err != nil {
			return err
		}
	}

	cw.Flush()
	if err := cw.Error(); err != nil {
		return err
	}

	return nil
}

func outputGoogleJSON(w io.Writer, links map[string]string) error {
	if err := json.NewEncoder(w).Encode(links); err != nil {
		return err
	}
	return nil
}

func outputICSAll(events []*calcon.Event) error {
	output := "ics"
	if flagOutput != "" {
		output = flagOutput
	}

	if err := os.MkdirAll(output, 0o744); err != nil {
		return err
	}

	for _, e := range events {
		if err := outputICS(output, e); err != nil {
			return err
		}
	}

	return nil
}

func outputICS(dir string, e *calcon.Event) (rerr error) {

	id := idRegexp.FindString(e.Title)
	if len(id) < 3 {
		return fmt.Errorf("Title %q does not have id", e.Title)
	}

	e.Title = e.Title[len(id):]
	id = id[1 : len(id)-1]

	e.Description = strings.ReplaceAll(html2text.HTML2Text(e.Description), "\r", "")

	fpath := filepath.Join(dir, id+".ics")
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}

	defer func() {
		rerr = multierr.Append(rerr, f.Close())
	}()

	if err := ics.Serialize(f, []*calcon.Event{e}); err != nil {
		return err
	}

	return nil
}
