# Hack Go Runtime #

This is a hack tool to copy const and type declarations from go source and inject hacked code.
With this hack, it's possible to expose any interesting type to wild world.

## Play with it ##

Download a go src tar file from https://golang.org/dl/, extract all files to a folder and
run following command.

    ./hack.sh path-to-go-src

## How it works ##

The basic idea is quite simple: If we can get internal type declarations and make it public,
then we can use it easily. So what I do is to parse all interesting go source files and copy
types to a folder and inject some hacked code to make some types public.

The `generator.go` is the key file to make it happen. It loads all hacks and walk through
all AST tree nodes in a file to find out all constants and types. It works quite well with
internal package and ignored files, as these features are required for go1.6 and greater.

There is only 1 hack implemented in `runtime_hacker.go`. If we want to add more hacks, just
implement more hacker files.
