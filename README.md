# Prerequisites

This project only assumes a modern-ish version of Docker (with Compose).

# Running The Project

Run/recompile the project and view logs:

```shell
make up 
```


# The Project Layout
```sh
cmd/upload-service/
└── main.go                        # entrypoint

internal/app/upload-service/
├── api
│   ├── api.go                     # defines routes and includes middlewares
│   ├── handler.go
│   └── openapi.yaml               # defines openAPIs routes
├── models
│   ├── file.go
└── repository
    └── repository.go

```
<!-- 
your_project/
├── vendor/               # Contains GORM and PostgreSQL driver and third party libraries
│   ├── gorm.io/
│   ├── github.com/
│   └── ...
├── go.mod                # Dependencies declared here
├── go.sum                # Dependency checksums
├── main.go               # Your application code -->
