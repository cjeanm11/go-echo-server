## Secure Server Template

My personal server setup based on go-blueprints 

### Features
- Minimalist Go web framework: [Echo](https://github.com/labstack/echo)
- Database: [PostgreSQL](https://www.postgresql.org/)

### Middleware Overview

1. **Content-Security-Policy**: Restricts content sources to mitigate XSS attacks.
2. **Default Headers**: Enhances security with default response headers.
3. **HSTS**: Enforces HTTPS usage to prevent protocol downgrade attacks.
4. **NoSniff**: Prevents MIME-sniffing to mitigate content type risks.
5. **ReferrerPolicy**: Controls referrer information to protect user privacy.
6. **XssProtection**: Blocks XSS attacks by enabling the browser's XSS filter.
7. **Custom Headers**: Allows setting custom response headers for security.
8. **Health Check**: Provides a simple health check endpoint.
9. **CORS**: Implements CORS for controlled resource access.
10. **Logger**: Logs HTTP request details for monitoring and analysis.
11. **Recovery**: Recovers from panics to maintain server stability.

### Build Commands

**Build the Application** : Compile the server application.
```bash
make build
```

**Run the Application** : build the application
```bash
make run
```

Live Reload: Automatically rebuild and restart the server upon file changes.
```bash
make watch
```

### Database Management

**Create DB Container**: Launch a Docker container running PostgreSQL.
```bash
make docker-run
```

**Shutdown DB Container**: Stop and remove the PostgreSQL container.
```bash
make docker-down
```

### Testing and Clean-up
**Run Tests**: Execute the test suite to ensure functionality.

```bash
make test
```

**Clean Up**: Remove binary files generated from previous builds.
```bash
make clean
```
