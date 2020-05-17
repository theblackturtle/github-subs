package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/theblackturtle/github-subs/domainparser"
	"github.com/theblackturtle/github-subs/http"
	"github.com/theblackturtle/github-subs/stringset"
)

var nameStripRE = regexp.MustCompile(`^u[0-9a-f]{4}|20|22|25|2b|2f|3d|3a|40`)

func main() {
	var APIKey string
	var domain string
	var delay int
	flag.StringVar(&APIKey, "api", "", "Github API Key")
	flag.StringVar(&domain, "d", "", "Domain to search")
	flag.IntVar(&delay, "delay", 7, "Delay each request (second)")
	flag.Parse()

	if APIKey == "" {
		fmt.Fprintln(os.Stderr, "Please check your api key")
		os.Exit(1)
	}
	if domain == "" {
		fmt.Fprintln(os.Stderr, "Please check your domain")
		os.Exit(1)
	}
	headers := map[string]string{
		"Authorization": "token " + APIKey,
		"Content-Type":  "application/json",
	}

	// Setup domainparser
	domainparser := domainparser.NewDomainParser()
	domainparser.AddDomain(domain)
	re := domainparser.DomainRegex(domain)

	nameFilter := stringset.NewStringFilter()
	fetchNames := func(u string) {
		time.Sleep(time.Duration(delay) * time.Second)
		page, err := http.RequestWebPage(u, nil, nil, "", "")
		if err != nil {
			return
		}
		// Extract the subdomain names from the page
		for _, sd := range re.FindAllString(page, -1) {
			if name := cleanName(sd); name != "" && !nameFilter.Duplicate(name) {
				fmt.Println(name)
			}
		}
	}

	urlFilter := stringset.NewStringFilter()

loop:
	for i := 1; i <= 100; i++ {
		time.Sleep(time.Duration(delay) * time.Second)

		url := buildURL(domain, i)
		var result struct {
			Total int `json:"total_count"`
			Items []struct {
				URL   string  `json:"html_url"`
				Score float64 `json:"score"`
			} `json:"items"`
		}
		page, err := http.RequestWebPage(url, nil, headers, ",", "")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break loop
		}
		if err := json.Unmarshal([]byte(page), &result); err != nil {
			fmt.Fprintln(os.Stderr, err)
			break loop
		}
		for _, item := range result.Items {
			if t := modifyURL(item.URL); t != "" && !urlFilter.Duplicate(t) {
				go fetchNames(t)
			}
		}
	}

}

func buildURL(domain string, page int) string {
	pn := strconv.Itoa(page)
	u, _ := url.Parse("https://api.github.com/search/code")

	u.RawQuery = url.Values{
		"s":        {"indexed"},
		"type":     {"Code"},
		"o":        {"desc"},
		"q":        {"\"" + domain + "\""},
		"page":     {pn},
		"per_page": {"100"},
	}.Encode()
	return u.String()
}

func modifyURL(url string) string {
	m := strings.Replace(url, "https://github.com/", "https://raw.githubusercontent.com/", 1)
	m = strings.Replace(m, "/blob/", "/", 1)
	return m
}

// Clean up the names scraped from the web.
func cleanName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))

	for {
		name = strings.Trim(name, "-")

		if i := nameStripRE.FindStringIndex(name); i != nil {
			name = name[i[1]:]
		} else {
			break
		}
	}

	// Remove dots at the beginning of names
	if len(name) > 1 && name[0] == '.' {
		name = name[1:]
	}
	return name
}
