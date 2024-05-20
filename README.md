## Secure Server Template

My personal server setup based on go-blueprints. The Server setup manages HTTP and gRPC servers with optional TLS configurations.

### Features
- Go web framework: [Echo](https://github.com/labstack/echo)
- Database: [PostgreSQL](https://www.postgresql.org/)

### Middleware Overview

- **Content-Security-Policy**: Restricts content sources to mitigate XSS attacks.
- **Default Headers**: Enhances security with default response headers.
- **HSTS**: Enforces HTTPS usage to prevent protocol downgrade attacks.
- **NoSniff**: Prevents MIME-sniffing to mitigate content type risks.
- **ReferrerPolicy**: Controls referrer information to protect user privacy.
- **XssProtection**: Blocks XSS attacks by enabling the browser's XSS filter.
- **Custom Headers**: Allows setting custom response headers for security.
- **Health Check**: Provides a simple health check endpoint.
- **CORS**: Implements CORS for controlled resource access.
- **Logger**: Logs HTTP request details for monitoring and analysis.
- **Recovery**: Recovers from panics to maintain server stability.

### Build and Management
- **Build**: `make build`
- **Run**: `make run`
- **Docker**: `make docker-run, make docker-down`
- **Test**: `make test`
- **Clean**: `make clean`
- **Certs**: `make gen-cert, make clear-cert`
- **Live Reload**: `make watch`