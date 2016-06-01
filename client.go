package sentryimporter

import "net/http"
import "time"
import "fmt"
import "encoding/json"
import "io"
import "errors"
import "io/ioutil"

type Credentials map[string]string

type Client struct {
	Credentials
	httpClient *http.Client
	url        string
}

type Info struct {
	Count int    `json:"count"`
	Label string `json:"label"`
}

type Project struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Tags struct {
	Logger     Info `json:"logger"`
	ServerName Info `json:"server_name"`
	Level      Info `json:"level"`
}

type Event struct {
	Permalink    string    `json:"permalink"`
	IsPublic     bool      `json:"isPublic"`
	Culprit      string    `json:"culprit"`
	Title        string    `json:"title"`
	Level        string    `json:"level"`
	Annotations  []string  `json:"annotations"`
	HasSeen      bool      `json:"hasSeen"`
	IsBookmarked bool      `json:"isBookmarked"`
	Project      Project   `json:"project"`
	AssignedTo   string    `json:"assignedTo"`
	Tags         Tags      `json:"tags"`
	TimeSpent    string    `json:"timeSpent"`
	ShareId      string    `json:"shareId"`
	Id           string    `json:"id"`
	FirstSeen    time.Time `json:"firstSeen"`
	LastSeen     time.Time `json:"lastSeen"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewClient(credentials Credentials) (*Client, error) {
	if credentials["organization_slug"] == "" {
		return nil, errors.New("credentials[\"organization_slug\"] is empty")
	}

	if credentials["project_slug"] == "" {
		return nil, errors.New("credentials[\"project_slug\"] is empty")
	}

	if credentials["api_key"] == "" {
		return nil, errors.New("credentials[\"api_key\"] is empty")
	}

	url := "https://app.getsentry.com"
	return &Client{
		Credentials: credentials,
		httpClient:  &http.Client{},
		url:         url,
	}, nil
}

func (this *Client) Events(events chan<- Event, from time.Time, to time.Time) (int, error) {
	groups := make(chan Event)
	go func() {
		err := this.Groups(groups, from, to)
		if err != nil {
			panic(err)
		}
	}()
	return this.EventsForGroups(groups, events, from, to)
}

func (this *Client) EventsForGroups(groups <-chan Event, events chan<- Event, from time.Time, to time.Time) (int, error) {
	// TODO: change resolution to 1m or 10s
	// 1m resolution is resolved in 500 Internal Server Error :(
	// 10s resolution always returns 0 events in stats
	defer close(events)
	eventsEndpoint := "%s/api/0/groups/%s/stats/?resolution=1h&since=%d&until=%d"
	since := from.Unix()
	until := to.Unix()
	counter := 0
	apiKey := this.Credentials["api_key"]

	for group := range groups {
		var stats [][2]int
		url := fmt.Sprintf(eventsEndpoint, this.url, group.Id, since, until)
		resp, err := performRequest(this.httpClient, "GET", url, apiKey, nil)
		if err != nil {
			return counter, err
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&stats)
		if err != nil {
			return counter, err
		}
		for _, stat := range stats {
			if stat[1] != 0 {
				for i := 0; i < stat[1]; i++ {
					e := group
					events <- e
					counter += 1
				}
			}
		}
	}

	return counter, nil
}

func (this *Client) Groups(groups chan<- Event, from time.Time, to time.Time) error {
	defer close(groups)

	since := 0
	until := to.Unix()
	groupsEndpoint := "%s/api/0/projects/%s/%s/groups/?since=%d&until=%d"
	organization := this.Credentials["organization_slug"]
	project := this.Credentials["project_slug"]
	apiKey := this.Credentials["api_key"]

	resource := paginatedSource{
		cursor: fmt.Sprintf(groupsEndpoint, this.url, organization, project, since, until),
		handler: func(url string) (*http.Response, error) {
			return performRequest(this.httpClient, "GET", url, apiKey, nil)
		},
	}

	scanner := newScanner(&resource)
	for scanner.scan() {
		err := decodeGroups(scanner, groups)
		if err != nil {
			return err
		}
	}
	return scanner.err()
}

func decodeGroups(scanner *scanner, groups chan<- Event) error {
	body := scanner.body()
	defer body.Close()
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	var k []Event
	err = json.Unmarshal(content, &k)
	if err != nil {
		return errors.New(string(content))
	}

	for _, v := range k {
		groups <- v
	}
	return nil
}

func performRequest(client *http.Client, verb, url, apiKey string, payload io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(verb, url, payload)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(apiKey, "")
	return client.Do(req)
}
