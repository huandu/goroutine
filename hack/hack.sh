#! /bin/sh

CUR_DIR=$(dirname $0)
GO_SRC="$1"

cd "$CUR_DIR"

if ! go build -o hack; then
    exit 1
fi

./hack -go-src="$GO_SRC" -import-path=github.com/huandu/goroutine/hack
