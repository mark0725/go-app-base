package web

import (
	"bytes"
	"os"
)

type WebLogWriter struct {
	accessLog *os.File
	errorLog  *os.File
	buf       *bytes.Buffer
}

// Write method for WebLogWriter to satisfy the io.Writer interface
func (w WebLogWriter) Write(p []byte) (n int, err error) {
	// In a production scenario, you might want a more sophisticated log parsing
	if bytes.Contains(p, []byte("ERROR")) {
		if w.errorLog != nil {
			w.errorLog.Write(p)
		}
	} else {
		if w.accessLog != nil {
			w.accessLog.Write(p)
		}
	}
	return w.buf.Write(p)
}
