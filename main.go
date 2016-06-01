package sentryimporter

import "time"
import "fmt"

func main() {
	credentials := map[string]string{
		"organization_slug": "YOUR ORGANIZATION SLUG",
		"project_slug":      "YOUR PROJECT SLUG",
		"api_key":           "YOUR PROJECT API KEY",
	}
	client, err := NewClient(credentials)
	if err != nil {
		panic(err)
	}

	eventsCh := make(chan Event)
	from := time.Now().AddDate(0, -6, 0)
	to := time.Now()

	go func() {
		_, err := client.Events(eventsCh, from, to)
		if err != nil {
			panic(err)
		}
	}()

	for event := range eventsCh {
		fmt.Printf(event.Title)
	}
}
