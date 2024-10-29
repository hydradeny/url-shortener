FROM golang:1.22 AS builder

WORKDIR /usr/local/go/src/


# add private repository through ssh
RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN git config --system --add url.ssh://git@github.com:.insteadOf https://github.com
RUN go env -w GOPRIVATE="github.com/hydradeny/*"

ADD ./go.* /usr/local/go/src/
RUN go clean --modcache

RUN --mount=type=ssh go mod download

ADD . /usr/local/go/src/
RUN CGO_ENABLED=0 go build -mod=readonly -o app cmd/auth/main.go

FROM alpine:3.19

COPY --from=builder /usr/local/go/src/app /
CMD ["/app"]
