# Logfmt Encoder

This package provides a logfmt encoder for [zap][zap].

[![GoDoc](https://godoc.org/github.com/sykesm/zap-logfmt?status.svg)](https://godoc.org/github.com/sykesm/zap-logfmt)

## Usage

The encoder is easy to configure. Simply create a new core with an instance of the logfmt encoder and use it with your preferred logging interface.

```go
package main

import (
    "os"

    "github.com/sykesm/zap-logfmt"
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

## Arrays, Objects, and Reflected Fields

While it's best to avoid complex data types in log fields, there are times when they sneak in. When complex fields are included in log records, they will be encoded, but they won't be very pretty.

### Arrays

Arrays are encoded as a comma separated list of values within square brackets. This format is very similar to JSON encoding. Arrays of simple scalars remain quite readable but including elements that require quoting will result in very ugly records.

### Objects

Objects are encoded as a space separated list of _key=value_ pairs. Because this format includes an equals sign, the encoded object will require quoting. If any value in the object requires quoting, the required escapes will make the encoded field pretty difficult for humans to read.

### Channels and Functions

Channels and functions are encoded as their type and their address. Therearen't many meaningful ways to log channels and functions...

### Maps and Structs

Maps and structs are encoded as strings that contain the result of `fmt.Sprint`.

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
