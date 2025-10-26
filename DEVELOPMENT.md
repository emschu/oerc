# OERC Development

This project is shipped with a `Makefile` to ease the development and testing process.
At first, you should run `make setup` and you need the usual Golang/Python/Node/Java toolchains.
Be sure to run `make build` (if you updated something there) before filing a pull request.

### Build/Configuration Instructions

#### Prerequisites
- Go 1.24+
- PostgreSQL 12+
- Node.js (for frontend development)

#### Setup
1. Clone the repository
2. Run `make setup` to install required dependencies:
   ```bash
   make setup
   ```
   This will:
    - Download Go dependencies
    - Install required tools (go.rice, revive)
    - Install schemathesis for API testing

3. Configure the database:
    - For development, you can use a PostgreSQL Docker container:
      ```bash
      docker run --name oer-postgres -p 5432:5432 -e POSTGRES_PASSWORD=root -e POSTGRES_DB=oer_server_dev -d postgres:14-alpine
      ```
    - Create a configuration file at `~/.oerc.yaml` or use `-c <path>` to specify a custom location
    - Ensure database connection details are correctly set in the configuration file

#### Building
- For development build:
  ```bash
  make build
  ```
  This creates a binary at `bin/oerc` with race detection enabled

- For release builds:
  ```bash
  make release
  ```
  This creates binaries for multiple platforms in the `bin/` directory

- For frontend development:
  ```bash
  make frontend
  ```
  This builds the frontend and generates the static rice box file

### Testing Information

#### Running Tests
- Run all tests:
  ```bash
  make test
  ```
  This runs unit tests with race detection and tracing

- Run tests with coverage:
  ```bash
  make cover
  ```
  This generates a coverage report at `cover.html`

- Run integration tests:
  ```bash
  make integration-test-prepare
  make integration-test
  ```
  This starts a local server and runs OpenAPI schema conformity tests against it

### Project Structure
- Main application code is in the root directory
- Frontend client code is in the `client/` directory
- API specifications are in the `spec/` directory
- Configuration templates are in the `config/` directory

### Database
- The project uses GORM as an ORM
- Models are defined in `models.go`
- For testing, an in-memory SQLite database is used

#### Development PostgreSQL Container
```bash
docker run --name oer-postgres -p 5432:5432 -e POSTGRES_PASSWORD=root -e POSTGRES_DB=oer_server_dev -d postgres:14-alpine 
```