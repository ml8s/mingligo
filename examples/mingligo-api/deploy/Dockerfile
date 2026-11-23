FROM golang:1.26.11-alpine@sha256:ae52b5bc6ed045a9ed26ab03dc19327a802a9cc9aa8b27ed57f5c6b60579b0bb AS build
WORKDIR /src
ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}
ARG GIT_SHA=unknown
ARG COMMIT=unknown
ARG BUILD_TIME=unknown
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" \
    go build -ldflags "-X main.BuildID=${GIT_SHA}" -o /out/mingligo-api ./examples/mingligo-api

FROM alpine:3.21@sha256:48b0309ca019d89d40f670aa1bc06e426dc0931948452e8491e3d65087abc07d AS certs
ARG APK_MIRROR=https://mirrors.tuna.tsinghua.edu.cn
RUN sed -i "s#https://dl-cdn.alpinelinux.org#${APK_MIRROR}#g" /etc/apk/repositories \
 && apk add --no-cache ca-certificates
RUN adduser -D -u 1000 mingligo

FROM scratch
ARG GIT_SHA=unknown
ARG COMMIT=unknown
ARG BUILD_TIME=unknown
LABEL org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.title="mingligo-api" \
      org.opencontainers.image.description="MingliGo deterministic metaphysics fact API"
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=certs /etc/passwd /etc/passwd
COPY --from=build /out/mingligo-api /usr/local/bin/mingligo-api
COPY --from=build /src/examples/mingligo-api/openapi.json /usr/local/bin/openapi.json
COPY --from=build /src/pkg/city/ATTRIBUTION.txt /usr/share/licenses/mingligo-api/ATTRIBUTION.txt
USER mingligo
EXPOSE 18082
ENTRYPOINT ["/usr/local/bin/mingligo-api"]
