# Logfmt Encoder

This package provides a logfmt encoder for [zap][zap].

## Usage

The encoder is easy to configure. Simply create a new core with an instance of the logfmt encoder and use it with your preferred logging interface.

```go
package main

import (
    "os"

    "github.com/allir/zap-logfmt"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func main() {
    config := zap.NewProductionEncoderConfig()
    logger := zap.New(zapcore.NewCore(
        zaplogfmt.NewEncoder(config),
        zapcore.Lock(os.Stdout),
        zapcore.DebugLevel,
    ))
    defer logger.Sync()

    logger.Info("Hello World")
}
```

To use RFC3339 output for the time instead of an integer timestamp, you provide EncodeTime to the EncoderConfig:

```go
package main

import (
    "os"

    "github.com/allir/zap-logfmt"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func main() {
    config := zap.NewProductionEncoderConfig()
    config.EncodeTime = zapcore.RFC3339TimeEncoder
    logger := zap.New(zapcore.NewCore(
        zaplogfmt.NewEncoder(config),
        zapcore.Lock(os.Stdout),
        zapcore.DebugLevel,
    ))
    defer logger.Sync()

    logger.Info("Hello World")
}
```

An alternative way to set up the logger by registering the encoder and using it with the config builder. Also setting the time encoding to RFC3339.

```go
package main

import (
    zaplogfmt "github.com/allir/zap-logfmt"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func main() {
    zaplogfmt.Register()
    zapConfig := zap.NewProductionConfig()
    zapConfig.EncoderConfig = zap.NewProductionEncoderConfig()
    zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
    zapConfig.Encoding = "logfmt"

    logger, err := zapConfig.Build()
    if err != nil {
        panic(err)
    }
    defer logger.Sync()

    logger.Info("Hello World")
}
```

## Limitations

It is not possible to log an array, channel, function, map, slice, or
struct. Functions and channels since they don't really have a suitable
representation to begin with. Logfmt does not have a method of
outputting arrays or maps so arrays, slices, maps, and structs cannot be
rendered.

## Namespaces

Namespaces are supported. If a namespace is opened, all of the keys will
be prepended with the namespace name. For example, with the namespace
`foo` and the key `bar`, you would get a key of `foo.bar`.

## Attribution

This is a fork of the original encoder from [github.com/jsternberg/zap-logfmt][jsternberg]. And pulling in and combining additional fixes from other sources such as;

* [github.com/jdechicchis/zap-logfmt][jdechicchis]
* [github.com/sykesm/zap-logfmt][sykesm]
* [github.com/indra-kargo/zap-logfmt][indra-kargo]

[zap]: https://github.com/uber-go/zap
[jsternberg]: https://github.com/jsternberg/zap-logfmt
[jdechicchis]: https://github.com/jdechicchis/zap-logfmt
[sykesm]: https://github.com/sykesm/zap-logfmt
[indra-kargo]: https://github.com/indra-kargo/zap-logfmt
