package gobatis

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	sugar *zap.SugaredLogger
)

func init() {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), os.Stdout, zapcore.InfoLevel)
	logger := zap.New(core).WithOptions()
	defer logger.Sync()
	sugar = logger.Sugar()
}

func SetCustomLog(log *zap.SugaredLogger) {
	sugar = log
}
