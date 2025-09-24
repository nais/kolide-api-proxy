# Kolide API proxy

Proxy for a small selection of the Kolide API.

## Local Development

### Clone the repository:

```bash
git clone git@github.com:nais/kolide-api-proxy.git
cd kolide-api-proxy
```

### Install tools

Install tools using [mise](https://mise.jdx.dev/):

```bash
mise install
```

### Copy environment variables

Create a copy of the example environment file and adjust as needed:

```bash
cp .env.example .env
```

### Run the application

Run the application to start the HTTP server:

```bash
go run main.go
```

Point your HTTP client of choice to `http://localhost:8080/api/devices` (or the `HTTP_LISTEN_ADDRESS` you have specified in your `.env` file), authenticate using `username:password` (or the `PROXY_USERNAME:PROXY_PASSWORD` combination you have specified in your `.env` file) to view devices.

### Run tasks

Run `mise run` to see all available tasks. These are some of the most common ones used for local development:

```bash
mise run check # run all static code analysis tools
```

## License

MIT, see [LICENSE.txt](./LICENSE.txt) for details.
