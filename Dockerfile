# syntax=docker/dockerfile:1

# Imágenes base fijadas por digest: la etiqueta dice la versión, el digest garantiza el contenido.
# Para actualizarlas: docker buildx imagetools inspect <imagen:etiqueta>

# --- 1. CSS: Tailwind lee los .templ para saber qué clases generar ---
FROM --platform=$BUILDPLATFORM node:22-alpine@sha256:fe4f6d74e20370acb04ad889077d14d0f64a9c8ea6a42bc60cddf1e844699914 AS css
WORKDIR /src
COPY package.json package-lock.json ./
# --ignore-scripts: ninguna dependencia ejecuta código al instalarse
RUN --mount=type=cache,target=/root/.npm npm ci --ignore-scripts --no-audit --no-fund
COPY static/css/input.css static/css/input.css
COPY internal/views/ internal/views/
RUN ./node_modules/.bin/tailwindcss -i static/css/input.css -o static/css/output.css --minify

# --- 2. Go: genera templ y compila un binario estático ---
FROM --platform=$BUILDPLATFORM golang:1.27.3-bookworm@sha256:972b4110ba0f47539a975c31e52c64765e1ff19f2d74469ae4a595c6676e80a5 AS build
WORKDIR /src
ENV CGO_ENABLED=0 GOTOOLCHAIN=local
# Debe coincidir con la versión de github.com/a-h/templ en go.mod
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/a-h/templ/cmd/templ@v0.3.1020
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
# El CSS tiene que estar antes de go build: static/ se embebe al compilar
COPY --from=css /src/static/css/output.css static/css/output.css
RUN templ generate
ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/ratelog ./cmd/app

# --- 3. Imagen final: solo el binario (más certificados CA, tzdata y usuario nonroot) ---
FROM gcr.io/distroless/static-debian13:nonroot@sha256:18f4cc26c9ffa7f20a359298be7deadac5c607c8e03a02b8119f1a5fc2f82c99
LABEL org.opencontainers.image.title="RateLog" \
      org.opencontainers.image.source="https://github.com/luism2302/RateLog"
COPY --from=build /out/ratelog /ratelog
USER nonroot:nonroot
EXPOSE 8000
# Comprueba el puerto por defecto (8000)
HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=3 CMD ["/ratelog", "-healthcheck"]
ENTRYPOINT ["/ratelog"]
