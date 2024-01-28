# mTLS Reverse proxy

Simple HTTP reverse proxy with support for mTLS connection.

Useful when your backend server requires TLS client certificates and your clients don't. This application can be used to close this gap.

### Container

Built and ready to use container image is available for `amd64`, `arm64` and `arm` architectures:

```
ghcr.io/eplightning/mtls-proxy:v0.1
```

### Configuration

`mtls-proxy` has several configuration options that can be specified using environment variables.

The only required environment variable is `TARGET_URL`, used to specify URL of backend server where the proxied requests should be forwarded to.

#### Defaults

By default, `mtls-proxy` will listen on HTTP connections on port `8080`.

`X-Forwarded-For`, `X-Forwarded-Host`, `X-Forwarded-Proto` headers will be set. `X-Forwarded-For` will be appended if it exists.

`Host` header will be set to `TARGET_URL` host.

TLS connections to backend server will be verified using system's root certificate pool.

#### Environment variables

```
LISTEN_ADDRESS    - listening address (default :8080)
TARGET_URL        - URL of backend server where the proxied requests should be forwarded to (required)
SET_FORWARDED_FOR - set X-Forwarded-For headers [true|false|1|0] (default true)
FORWARD_HOST      - use original request's Host (takes precedence over OVERRIDE_HOST) [true|false|1|0] (default false)
OVERRIDE_HOST     - Host header value to use (default none)
EXTRA_HEADERS     - Additional headers to set in format HeaderName:HeaderValue separated by new lines (\n) (default none)

TLS_SERVER_NAME      - Override TLS server name used for validation (default none)
TLS_SKIP_VERIFY      - Don't verify server certificates [true|false|1|0] (default false)
TLS_SERVER_CA_PATH   - Path to custom CA bundle to use for server certificate validation (default none)
TLS_CLIENT_CERT_PATH - Path to client certificate PEM to use (default none)
TLS_CLIENT_KEY_PATH  - Path to client key PEM to use (default none)
```

Note that both `TLS_CLIENT_CERT_PATH` and `TLS_CLIENT_KEY_PATH` need to be specified in order to client certificate to be presented.

#### Config reload

Files specified via `TLS_CLIENT_CERT_PATH` and `TLS_CLIENT_KEY_PATH` will be reloaded every 5 minutes.

At the moment `TLS_SERVER_CA_PATH` is only loaded once at the application start.
