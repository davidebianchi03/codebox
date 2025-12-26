package proxy

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func CreateReverseProxy(
	targetStr string,
	timeout time.Duration,
	keepAlive time.Duration,
	skipSSLVerify bool,
	headers map[string][]string,
) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse target, %s", err)
	}

	director := func(req *http.Request) {
		// force custom endpoint
		req.URL = target

		req.Host = target.Host
		if ip, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			req.Header.Set("X-Forwarded-For", ip)
		}
		req.Header.Set("X-Forwarded-Proto", req.URL.Scheme)
		req.Header.Set("X-Forwarded-Host", req.Host)
		for key, value := range headers {
			req.Header[key] = value
		}
	}

	return &httputil.ReverseProxy{Director: director, Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   timeout * time.Second,
			KeepAlive: keepAlive * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: skipSSLVerify},
	}}, nil
}
