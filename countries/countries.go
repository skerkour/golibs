package countries

import (
	_ "embed"
	"encoding/json"
)

//go:embed countries.json
var Bytes []byte

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

var countriesMap map[string]Country
var countriesList []Country

func GetMap() (map[string]Country, error) {
	if countriesMap == nil {
		list, err := GetList()
		if err != nil {
			return countriesMap, err
		}

		countriesMap = map[string]Country{}
		for _, country := range list {
			countriesMap[country.Code] = country
		}
	}

	return countriesMap, nil
}

func GetList() ([]Country, error) {
	var err error

	if countriesList == nil {
		countriesList = []Country{}

		err = json.Unmarshal(Bytes, &countriesList)
		if err != nil {

			return countriesList, err
		}
	}

	return countriesList, nil
}
