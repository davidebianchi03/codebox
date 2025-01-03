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

func CreateReverseProxy(targetStr string, timeout time.Duration, keepAlive time.Duration, skipSSLVerify bool, headers map[string][]string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse target, %s", err)
	}
	director := func(req *http.Request) {
		req.URL = target
		for key, value := range headers {
			req.Header[key] = value
		}

		// TODO: set x forwarded headers
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
