FROM golang:1.17-alpine as build
WORKDIR /ocp-src
COPY . .
RUN go mod vendor && \
        CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /ocp cmd/ocp-metrics/main.go && \
        apk add upx binutils && \
        strip /ocp && \
        upx /ocp && \
        ls -alh /ocp

FROM scratch
LABEL org.opencontainers.image.source https://github.com/rtgnx/octopus
ENTRYPOINT ["/ocp"]
COPY --from=build /ocp /ocp
