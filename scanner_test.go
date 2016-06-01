package sentryimporter

import "testing"
import "errors"
import "net/http"
import "bytes"
import "io/ioutil"

func testNewScanner(t *testing.T) {
	source := &paginatedSource{}

	result := newScanner(source)

	if result == nil {
		t.Error("expected new instance of scanner to be created, got (nil)")
	}
}

func TestScanError(t *testing.T) {
	source := &paginatedSource{
		handler: func(url string) (*http.Response, error) {
			return nil, errors.New("end")
		},
	}
	aScanner := newScanner(source)

	result := aScanner.scan()

	if result {
		t.Errorf("expected result to be (false), got (true)")
	}
}

func TestScanSuccess(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString("Hello World"))
	source := &paginatedSource{
		handler: func(url string) (*http.Response, error) {
			return &http.Response{Body: body}, nil
		},
	}
	aScanner := newScanner(source)

	result := aScanner.scan()

	if !result {
		t.Errorf("expected result to be (true), got (false)")
	}

	if aScanner.body() != body {
		t.Error("expected result to be a given body")
	}
}

func TestBody(t *testing.T) {
	response := &http.Response{}
	aScanner := scanner{resp: response}

	result := aScanner.body()

	if result != response.Body {
		t.Error("expected to get response body, got (%v)", result)
	}
}

func TestErr(t *testing.T) {
	message := "error message"
	aScanner := scanner{errMsg: errors.New(message)}

	result := aScanner.err()

	if result.Error() != message {
		t.Errorf("expected (%s), got (%s)", message, result)
	}
}
