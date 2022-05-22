package util_test

import (
	"testing"

	"pig/util/logger"
)

func TestLogger(t *testing.T) {
	logger := logger.Logger{
		Level: logger.Info,
	}
	logger.Trace("trace", 1, []byte{0})
	logger.Debug("debug", 2, []byte{0, 1})
	logger.Info("info", 3, []byte{0, 1, 2, 3})
	logger.Warn("warn", 4, []byte{0, 1, 2, 3, 4})
	logger.Error("error", 5, []byte{0, 1, 2, 3, 4, 5})
	logger.Fatal("fatal", 6, []byte{0, 1, 2, 3, 4, 5, 6})
}
