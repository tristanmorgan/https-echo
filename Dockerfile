#
# Builder
#
FROM golang:alpine AS builder

WORKDIR /go/src/github.com/tristanmorgan/https-echo/

COPY . .

ARG LD_FLAGS="-s -w"
ENV LD_FLAGS="${LD_FLAGS}"

RUN CGO_ENABLED="0" go build -v -a -trimpath -o "/https-echo" -ldflags "${LD_FLAGS}"

#
# Final
#
FROM scratch
LABEL maintainer="Tristan Morgan <tristan.morgan@hashicorp.com>"
LABEL Description="HTTPS_ECHO, echo url and redirect http to https"
EXPOSE 80
WORKDIR /
COPY --from=builder /https-echo /https-echo
ENTRYPOINT ["/https-echo"]
