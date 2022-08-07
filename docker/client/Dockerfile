FROM golang:1.16-alpine
ARG ENV
ARG SERVICE
ARG PATH
RUN apk add --no-cache git
ENV SERVICE ${SERVICE}
# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .


# Build the Go app
RUN go build -o ./${SERVICE} ./cmd/${SERVICE}


# This container exposes port 8080 to the outside world
EXPOSE 8080

CMD ["sh", "-c", "./${SERVICE} --addr=server"]