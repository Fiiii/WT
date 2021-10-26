# Build the Go Binary.
FROM golang:1.17 as build_wt-api
ENV CGO_ENABLED 0

COPY . /wt

# Build the service binary.
WORKDIR /wt/app
RUN go build

# Run the Go Binary in Alpine.
FROM alpine
COPY --from=build_wt-api /wt/app/app /wt/app
WORKDIR /wt

CMD ["./app"]