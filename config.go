package main

import (
	"log"
	"net/url"
	"os"
	"strings"
)

type appConfig struct {
	listen          string
	target          *url.URL
	extraHeaders    map[string]string
	setForwardedFor bool
	forwardHost     bool
	overrideHost    string

	tlsServerName     string
	tlsSkipVerify     bool
	tlsServerCAPath   string
	tlsClientKeyPath  string
	tlsClientCertPath string
}

func loadConfig() *appConfig {
	envString := func(key string, defaultVal string) string {
		if val, found := os.LookupEnv(key); found {
			return val
		}

		return defaultVal
	}

	envBool := func(key string, defaultVal bool) bool {
		val := os.Getenv(key)

		if strings.EqualFold(val, "true") || val == "1" {
			return true
		} else if strings.EqualFold(val, "false") || val == "0" {
			return false
		}

		return defaultVal
	}

	envRequiredURL := func(key string) *url.URL {
		val := os.Getenv(key)
		if val == "" {
			log.Fatalf("required env var %v not specified", key)
			return nil
		}

		u, err := url.Parse(val)
		if err != nil {
			log.Fatalf("required env var %v is not a valid URL: %v", key, err)
			return nil
		}

		return u
	}

	envStringKeyValueMap := func(key string) map[string]string {
		output := make(map[string]string)

		for _, pair := range strings.Split(os.Getenv(key), "\n") {
			if strings.TrimSpace(pair) == "" {
				continue
			}

			name, value, _ := strings.Cut(pair, ":")
			output[strings.TrimSpace(name)] = strings.TrimSpace(value)
		}

		return output
	}

	return &appConfig{
		listen:          envString("LISTEN_ADDRESS", ":8080"),
		target:          envRequiredURL("TARGET_URL"),
		setForwardedFor: envBool("SET_FORWARDED_FOR", true),
		forwardHost:     envBool("FORWARD_HOST", false),
		overrideHost:    envString("OVERRIDE_HOST", ""),
		extraHeaders:    envStringKeyValueMap("EXTRA_HEADERS"),

		tlsServerName:     envString("TLS_SERVER_NAME", ""),
		tlsSkipVerify:     envBool("TLS_SKIP_VERIFY", false),
		tlsServerCAPath:   envString("TLS_SERVER_CA_PATH", ""),
		tlsClientCertPath: envString("TLS_CLIENT_CERT_PATH", ""),
		tlsClientKeyPath:  envString("TLS_CLIENT_KEY_PATH", ""),
	}
}
