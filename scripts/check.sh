#!/usr/bin/env sh
set -eu

unformatted="$(find . -name '*.go' \
  -not -path './.git/*' \
  -not -path './.gocache/*' \
  -not -path './subscriptions_tracker/*' \
  -exec gofmt -l {} +)"
if [ -n "$unformatted" ]; then
  echo "The following files need gofmt:"
  echo "$unformatted"
  exit 1
fi

go vet ./...
go test ./...
