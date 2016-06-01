package sentryimporter

import "testing"

func TestParseEmptyDirection(t *testing.T) {
	_, err := parseDirection("")

	if err == nil || err.Error() != WRONG_HEADER_FORMAT {
		t.Errorf("excpected to get (%s), got (nil)", WRONG_HEADER_FORMAT)
	}
}

func TestParseMalformedSizeDirection(t *testing.T) {
	_, err := parseDirection("dir; d")

	if err == nil || err.Error() != WRONG_HEADER_FORMAT {
		t.Errorf("excpected to get (%s), got (nil)", WRONG_HEADER_FORMAT)
	}
}

func TestParseMalformedContentDirection(t *testing.T) {
	_, err := parseDirection("l; rel; res; curs")

	if err == nil || err.Error() != WRONG_HEADER_FORMAT {
		t.Errorf("excpected to get (%s), got (nil)", WRONG_HEADER_FORMAT)
	}
}

func TestParseDirection(t *testing.T) {
	result, err := parseDirection("<http://app.getsentry.com>; rel=\"Next\"; results=\"true\"; cursor=\"23412\"")

	if err != nil {
		t.Errorf("excpected to get (nil), got (%v)", err)
	}

	if result.url != "http://app.getsentry.com" {
		t.Errorf("excpected to get (http://app.getsentry.com), got (%s)", result.url)
	}

	if result.rel != "Next" {
		t.Errorf("excpected to get (Next), got (%s)", result.rel)
	}

	if result.results != true {
		t.Errorf("excpected to get (true), got (%v)", result.results)
	}

	if result.cursor != "23412" {
		t.Errorf("excpected to get (23412), got (%s)", result.cursor)
	}
}

func TestParsePaginationHeader(t *testing.T) {
	header := "<http://app.getsentry.com/previous>; rel=\"Previous\"; results=\"true\"; cursor=\"previous\", <http://app.getsentry.com/next>; rel=\"Next\"; results=\"true\"; cursor=\"next\""
	result, err := parsePaginationHeader(header)

	if err != nil {
		t.Errorf("excpected to get (nil), got (%v)", err)
	}

	if result.previous == nil {
		t.Error("excpected not to get (nil), got (nil)")
	}

	if result.next == nil {
		t.Error("excpected not to get (nil), got (nil)")
	}
}
