package main

import (
	"math/rand"

	"pig/lib/cop"
	"pig/util/logger"
)

// do something `testing` unable
func main() {
	logger := logger.Dog

	logger.Trace("trace", 1, []byte{0})
	logger.Debug("debug", 2, []byte{0, 1})
	logger.Info("info", 3, []byte{0, 1, 2, 3})
	// log.Dog.Level = log.Fatal
	logger.Warn("warn", 4, []byte{0, 1, 2, 3, 4})
	logger.Error("error", 5, []byte{0, 1, 2, 3, 4, 5})
	logger.Fatal("fatal", 6, []byte{0, 1, 2, 3, 4, 5, 6})

	// when an integer overflowed, it'll come to 0
	x := uint8(0xff)
	logger.Trace(x + 1)

	for i := 0; i < 10; i++ {
		logger.Info(rand.Uint32())
	}
	var data []byte = nil
	logger.Infof("len %v, %v", len(data), data)

	header := &cop.Header{}

	logger.Trace(header)
	logger.Debug(header)
	logger.Info(header)
	logger.Warn(header)
	logger.Error(header)
}
