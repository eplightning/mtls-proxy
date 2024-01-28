package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

func configureServer(conf *appConfig) *http.Server {
	rt, err := httpTransport(conf)
	if err != nil {
		log.Fatalf("error when configuring HTTP transport: %v", err)
		return nil
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			httpRewrite(r, conf)
		},
		Transport: rt,
	}

	return &http.Server{
		Addr:    conf.listen,
		Handler: proxy,
	}
}

func httpRewrite(r *httputil.ProxyRequest, conf *appConfig) {
	if conf.setForwardedFor {
		r.Out.Header["X-Forwarded-For"] = r.In.Header["X-Forwarded-For"]
		r.SetXForwarded()
	}

	r.SetURL(conf.target)

	if conf.forwardHost {
		r.Out.Host = r.In.Host
	} else if conf.overrideHost != "" {
		r.Out.Host = conf.overrideHost
	}

	for name, val := range conf.extraHeaders {
		r.Out.Header[name] = []string{val}
	}
}

func httpTransport(conf *appConfig) (http.RoundTripper, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: conf.tlsSkipVerify,
		ServerName:         conf.tlsServerName,
	}

	if conf.tlsServerCAPath != "" {
		caFile, err := os.ReadFile(conf.tlsServerCAPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read CA from %v: %w", conf.tlsServerCAPath, err)
		}

		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM(caFile); !ok {
			return nil, fmt.Errorf("unable to parse any certificates from %v", conf.tlsServerCAPath)
		}

		tlsConfig.RootCAs = pool
	}

	if conf.tlsClientCertPath != "" && conf.tlsClientKeyPath != "" {
		cert := NewDynamicClientCert(conf.tlsClientCertPath, conf.tlsClientKeyPath, 5*time.Minute)
		tlsConfig.GetClientCertificate = cert.GetClientCertificate
	}

	rt := http.DefaultTransport.(*http.Transport).Clone()
	rt.TLSClientConfig = tlsConfig

	return rt, nil
}
