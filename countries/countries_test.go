package countries

import "testing"

func TestGetMap(t *testing.T) {
	expectedNumberOfCountries := 245
	countries, err := GetMap()
	if err != nil {
		t.Error(err)
	}

	if len(countries) != expectedNumberOfCountries {
		t.Errorf("Invalid number of countries. Got %d, expected: %d", len(countries), expectedNumberOfCountries)
	}
}

func TestGetList(t *testing.T) {
	expectedNumberOfCountries := 245
	countries, err := GetList()
	if err != nil {
		t.Error(err)
	}

	if len(countries) != expectedNumberOfCountries {
		t.Errorf("Invalid number of countries. Got %d, expected: %d", len(countries), expectedNumberOfCountries)
	}
}
