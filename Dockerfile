# syntax=docker/dockerfile:1

# Base images pinned by digest: the tag shows the version, the digest guarantees the content.
# To update them: docker buildx imagetools inspect <image:tag>

# --- 1. CSS: Tailwind reads the .templ files to know which classes to generate ---
FROM --platform=$BUILDPLATFORM node:22-alpine@sha256:fe4f6d74e20370acb04ad889077d14d0f64a9c8ea6a42bc60cddf1e844699914 AS css
WORKDIR /src
COPY package.json package-lock.json ./
# --ignore-scripts: no dependency runs code on install
RUN --mount=type=cache,target=/root/.npm npm ci --ignore-scripts --no-audit --no-fund
COPY static/css/input.css static/css/input.css
COPY internal/views/ internal/views/
RUN ./node_modules/.bin/tailwindcss -i static/css/input.css -o static/css/output.css --minify

# --- 2. Go: generates templ code and builds a static binary ---
FROM --platform=$BUILDPLATFORM golang:1.27.3-bookworm@sha256:972b4110ba0f47539a975c31e52c64765e1ff19f2d74469ae4a595c6676e80a5 AS build
WORKDIR /src
ENV CGO_ENABLED=0 GOTOOLCHAIN=local
# Must match the github.com/a-h/templ version in go.mod
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/a-h/templ/cmd/templ@v0.3.1020
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
# The CSS must exist before go build: static/ is embedded at compile time
COPY --from=css /src/static/css/output.css static/css/output.css
RUN templ generate
ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/ratelog ./cmd/app

# --- 3. Final image: just the binary (plus CA certificates, tzdata and a nonroot user) ---
FROM gcr.io/distroless/static-debian13:nonroot@sha256:18f4cc26c9ffa7f20a359298be7deadac5c607c8e03a02b8119f1a5fc2f82c99
LABEL org.opencontainers.image.title="RateLog" \
      org.opencontainers.image.source="https://github.com/luism2302/RateLog"
COPY --from=build /out/ratelog /ratelog
USER nonroot:nonroot
EXPOSE 8000
# Checks the default port (8000)
HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=3 CMD ["/ratelog", "-healthcheck"]
ENTRYPOINT ["/ratelog"]
