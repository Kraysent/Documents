### Docarchive

This is a backend and frontend of a simple document storage. Notable features include but are not limited to:
- Go backend, including [scs](https://github.com/alexedwards/scs) for a session management after the authentication through Google OAuth, [chi](https://pkg.go.dev/github.com/go-chi/chi@v1.5.5) router. [Golang-migrate](https://pkg.go.dev/github.com/golang-migrate/migrate/v4@v4.16.2) and [embedded-postgres](https://pkg.go.dev/github.com/fergusstrange/embedded-postgres@v1.25.0) are used for testing.
- Prometheus metrics that are exported directly from the application to Grafana.
- Terraform configuration to set up Yandex Cloud Postgres storage, networks and docker registry.
- CI/CD pipeline for testing, building the image and upload it to Yandex Cloud docker registry.
- Simple NGINX configuration to route traffic from client to frontend or backend depending on the path.
- Simple frontend written in Typescript with React framework. This frontend is mostly read-only.
