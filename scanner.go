package sentryimporter

import (
	"io"
	"net/http"
)

type scanner struct {
	source *paginatedSource
	resp   *http.Response
	errMsg error
}

func newScanner(source *paginatedSource) *scanner {
	return &scanner{source: source}
}

func (this *scanner) scan() bool {
	if this.errMsg != nil {
		return false
	}

	if this.resp != nil {
		linkString := this.resp.Header.Get("Link")
		header, err := parsePaginationHeader(linkString)
		if err != nil || !header.next.results {
			this.errMsg = err
			return false
		}
		this.source.cursor = header.next.url
	}

	this.resp, this.errMsg = this.source.fetch()
	return this.errMsg == nil
}

func (this *scanner) body() io.ReadCloser {
	return this.resp.Body
}

func (this *scanner) err() error {
	return this.errMsg
}
