package http

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/caffix/cloudflare-roundtripper/cfrt"
)

const (
	UserAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36"
	Accept     = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
	AcceptLang = "en-US,en;q=0.8"
)

var (
	defaultClient *http.Client
)

func init() {
	defaultClient = &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          200,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   20 * time.Second,
			ExpectContinueTimeout: 20 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}
	defaultClient.Transport, _ = cfrt.New(defaultClient.Transport)
}

// RequestWebPage returns a string containing the entire response for
// the urlstring parameter when successful.
func RequestWebPage(urlstring string, body io.Reader, hvals map[string]string, uid, secret string) (string, error) {
	method := "GET"
	if body != nil {
		method = "POST"
	}
	req, err := http.NewRequest(method, urlstring, body)
	if err != nil {
		return "", err
	}
	if uid != "" && secret != "" {
		req.SetBasicAuth(uid, secret)
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", Accept)
	req.Header.Set("Accept-Language", AcceptLang)
	if hvals != nil {
		for k, v := range hvals {
			req.Header.Set(k, v)
		}
	}

	resp, err := defaultClient.Do(req)
	if err != nil {
		return "", err
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New(resp.Status)
	}

	in, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(in), nil
}
