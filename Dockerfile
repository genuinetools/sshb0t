# Builder image
FROM golang:alpine AS builder
WORKDIR /go/src/github.com/jessfraz/sshb0t
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o sshb0t

# Final image
FROM alpine
LABEL maintainer "Jessica Frazelle <jess@linux.com>"
RUN	apk add --no-cache \
    ca-certificates \
    git
COPY --from=builder /go/src/github.com/jessfraz/sshb0t/sshb0t /usr/bin/sshb0t
ENTRYPOINT ["sshb0t"]
