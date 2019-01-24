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
			zapcore.Lock(os.Stdout),
			zapcore.DebugLevel),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	).Named("main")
	defer logger.Sync()

	logger.Info("Hello World")

	// Output: level=info logger=main caller=zap-logfmt/example_test.go:24 msg="Hello World"
}

func Example_usage_register() {
	zaplogfmt.Register()
	zapConfig := zap.NewProductionConfig()
	zapConfig.Encoding = "logfmt"
	zapConfig.OutputPaths = []string{"stdout"} // Need to log to stdout for the test

	zapConfig.EncoderConfig = zap.NewProductionEncoderConfig()
	zapConfig.EncoderConfig.TimeKey = ""
	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger = logger.Named("reg")

	logger.Info("Hello World")

	// Output: level=info logger=reg caller=zap-logfmt/example_test.go:47 msg="Hello World"
}

func Example_array() {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = ""

	logger := zap.New(zapcore.NewCore(
		zaplogfmt.NewEncoder(config),
		os.Stdout,
		zapcore.DebugLevel,
	))

	logger.Info("counting", zap.Ints("values", []int{0, 1, 2, 3}))

	// Output: level=info msg=counting values=[0,1,2,3]
}
