package languages

import (
	_ "embed"
	"encoding/json"
)

//go:embed languages.json
var Bytes []byte

type Language struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	NativeName string `json:"native_name"`
}

var langs map[string]Language

func Get() (map[string]Language, error) {
	var err error

	if langs == nil {
		langs = map[string]Language{}
		err = json.Unmarshal(Bytes, &langs)
		if err != nil {

			return langs, err
		}
	}

	return langs, nil
}
