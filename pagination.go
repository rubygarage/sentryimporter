package sentryimporter

import (
	"errors"
	"net/http"
	"strings"
)

const WRONG_HEADER_FORMAT = "Wrong header format"

type paginationInfo struct {
	url     string
	rel     string
	results bool
	cursor  string
}

type paginationHeader struct {
	previous *paginationInfo
	next     *paginationInfo
}

type paginatedSource struct {
	cursor  string
	handler func(string) (*http.Response, error)
}

func parseDirection(direction string) (*paginationInfo, error) {
	if direction == "" {
		return nil, errors.New(WRONG_HEADER_FORMAT)
	}

	parse := func(str string) string {
		value := strings.Split(str, "=")
		return value[1][1 : len(value[1])-1]
	}

	vals := strings.Split(direction, "; ")
	if len(vals) != 4 {
		return nil, errors.New(WRONG_HEADER_FORMAT)
	}

	link, rel, results, cursor := vals[0], vals[1], vals[2], vals[3]
	if len(link) < 2 || len(rel) < 6 || len(results) < 8 || len(cursor) < 6 {
		return nil, errors.New(WRONG_HEADER_FORMAT)
	}

	return &paginationInfo{
		url:     link[1 : len(link)-1],
		rel:     parse(rel),
		results: parse(results) == "true",
		cursor:  parse(cursor),
	}, nil
}

func parsePaginationHeader(value string) (*paginationHeader, error) {
	directions := strings.Split(value, ", ")
	previous, err := parseDirection(directions[0])
	if err != nil {
		return nil, err
	}
	next, err := parseDirection(directions[1])
	if err != nil {
		return nil, err
	}
	return &paginationHeader{previous: previous, next: next}, nil
}

func (this *paginatedSource) fetch() (*http.Response, error) {
	return this.handler(this.cursor)
}
