sentryimporter
=================

Golang client for importing Sentry Events.

## Usage

```golang
package main

import "time"
import "fmt"

import "github.com/RubyGarage/sentryimporter"

func main() {
	credentials := map[string]string{
		"organization_slug": "YOUR ORGANIZATION SLUG",
		"project_slug":      "YOUR PROJECT SLUG",
		"api_key":           "YOUR PROJECT API KEY",
	}
	client, err := sentryimporter.NewClient(credentials)
  if err != nil {
    panic(err)
  }

	eventsCh := make(chan sentryimporter.Event)
	from := time.Now().AddDate(0, -6, 0)
	to := time.Now()

	go client.Events(eventsCh, from, to)

	for event := range eventsCh {
		fmt.Printf(event.Title)
	}
}

```
