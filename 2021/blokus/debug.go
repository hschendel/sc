package blokus

import (
	"io"
	"os"
)

type loggingReader struct {
	Out *os.File
	R   io.Reader
}

func (l *loggingReader) Read(p []byte) (n int, err error) {
	n, err = l.R.Read(p)
	if n > 0 {
		l.Out.Write(p[:n])
		l.Out.Sync()
	}
	return
}

type loggingWriter struct {
	Out *os.File
	W   io.Writer
}

func (l *loggingWriter) Write(p []byte) (n int, err error) {
	n, err = l.W.Write(p)
	if n > 0 {
		l.Out.Write(p[:n])
		l.Out.Sync()
	}
	return
}

