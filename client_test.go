package sentryimporter

import "testing"

func TestNewClientWithoutOrganization(t *testing.T) {
	credentials := map[string]string{
		"project_slug": "projectslug",
		"api_key":      "apikey",
	}
	_, err := NewClient(credentials, nil)
	expected := "credentials[\"organization_slug\"] is empty"

	if err.Error() != expected {
		t.Errorf("expected to get (%s), got (%s)", err.Error(), expected)
	}
}

func TestNewClientWithoutProject(t *testing.T) {
	credentials := map[string]string{
		"organization_slug": "organization",
		"api_key":           "apikey",
	}
	_, err := NewClient(credentials, nil)
	expected := "credentials[\"project_slug\"] is empty"

	if err.Error() != expected {
		t.Errorf("expected to get (%s), got (%s)", err.Error(), expected)
	}
}

func TestNewClientWithoutApiKey(t *testing.T) {
	credentials := map[string]string{
		"organization_slug": "organization",
		"project_slug":      "projectslug",
	}
	_, err := NewClient(credentials, nil)
	expected := "credentials[\"api_key\"] is empty"

	if err.Error() != expected {
		t.Errorf("expected to get (%s), got (%s)", err.Error(), expected)
	}
}

func TestNewClient(t *testing.T) {
	credentials := map[string]string{
		"organization_slug": "organization",
		"project_slug":      "projectslug",
		"api_key":           "apikey",
	}
	client, err := NewClient(credentials, nil)

	if err != nil {
		t.Errorf("expected to get (nil), got (%s)", err)
	}

	if client == nil {
		t.Error("expected client not to be nil")
	}
}
