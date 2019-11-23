package domainparser

import (
	"regexp"
	"strings"
)

// SUBRE is a regular expression that will match on all subdomains once the domain is appended.
const SUBRE = "(([a-zA-Z0-9]{1}|[_a-zA-Z0-9]{1}[_a-zA-Z0-9-]{0,61}[a-zA-Z0-9]{1})[.]{1})+"

// SubdomainRegex returns a Regexp object initialized to match
// subdomain names that end with the domain provided by the parameter.
func SubdomainRegex(domain string) *regexp.Regexp {
	// Change all the periods into literal periods for the regex
	d := strings.Replace(domain, ".", "[.]", -1)

	return regexp.MustCompile(SUBRE + d)
}
