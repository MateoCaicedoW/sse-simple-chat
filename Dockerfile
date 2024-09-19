FROM golang:1.22-alpine AS builder
RUN apk --update add build-base

WORKDIR /src/app
ADD go.mod .
RUN go mod download

ADD . .

# Copy everything in assets to the public folder
RUN cp internal/assets/* public/

# Building TailwindCSS with tailo
RUN go run github.com/paganotoni/tailo/cmd/build@v1.0.8

# Building the migrate command
RUN go build -tags osusergo,netgo -buildvcs=false -o bin/migrate ./cmd/migrate

# Building the app
RUN go build -tags osusergo,netgo -buildvcs=false -o bin/app ./cmd/app

FROM alpine
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /bin/

# Copying binaries
COPY --from=builder /src/app/bin/app .
COPY --from=builder /src/app/bin/migrate .

# Specifying the shell to use
SHELL ["/bin/ash", "-c"]
CMD migrate && app