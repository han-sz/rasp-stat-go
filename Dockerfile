FROM golang

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
ENV GIN_MODE=release
RUN go build -o rasp-stat-go

EXPOSE 4322

ENTRYPOINT "rasp-stat-go"