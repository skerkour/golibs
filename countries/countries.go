package countries

import _ "embed"

//go:embed countries.json
var ConfigCountries []byte
