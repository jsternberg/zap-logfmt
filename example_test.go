package zaplogfmt_test

import (
	"os"

	zaplogfmt "github.com/allir/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Example_usage() {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = ""
	logger := zap.New(
		zapcore.NewCore(
			zaplogfmt.NewEncoder(config),
			os.Stdout,
			zapcore.DebugLevel),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	).Named("main")

	logger.Info("Hello World")

	// Output: level=info logger=main caller=zap-logfmt/example_test.go:23 msg="Hello World"
}
