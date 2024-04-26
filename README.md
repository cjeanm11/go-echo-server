## Server Template

My personal server setup based on go-blueprints 

### Features
- [Echo](https://github.com/labstack/echo): minimalist Go web framework.
- [PostgreSQL](https://www.postgresql.org/): Open source object-relational database system.

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
