package email

import _ "embed"

//go:embed domains_blocklist.txt
var BlocklistBytes []byte
