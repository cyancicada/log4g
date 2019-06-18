package log4g

import (
	"encoding/json"
	"sync/atomic"
	"time"
)

type LessLogger struct {
	threshold int
	lastTime  int64
	discarded uint32
}

func NewLessLogger(milliseconds int) *LessLogger {
	return &LessLogger{
		threshold: milliseconds,
		lastTime:  0,
	}
}

func (logger *LessLogger) Error(v ...interface{}) {
	logger.logOrDiscard(func() {
		Error(v...)
	})
}
func (logger *LessLogger) Log(keyValues ...interface{}) error {
	return json.NewEncoder(InfoLog).Encode(keyValues)
}

func (logger *LessLogger) Errorf(format string, v ...interface{}) {
	logger.logOrDiscard(func() {
		ErrorFormat(format, v...)
	})
}

func (logger *LessLogger) logOrDiscard(fn func()) {
	if logger == nil || logger.threshold <= 0 {
		fn()
		return
	}

	currentMillis := time.Now().UnixNano() / int64(time.Millisecond)
	if currentMillis-logger.lastTime > int64(logger.threshold) {
		logger.lastTime = currentMillis
		discarded := atomic.SwapUint32(&logger.discarded, 0)
		if discarded > 0 {
			ErrorFormat("Discarded %d error messages", discarded)
		}

		fn()
	} else {
		atomic.AddUint32(&logger.discarded, 1)
	}
}
