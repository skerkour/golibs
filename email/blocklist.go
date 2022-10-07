package email

import (
	"bufio"
	"bytes"
	_ "embed"
	"strings"
)

//go:embed domains_blocklist.txt
var BlocklistBytes []byte

var blocklist map[string]bool

func IsInBlocklist(email string) bool {
	if blocklist == nil {
		mailBlocklistFileReader := bytes.NewReader(BlocklistBytes)
		mailBlocklistScanner := bufio.NewScanner(mailBlocklistFileReader)
		for mailBlocklistScanner.Scan() {
			blocklist[strings.TrimSpace(mailBlocklistScanner.Text())] = true
		}
	}

	return blocklist[email]
}
