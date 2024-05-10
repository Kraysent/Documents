### Docarchive

This is a backend and frontend of a simple document storage. Notable features include but are not limited to:
- Go backend, including [scs](https://github.com/alexedwards/scs) for a session management after the authentication through Google OAuth, [chi](https://pkg.go.dev/github.com/go-chi/chi@v1.5.5) router. [Golang-migrate](https://pkg.go.dev/github.com/golang-migrate/migrate/v4@v4.16.2) and [embedded-postgres](https://pkg.go.dev/github.com/fergusstrange/embedded-postgres@v1.25.0) are used for testing.
- Prometheus metrics that are exported directly from the application to Grafana.
- Terraform configuration that sets up Yandex Cloud Postgres storage, networks and docker registry.
- CI/CD pipeline for testing, building the image and uploading it to Yandex Cloud docker registry.
- Simple NGINX configuration that routes traffic from client to the frontend or the backend depending on the path.
- Simple frontend written in Typescript with React framework. This frontend is mostly read-only.
