package cortexbot

import (
	"regexp"
	"strings"

	valid "github.com/asaskevich/govalidator"
)

const (
	// Modified version of govalidator.DNSName that does not allow domain without tld (like localhost)
	DNSName string = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})+[\._]?$`
	Hash    string = `^[a-fA-F0-9]{32,}$`
)

var (
	rxDNS  = regexp.MustCompile(DNSName)
	rxHash = regexp.MustCompile(Hash)
)

// IsHash checks if a given string is a hash
// BUG(ilyaglow): not supported hashes are hashes shorter 32 letters and ssdeep
func IsHash(str string) bool {
	return rxHash.MatchString(str)
}

// IsDNSName is a modified version of function IsDNSName from https://github.com/asaskevich/govalidator/blob/master/patterns.go
func IsDNSName(str string) bool {
	if str == "" || len(strings.Replace(str, ".", "", -1)) > 255 {
		return false
	}
	return !valid.IsIP(str) && rxDNS.MatchString(str)
}
