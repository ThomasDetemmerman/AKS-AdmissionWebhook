# Build stage
FROM golang:1.16 AS build
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
WORKDIR /work
COPY . /work
RUN --mount=type=cache,target=/root/.cache/go-build,sharing=private
RUN go build -o bin/admission-webhook .

# Run stage
FROM alpine:latest  AS run
COPY --from=build /work/bin/admission-webhook /usr/local/bin/
CMD ["admission-webhook"]
