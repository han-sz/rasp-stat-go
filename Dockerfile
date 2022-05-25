FROM golang

WORKDIR /app

# Default: 4322
ARG PORT 
# Default: 10 (data points)
ENV BUFFER
# Default: 6 (seconds)
ENV INTERVAL
ENV GIN_MODE=release

COPY Makefile
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY rasp-stat ./
RUN make prod-arm

EXPOSE $PORT

ENTRYPOINT "build/rasp-stat_arm-linux"