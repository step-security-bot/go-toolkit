FROM golang:1-bookworm as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Download the example publications
# RUN git lfs pull && ls -alh publications

# Run tests
FROM builder AS tester
RUN go test ./...

# Run goreleaser
RUN goreleaser build --single-target --id rwp --skip-validate --snapshot

# Produces very small images
FROM gcr.io/distroless/static-debian12 AS packager

# Add Fedora's mimetypes (pretty up-to-date and expansive)
# since the distroless container doesn't have any. Go uses
# this file as part of its mime package, and readium/go-toolkit
# has a mediatype package that falls back to Go's mime
# package to discover a file's mimetype when all else fails.
ADD https://pagure.io/mailcap/raw/master/f/mime.types /etc/

# Add two demo EPUBs to the container by default
ADD https://readium-playground-files.storage.googleapis.com/demo/moby-dick.epub /srv/publications/
ADD https://readium-playground-files.storage.googleapis.com/demo/BellaOriginal3.epub /srv/publications/

# Copy built Go binary
COPY --from=builder /app/dist/rwp_linux_amd64_v3/rwp /opt/

EXPOSE 15080

USER nonroot:nonroot

ENTRYPOINT ["/opt/rwp"]
CMD ["serve", "/srv/publications"]