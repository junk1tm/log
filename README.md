# log

[![ci-img]][ci]
[![docs-img]][docs]
[![report-img]][report]
[![codecov-img]][codecov]

Simple, structured and opinionated logging interface with levels.

## Rationale

According to the [Dependency inversion principle][dip], applications should depend on abstractions rather than specific
implementations, thus avoiding unnecessary coupling. And logger, being just another dependency, shouldn't be used
directly but should be hidden behind an interface.

Unfortunately, Go's [log package][std-log] doesn't provide any logging interface (see related
issues [#13182][issue-13182], [#28412][issue-28412]). Standard [log.Logger][std-logger] is implemented as a type
with `Print[f|ln]` methods, and it's usually enough for small applications and tools. However, when it comes to more
complex systems, this implementation lacks some important features including levels and structured logging. So, most
applications use third-party loggers, such as [zap][zap], [logrus][logrus] or [zerolog][zerolog].

This package doesn't aim to provide a universal logging API (I doubt designing such a thing is possible). Instead, it
combines the key features of the most popular third-party logging libraries, levels and structured logging, into an
opinionated interface. For simple applications that don't require complex logging, you might
prefer [io.Writer][io-writer].

## Features

* Simple API
* Dependency-free (implementations are optional)
* Support for most basic types
* Support for user-defined types implementing [Loggable][loggable]
* Support for [child loggers][with-fields]
* Support for [hooks][with-hooks]
* Implementations for the most popular logging libraries:
  * [zap][zap-impl]
  * [logrus][logrus-impl]
  * [zerolog][zerolog-impl]

## Install

```bash
go get github.com/junk1tm/log
```

## Interface

```go
// Logger is a minimal logging interface with levels.
// It uses structured logging with Field parameters to ensure type safety.
type Logger interface {
    // Debug writes logs at DEBUG level.
    // It is used to log information useful for developers.
    Debug(msg string, fields ...Field)
    // Info writes logs at INFO level.
    // It is used to log information useful for users.
    Info(msg string, fields ...Field)
    // Error writes logs at ERROR level.
    // It is used to handle errors by logging them.
    Error(msg string, fields ...Field)
}
```

where Field is either a builtin type or an implementation of the following interface:

```go
// Loggable represents an entity that can be logged.
// It provides a type-safe way to add any user-defined type to log context.
// Use Object function to pass Loggable as a Field to Logger's methods.
type Loggable interface {
    // ToLog returns a list of fields to log.
    ToLog() []Field
}
```

## Example

```go
package main

import (
	"github.com/example/logger" // logging backend of your choice (e.g., zap)
	"github.com/junk1tm/log"
	"github.com/junk1tm/log/exampleimpl"
)

func main() {
	l := logger.New()
	a := app{logger: exampleimpl.NewLogger(l)}
	a.run()
}

type app struct {
	logger log.Logger // do not depend on github.com/example/logger directly
}

func (a *app) run() {
	a.logger.Info("running app", log.String("some", "context"))
}
```

## FAQ

### No `warn` level?

I agree with [Dave Cheney's thoughts][cheney-post] on this one.

### No termination levels (`fatal`/`panic`)?

I don't think it's logger's responsibility to terminate a program. To avoid writing something like

```go
if err != nil {
    logger.Error("could not do something", log.Error(err))
    panic(err) // or os.Exit(1)
}
```

on each error, I recommend following [exit-once pattern][exit-once]. If you really want to use `Fatal`/`Panic`
methods, the original unwrapped logger is available in `main.go` (see [example](#Example)), where all the termination
happens, so you might just use its methods (most loggers provide them).

### How to deal with repetitive keys?

If you find yourself using repetitive keys a lot, you might as well define a helper package with your domain-specific
fields, for example:

```go
package fields

import "github.com/junk1tm/log"

func UserID(id int) log.Field { return log.Int("user_id", id) }
```

### Why add hooks if most loggers already support them?

For the same reason the logging interface is introduced in the first place: to prevent coupling between your code
(hooks, in this case) and a logging library of your choice.

[ci]: https://github.com/junk1tm/log/actions/workflows/go.yml
[ci-img]: https://github.com/junk1tm/log/actions/workflows/go.yml/badge.svg
[docs]: https://pkg.go.dev/github.com/junk1tm/log
[docs-img]: https://pkg.go.dev/badge/github.com/junk1tm/log.svg
[report]: https://goreportcard.com/report/github.com/junk1tm/log
[report-img]: https://goreportcard.com/badge/github.com/junk1tm/log
[codecov]: https://codecov.io/gh/junk1tm/log
[codecov-img]: https://codecov.io/gh/junk1tm/log/branch/main/graph/badge.svg
[dip]: https://en.wikipedia.org/wiki/Dependency_inversion_principle
[std-log]: https://pkg.go.dev/log
[std-logger]: https://pkg.go.dev/log#Logger
[issue-13182]: https://github.com/golang/go/issues/13182
[issue-28412]: https://github.com/golang/go/issues/28412
[zap]: https://github.com/uber-go/zap
[logrus]: https://github.com/sirupsen/logrus
[zerolog]: https://github.com/rs/zerolog
[io-writer]: https://pkg.go.dev/io#Writer
[loggable]: https://pkg.go.dev/github.com/junk1tm/log#Loggable
[with-fields]: https://pkg.go.dev/github.com/junk1tm/log#WithFields
[with-hooks]: https://pkg.go.dev/github.com/junk1tm/log#WithHooks
[zap-impl]: https://pkg.go.dev/github.com/junk1tm/log/zapimpl
[logrus-impl]: https://pkg.go.dev/github.com/junk1tm/log/logrusimpl
[zerolog-impl]: https://pkg.go.dev/github.com/junk1tm/log/zerologimpl
[cheney-post]: https://dave.cheney.net/2015/11/05/lets-talk-about-logging
[exit-once]: https://github.com/uber-go/guide/blob/master/style.md#exit-once
