FROM alpine AS prerequisite

###

FROM golang:1.15.5 AS base

WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download

###

FROM base AS build

COPY . .
RUN make build

###

FROM prerequisite

COPY --from=build /src/bin/api /

ENTRYPOINT ["/api"]
