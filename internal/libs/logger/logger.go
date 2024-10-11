package logger

import (
	"fmt"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
)

type CallerFormatter struct {
	*log.TextFormatter
}

func (f *CallerFormatter) Format(entry *log.Entry) ([]byte, error) {
	entry.Data["caller"] = f.getCaller()
	return f.TextFormatter.Format(entry)
}

func (f *CallerFormatter) getCaller() string {
	// Skip this function, and fetch the caller
	_, file, line, ok := runtime.Caller(8)
	if !ok {
		return "unknown:0"
	}

	return filepath.Base(file) + ":" + fmt.Sprintf("%d", line)
}
