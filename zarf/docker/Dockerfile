FROM golang:1.17 as build_wt
ENV CGO_ENABLED 0
ARG BUILD_REF

COPY . /service

WORKDIR /service/app
RUN go build -ldflags "-X main.build=${BUILD_REF}"

FROM alpine:3.14
ARG BUILD_DATE
ARG BUILD_REF
COPY --from=build_wt /service/app/app /service/app
WORKDIR /service
CMD ["./service"]