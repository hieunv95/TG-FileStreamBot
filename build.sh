#!/bin/sh
CGO_ENABLED=0 go build -o /app/fsb -ldflags="-w -s" ./cmd/fsb
