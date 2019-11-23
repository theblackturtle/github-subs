package domainparser

import (
	"regexp"
	"strings"
	"sync"

	"github.com/theblackturtle/github-subs/stringset"
)

type DomainParser struct {
	sync.Mutex
	regexps map[string]*regexp.Regexp
	domains []string
}

func NewDomainParser() *DomainParser {
	return &DomainParser{}
}

// DomainRegex returns the Regexp object for the domain name identified by the parameter.
func (c *DomainParser) DomainRegex(domain string) *regexp.Regexp {
	c.Lock()
	defer c.Unlock()

	if re, found := c.regexps[domain]; found {
		return re
	}
	return nil
}

// AddDomains appends the domain names provided in the parameter to the list in the configuration.
func (c *DomainParser) AddDomains(domains []string) {
	for _, d := range domains {
		c.AddDomain(d)
	}
}

// AddDomain appends the domain name provided in the parameter to the list in the configuration.
func (c *DomainParser) AddDomain(domain string) {
	c.Lock()
	defer c.Unlock()

	// Check that the domain string is not empty
	d := strings.TrimSpace(domain)
	if d == "" {
		return
	}
	// Check that it is a domain with at least two labels
	labels := strings.Split(d, ".")
	if len(labels) < 2 {
		return
	}
	// Check that none of the labels are empty
	for _, label := range labels {
		if label == "" {
			return
		}
	}

	// Check that the regular expression map has been initialized
	if c.regexps == nil {
		c.regexps = make(map[string]*regexp.Regexp)
	}

	// Create the regular expression for this domain
	c.regexps[d] = SubdomainRegex(d)
	if c.regexps[d] != nil {
		// Add the domain string to the list
		c.domains = append(c.domains, d)
	}

	c.domains = stringset.Deduplicate(c.domains)
}

// Domains returns the list of domain names currently in the configuration.
func (c *DomainParser) Domains() []string {
	c.Lock()
	defer c.Unlock()

	return c.domains
}
