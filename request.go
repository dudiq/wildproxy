package main

import (
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Function that modifes the request
// Proxies requests in the form of
// http://www.example.com/https://www.upstream.com/foo/bar
// TODO support http scheme?
func proxyRequest(req *http.Request) {
	targetUrl := strings.TrimPrefix(req.URL.Path, "/")
	if len(req.URL.RawQuery) > 0 {
		targetUrl += "?" + req.URL.RawQuery
	}

	if !strings.HasPrefix(targetUrl, "http") {
		targetUrl = "https://" + targetUrl
	}

	u, err := url.Parse(targetUrl)
	if err != nil {
		log.Errorf("Request URL error: %s\n", err)
		return
	}

	req.URL = u
	req.Host = u.Host

	setOutHeaders(req)

	log.Infof("Proxying request to %s", req.URL)
}

func setOutHeaders(req *http.Request) {
	req.Header.Set("X-Forwarded-Proto", req.URL.Scheme)
	req.Header.Set("X-Forwarded-Host", req.Host)

	// set by transport
	req.Header.Del("Accept-Encoding")

	urlHeaders := []string{"Referer", "Origin"}
	for _, name := range urlHeaders {
		trimRootUrl(name, req)
	}

}

// Strips off root URL from headers like Referer
func trimRootUrl(name string, req *http.Request) {
	hv := req.Header.Get(name)
	if hv == "" {
		return
	}
	rootPath := rootUrl.String()
	if !strings.HasSuffix(rootPath, "/") {
		rootPath += "/"
	}
	hv = strings.TrimPrefix(hv, rootPath)
	req.Header.Set(name, hv)
}
