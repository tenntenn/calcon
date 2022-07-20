# calcon

calcon gets events of Google calendar via API and outputs ics files or Google calendar tempalte links (CSV or JSON).

## Install

```
$ go install github.com/tenntenn/calcon/cmd/calcon@latest
```

## How to use

```
# Output to ics files
$ export GOOGLE_APPLICATION_CREDENTIAL=<path to credential file> 
$ calcon -output ics -format ics <Google Calendar ID>

# Output to Google calndear links(CSV)
$ export GOOGLE_APPLICATION_CREDENTIAL=<path to credential file> 
$ calcon -output calendar.csv -format google-csv <Google Calendar ID>

# Output to Google calndear links(JSON)
$ export GOOGLE_APPLICATION_CREDENTIAL=<path to credential file> 
$ calcon -output calendar.json -format google-json <Google Calendar ID>
```
