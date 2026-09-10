package opennox

import (
	"bufio"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/opennox/libs/env"
	"github.com/opennox/libs/log"
)

var (
	logFile    io.Closer
	logBuf     bufio.Writer
	logFileHnd = log.NewTextHandler(nil)
)

func init() {
	logFileHnd.SetLevel(slog.LevelDebug)
	log.AddHandler(logFileHnd)
}

func closeLog() {
	setLogFile(nil)
}

type syncLogWriter struct {
	w io.WriteCloser
}

func (s syncLogWriter) Write(p []byte) (int, error) {
	n, err := s.w.Write(p)
	if f, ok := s.w.(interface{ Sync() error }); ok {
		_ = f.Sync()
	}
	return n, err
}

func (s syncLogWriter) Close() error {
	return s.w.Close()
}

func setLogFile(w io.WriteCloser) {
	if f := logFile; f != nil {
		_ = f.Close()
	}
	if w == nil {
		logFile = nil
		logFileHnd.SetWriter(nil)
		return
	}
	logFile = w
	logFileHnd.SetWriter(syncLogWriter{w: w})
}

func defaultLogDir() string {
	logDir := filepath.Dir(os.Args[0])
	if sdir := env.AppUserDir(); sdir != "" {
		logDir = sdir
	}
	return filepath.Join(logDir, "logs")
}

func writeLogsToDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	name := "opennox.log"
	if isDedicatedServer {
		name = "opennox-server.log"
	}
	f, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	setLogFile(f)
	return nil
}
