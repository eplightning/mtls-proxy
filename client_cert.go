package main

import (
	"crypto/tls"
	"sync"
	"time"
)

type DynamicClientCert struct {
	certPath string
	keyPath  string
	interval time.Duration

	nextRefresh time.Time
	mut         sync.Mutex
	cert        *tls.Certificate
}

func NewDynamicClientCert(certPath, keyPath string, interval time.Duration) *DynamicClientCert {
	return &DynamicClientCert{
		certPath: certPath,
		keyPath:  keyPath,
		interval: interval,
	}
}

func (c *DynamicClientCert) GetClientCertificate(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
	c.mut.Lock()
	defer c.mut.Unlock()

	if c.cert == nil || c.nextRefresh.Before(time.Now()) {
		cert, err := tls.LoadX509KeyPair(c.certPath, c.keyPath)
		if err != nil {
			return nil, err
		}

		c.cert = &cert
		c.nextRefresh = time.Now().Add(c.interval)
	}

	return c.cert, nil
}
