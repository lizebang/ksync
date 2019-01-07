package log

import (
	"go.uber.org/zap"
)

var (
	sugar *zap.SugaredLogger

	Debug  func(args ...interface{})
	Info   func(args ...interface{})
	Warn   func(args ...interface{})
	Error  func(args ...interface{})
	DPanic func(args ...interface{})
	Panic  func(args ...interface{})
	Fatal  func(args ...interface{})

	Debugf  func(template string, args ...interface{})
	Infof   func(template string, args ...interface{})
	Warnf   func(template string, args ...interface{})
	Errorf  func(template string, args ...interface{})
	DPanicf func(template string, args ...interface{})
	Panicf  func(template string, args ...interface{})
	Fatalf  func(template string, args ...interface{})

	Debugw  func(msg string, keysAndValues ...interface{})
	Infow   func(msg string, keysAndValues ...interface{})
	Warnw   func(msg string, keysAndValues ...interface{})
	Errorw  func(msg string, keysAndValues ...interface{})
	DPanicw func(msg string, keysAndValues ...interface{})
	Panicw  func(msg string, keysAndValues ...interface{})
	Fatalw  func(msg string, keysAndValues ...interface{})
)

func init() {
	logger, err := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout", "ksync.log"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.FatalLevel),
	)
	if err != nil {
		panic(err)
	}
	sugar = logger.Sugar()

	Debug = sugar.Debug
	Info = sugar.Info
	Warn = sugar.Warn
	Error = sugar.Error
	DPanic = sugar.DPanic
	Panic = sugar.Panic
	Fatal = sugar.Fatal

	Debugf = sugar.Debugf
	Infof = sugar.Infof
	Warnf = sugar.Warnf
	Errorf = sugar.Errorf
	DPanicf = sugar.DPanicf
	Panicf = sugar.Panicf
	Fatalf = sugar.Fatalf

	Debugw = sugar.Debugw
	Infow = sugar.Infow
	Warnw = sugar.Warnw
	Errorw = sugar.Errorw
	DPanicw = sugar.DPanicw
	Panicw = sugar.Panicw
	Fatalw = sugar.Fatalw
}
