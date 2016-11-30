package zlogfmt

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/uber-go/zap"
)

const (
	defaultTimeKey    = "ts"
	defaultLevelKey   = "level"
	defaultMessageKey = "msg"
	hex               = "0123456789abcdef"
	_initialBufSize   = 1024
)

var logfmtPool = sync.Pool{New: func() interface{} {
	return &logfmtEncoder{
		bytes: make([]byte, 0, _initialBufSize),
	}
}}

type logfmtEncoder struct {
	bytes      []byte
	timeKey    string
	levelKey   string
	messageKey string
}

func NewLogfmtEncoder(options ...LogfmtOption) zap.Encoder {
	enc := logfmtPool.Get().(*logfmtEncoder)
	enc.truncate()

	enc.timeKey = defaultTimeKey
	enc.levelKey = defaultLevelKey
	enc.messageKey = defaultMessageKey
	for _, opt := range options {
		opt.apply(enc)
	}
	return enc
}

func (enc *logfmtEncoder) Free() {
	logfmtPool.Put(enc)
}

func (enc *logfmtEncoder) AddString(key, val string) {
	enc.addKey(key)
	if strings.IndexFunc(val, needsQuotedValueRune) != -1 {
		enc.safeAddString(val)
	} else {
		enc.bytes = append(enc.bytes, val...)
	}
}

func (enc *logfmtEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	if val {
		enc.bytes = append(enc.bytes, []byte("true")...)
	} else {
		enc.bytes = append(enc.bytes, []byte("false")...)
	}
}

func (enc *logfmtEncoder) AddInt(key string, val int) {
	enc.AddInt64(key, int64(val))
}

func (enc *logfmtEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.bytes = strconv.AppendInt(enc.bytes, val, 10)
}

func (enc *logfmtEncoder) AddUint(key string, val uint) {
	enc.AddUint64(key, uint64(val))
}

func (enc *logfmtEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.bytes = strconv.AppendUint(enc.bytes, val, 10)
}

func (enc *logfmtEncoder) AddUintptr(key string, val uintptr) {
}

func (enc *logfmtEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.bytes = strconv.AppendFloat(enc.bytes, val, 'g', 3, 64)
}

func (enc *logfmtEncoder) AddMarshaler(key string, obj zap.LogMarshaler) error {
	return errors.New("unimplemented")
}

func (enc *logfmtEncoder) AddObject(key string, obj interface{}) error {
	return errors.New("unimplemented")
}

func (enc *logfmtEncoder) Clone() zap.Encoder {
	clone := logfmtPool.Get().(*logfmtEncoder)
	clone.truncate()
	clone.bytes = make([]byte, 0, cap(clone.bytes))
	clone.bytes = append(clone.bytes, enc.bytes...)
	clone.timeKey = enc.timeKey
	clone.levelKey = enc.levelKey
	clone.messageKey = enc.messageKey
	return clone
}

func (enc *logfmtEncoder) WriteEntry(sink io.Writer, msg string, level zap.Level, t time.Time) error {
	if sink == nil {
		return errors.New("can't write encoded message to a nil WriteSyncer")
	}

	final := logfmtPool.Get().(*logfmtEncoder)
	final.truncate()
	final.AddString(enc.timeKey, t.Format(time.RFC3339Nano))
	final.AddString(enc.levelKey, level.String())
	final.AddString(enc.messageKey, msg)
	if len(enc.bytes) > 0 {
		final.bytes = append(final.bytes, ' ')
		final.bytes = append(final.bytes, enc.bytes...)
	}
	final.bytes = append(final.bytes, '\n')

	expectedBytes := len(final.bytes)
	n, err := sink.Write(final.bytes)
	final.Free()
	if err != nil {
		return err
	} else if n != expectedBytes {
		return fmt.Errorf("incomplete write: only wrote %v of %v bytes", n, expectedBytes)
	}
	return nil
}

func (enc *logfmtEncoder) truncate() {
	enc.bytes = enc.bytes[:0]
}

func (enc *logfmtEncoder) addKey(key string) {
	if len(enc.bytes) > 0 {
		enc.bytes = append(enc.bytes, ' ')
	}
	enc.bytes = append(enc.bytes, key...)
	enc.bytes = append(enc.bytes, '=')
}

func (enc *logfmtEncoder) safeAddString(val string) {
	enc.bytes = append(enc.bytes, '"')
	start := 0
	for i := 0; i < len(val); {
		if b := val[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' {
				i++
				continue
			}
			if start < i {
				enc.bytes = append(enc.bytes, val[start:i]...)
			}
			switch b {
			case '\\', '"':
				enc.bytes = append(enc.bytes, '\\', b)
			case '\n':
				enc.bytes = append(enc.bytes, '\\', 'n')
			case '\r':
				enc.bytes = append(enc.bytes, '\\', 'r')
			case '\t':
				enc.bytes = append(enc.bytes, '\\', 't')
			default:
				enc.bytes = append(enc.bytes, `\u00`...)
				enc.bytes = append(enc.bytes, hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(val[i:])
		if c == utf8.RuneError {
			if start < i {
				enc.bytes = append(enc.bytes, val[start:i]...)
			}
			enc.bytes = append(enc.bytes, "\ufffd"...)
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(val) {
		enc.bytes = append(enc.bytes, val[start:]...)
	}
	enc.bytes = append(enc.bytes, '"')
}

func needsQuotedValueRune(r rune) bool {
	return r <= ' ' || r == '=' || r == '"' || r == utf8.RuneError
}

type LogfmtOption interface {
	apply(*logfmtEncoder)
}
