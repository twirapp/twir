package shell

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type PrefixWriter struct {
	w      io.Writer
	prefix string
	mu     sync.Mutex
	buf    []byte
}

func NewPrefixWriter(w io.Writer, prefix string) *PrefixWriter {
	return &PrefixWriter{w: w, prefix: prefix}
}

func (pw *PrefixWriter) Write(p []byte) (n int, err error) {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	pw.buf = append(pw.buf, p...)

	for {
		idx := -1
		for i, b := range pw.buf {
			if b == '\n' {
				idx = i
				break
			}
		}
		if idx < 0 {
			break
		}

		line := pw.buf[:idx]
		pw.buf = pw.buf[idx+1:]

		if len(line) == 0 {
			fmt.Fprintln(pw.w)
			continue
		}

		fmt.Fprintf(pw.w, "[%s] %s\n", pw.prefix, string(line))
	}

	return len(p), nil
}

func (pw *PrefixWriter) Flush() {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	if len(pw.buf) > 0 {
		fmt.Fprintf(pw.w, "[%s] %s", pw.prefix, string(pw.buf))
		pw.buf = nil
	}
}

var _ io.Writer = (*PrefixWriter)(nil)

func StdoutFor(prefix string) *PrefixWriter {
	return NewPrefixWriter(os.Stdout, prefix)
}

func StderrFor(prefix string) *PrefixWriter {
	return NewPrefixWriter(os.Stderr, prefix)
}
