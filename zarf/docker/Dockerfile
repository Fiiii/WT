FROM golang:1.17 as build_wt
ENV CGO_ENABLED 0
ARG BUILD_REF

COPY . /service

# Build the admin binary.
WORKDIR /service/app/tooling/admin
RUN go build -ldflags "-X main.build=${BUILD_REF}"

WORKDIR /service/app/services/wt-api
RUN go build -ldflags "-X main.build=${BUILD_REF}"

FROM alpine:3.14
ARG BUILD_DATE
ARG BUILD_REF
COPY --from=build_wt /service/zarf/keys/. /service/zarf/keys/.
COPY --from=build_wt /service/app/tooling/admin/admin /service/admin
COPY --from=build_wt /service/app/services/wt-api/wt-api /service/wt-api
WORKDIR /service
CMD ["./wt-api"]

