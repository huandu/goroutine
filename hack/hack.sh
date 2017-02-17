#! /bin/sh

CUR_DIR=$(dirname $0)
GO_SRC="$1"

if [ -z "$GO_SRC" ] || [ "$GO_SRC" = "-h" ] || [ "$GO_SRC" = "--help" ]; then
    cat <<EOF
Build the hack tool and run it to hack go src.

Usage: $0 path-to-go-src
EOF
    exit 1
fi

cd "$CUR_DIR"

if ! go build -o hack; then
    exit 1
fi

./hack -go-src="$GO_SRC" -import-path=github.com/huandu/goroutine/hack -debug
