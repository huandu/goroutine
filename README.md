# Hacking goroutine #

[![Build Status](https://travis-ci.org/huandu/goroutine.png?branch=master)](https://travis-ci.org/huandu/goroutine)

Package goroutine is merely a hack.
It exports goroutine id to outside so that you can use it for whatever purpose.
However, it's highly recommended to not use this package in your daily life.
It may be broken at any go release as it's a hack.

## Usage ##

Get the latest version through `go get -u github.com/huandu/goroutine`.

Get current goroutine id with `goroutine.GoroutineId()`.

```go
// Get id of current goroutine.
var id int64 = goroutine.GoroutineId()
println(id)
```

See [godoc](https://godoc.org/github.com/huandu/goroutine) for more details.

## Supported builds ##

Package goroutine is not well tested due to lack of test machines.
Ideally, it should work on all go >= go1.5.

Tested platforms.
* Darwin (Mac OSX 10.11.6) + amd64 CPU
    * go1.5.1
    * go1.6.3
    * go1.7
    * go1.7.1
    * go1.7.3
* Travis CI (See https://travis-ci.org/huandu/goroutine)
    * go1.5
    * go1.5.1
    * go1.5.2
    * go1.5.3
    * go1.5.4
    * go1.6
    * go1.6.1
    * go1.6.2
    * go1.6.3
    * go1.7
    * go1.7.1
    * go1.7.2
    * go1.7.3

## How it works ##

Go program is a statically built binary. The runtime inside a binary is statically linked. It means, if I know Go version
and runtime source code for this version, I can copy struct declaration from runtime package source to my package and
cast runtime internal pointers to its underlying struct safely. As Go is an open source project, I can always find the right
struct for an interesting runtime pointer and then manipulate it.

In this package, I just get current goroutine pointer (copy the `getg()` implementation from compiler) and cast it to a right
struct. It sounds simple. However, the struct `g` refers many other internal types defined in `runtime` package. I cannot
simply copy some necessary types to my package to make it work. I have to scan all types and constants in `runtime` and its
internal packages to make the struct `g` well defined. Another challenge is that Go authors update `runtime` structs in nearly
every major version (or even in a minor version). I have to maintain hacked code for every Go release. I develop a
semi-automatical tool to make things easier. I guess I may need to think of other better way to avoid to generate hacked source
for every Go release.

NOTE: Starting from go1.7.2, Go compiler generates some constants definition for `runtime` package according to build flags and
environment when building. It makes current hack impossible to handle all posibile flag and environment combinations. I make a
hack to detour it and no impact to the major task of this package - get goroutine id. However, it's not a perfect solution.
If I want to do more tricks in runtime, such constants may bother me.

I'm think of a perfect solution. If you have any suggestion, please open issue and let me know. Many thanks.

## License ##

This package is licensed under MIT license. See LICENSE for details.
