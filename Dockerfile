FROM golang:1.15.3-alpine AS builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN apk add git

# Set the Current Working Directory inside the container
WORKDIR /src

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

COPY ./data/config.dev.yml ./data/config.yml

EXPOSE 8080

CMD ["go", "run", "main.go"]

# Build the Go app
RUN GIT_COMMIT=$(git rev-list -1 HEAD) && \
    go build -ldflags="-X 'sutjin/go-rest-template/internal/pkg/config.Version=$GIT_COMMIT' -X 'git.infokes.id/ehealthcare/iam/internal/routing-service/config.BuildTime=$(date)'" \
    -o ./out/app ./main.go
    
# Start fresh from a smaller image
FROM alpine:3.12
RUN apk add ca-certificates

WORKDIR /app

RUN mkdir log

COPY --from=builder /src/out/app /app/restapi
COPY --from=builder /src/data /app/data

COPY ./data/config.prod.yml ./data/config.yml

RUN chmod +x restapi

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
ENTRYPOINT ./restapi
