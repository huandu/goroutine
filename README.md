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

## Caveats ##

Package goroutine is not well tested due to lack of test machines.
Ideally, it should work on all go >= go1.5.

Tested platforms.
* go1.5.1 + Darwin (Mac OSX 10.11.6) + amd64 CPU
* go1.6.3 + Darwin (Mac OSX 10.11.6) + amd64 CPU
* go1.7 + Darwin (Mac OSX 10.11.6) + amd64 CPU
* go1.7.1 + Darwin (Mac OSX 10.11.6) + amd64 CPU

## License ##

This package is licensed under MIT license. See LICENSE for details.
