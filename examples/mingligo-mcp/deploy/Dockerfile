# 多阶段构建：镜像仅含 mingligo-mcp 运行时。
FROM golang:1.26.11-alpine@sha256:ae52b5bc6ed045a9ed26ab03dc19327a802a9cc9aa8b27ed57f5c6b60579b0bb AS build
WORKDIR /src
ENV GOPROXY=https://goproxy.cn,direct
ARG GIT_SHA=unknown
ARG COMMIT=unknown
ARG BUILD_TIME=unknown
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" \
    go build -ldflags "-X main.BuildID=${GIT_SHA}" -o /out/mingligo-mcp ./examples/mingligo-mcp/

FROM alpine:3.21@sha256:48b0309ca019d89d40f670aa1bc06e426dc0931948452e8491e3d65087abc07d AS certs
RUN apk add --no-cache ca-certificates
RUN adduser -D -u 1000 mingligo

FROM scratch
ARG GIT_SHA=unknown
ARG COMMIT=unknown
ARG BUILD_TIME=unknown
LABEL org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.title="mingligo-mcp" \
      org.opencontainers.image.description="MingliGo deterministic metaphysics fact library"
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=certs /etc/passwd /etc/passwd
COPY --from=build /out/mingligo-mcp /usr/local/bin/mingligo-mcp
COPY --from=build /src/pkg/city/ATTRIBUTION.txt /usr/share/licenses/mingligo-mcp/ATTRIBUTION.txt
USER mingligo
EXPOSE 8081
ENTRYPOINT ["/usr/local/bin/mingligo-mcp", "--addr", ":8081"]
