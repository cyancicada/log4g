package log4g

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	dateFormat      = "2006-01-02"
	hoursPerDay     = 24
	bufferSize      = 100
	defaultDirMode  = 0755
	defaultFileMode = 0600
)

var ErrLogFileClosed = errors.New("error: log file closed")

type (
	BackupRule interface {
		BackupFileName() string
		MarkRotated()
		OutdatedFiles() []string
		ShallRotate() bool
		GetPrefix() string
	}
	Placeholder  struct{}
	BackupLogger struct {
		filename  string
		backup    string
		fp        *os.File
		channel   chan []byte
		done      chan Placeholder
		rule      BackupRule
		compress  bool
		stdout    bool
		keepDays  int
		waitGroup sync.WaitGroup
		closeOnce sync.Once
	}

	DailyBackupRule struct {
		rotatedTime string
		filename    string
		prefix      string
		delimiter   string
		days        int
		gzip        bool
	}
)

func DefaultBackupRule(filename, prefix, delimiter string, days int, gzip bool) BackupRule {
	return &DailyBackupRule{
		rotatedTime: getNowDate(),
		filename:    filename,
		delimiter:   delimiter,
		days:        days,
		prefix:      prefix,
		gzip:        gzip,
	}
}

func (r *DailyBackupRule) BackupFileName() string {
	return fmt.Sprintf("%s%s%s", r.filename, r.delimiter, getNowDate())
}

func (r *DailyBackupRule) GetPrefix() string {
	return r.prefix
}

func (r *DailyBackupRule) MarkRotated() {
	r.rotatedTime = getNowDate()
}

func (r *DailyBackupRule) OutdatedFiles() []string {
	if r.days <= 0 {
		return nil
	}

	var pattern string
	if r.gzip {
		pattern = fmt.Sprintf("%s%s*.gz", r.filename, r.delimiter)
	} else {
		pattern = fmt.Sprintf("%s%s*", r.filename, r.delimiter)
	}

	files, err := filepath.Glob(pattern)
	if err != nil {
		ErrorFormat("failed to delete outdated log files, error: %s", err)
		return nil
	}

	var buf strings.Builder
	boundary := time.Now().Add(-time.Hour * time.Duration(hoursPerDay*r.days)).Format(dateFormat)
	if _, err := fmt.Fprintf(&buf, "%s%s%s", r.filename, r.delimiter, boundary); err != nil {
		ErrorFormat("failed to fmt.Fprintf %s%s%s err is %s", r.filename, r.delimiter, boundary, err)
		return nil
	}
	if r.gzip {
		buf.WriteString(".gz")
	}
	boundaryFile := buf.String()

	var outDates []string
	for _, file := range files {
		if file < boundaryFile {
			outDates = append(outDates, file)
		}
	}

	return outDates
}

func (r *DailyBackupRule) ShallRotate() bool {
	return len(r.rotatedTime) > 0 && getNowDate() != r.rotatedTime
}

func NewLogger(filename string, stdout bool, rule BackupRule, compress bool) (*BackupLogger, error) {
	l := &BackupLogger{
		filename: filename,
		channel:  make(chan []byte, bufferSize),
		done:     make(chan Placeholder),
		rule:     rule,
		stdout:   stdout,
		compress: compress,
	}
	if err := l.init(); err != nil {
		return nil, err
	}

	l.startWorker()
	return l, nil
}

func (l *BackupLogger) Close() error {
	var err error

	l.closeOnce.Do(func() {
		close(l.done)
		l.waitGroup.Wait()

		if err = l.fp.Sync(); err != nil {
			return
		}

		err = l.fp.Close()
	})

	return err
}

func (l *BackupLogger) Write(data []byte) (int, error) {
	logBytes := append([]byte(l.rule.GetPrefix()), data...)
	select {
	case l.channel <- logBytes:
		if l.stdout {
			fmt.Print(string(logBytes))
		}
		return len(logBytes), nil
	case <-l.done:
		stdoutErrOutput(l.rule.GetPrefix(), string(logBytes), 1)
		return 0, ErrLogFileClosed
	}
}

func (l *BackupLogger) getBackupFilename() string {
	if len(l.backup) == 0 {
		return l.rule.BackupFileName()
	} else {
		return l.backup
	}
}

func (l *BackupLogger) init() error {
	l.backup = l.rule.BackupFileName()

	if _, err := os.Stat(l.filename); err != nil {
		basePath := path.Dir(l.filename)
		if _, err = os.Stat(basePath); err != nil {
			if err = os.MkdirAll(basePath, defaultDirMode); err != nil {
				return err
			}
		}

		if l.fp, err = os.Create(l.filename); err != nil {
			return err
		}
	} else if l.fp, err = os.OpenFile(l.filename, os.O_APPEND|os.O_WRONLY, defaultFileMode); err != nil {
		return err
	}

	CloseOnExec(l.fp)

	return nil
}

func (l *BackupLogger) maybeCompressFile(file string) {
	if l.compress {
		defer func() {
			if r := recover(); r != nil {
				Server(r)
			}
		}()
		compressLogFile(file)
	}
}

func (l *BackupLogger) maybeDeleteOutdatedFiles() {
	files := l.rule.OutdatedFiles()
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			ErrorFormat("failed to remove outdated file: %s", file)
		}
	}
}

func (l *BackupLogger) postRotate(file string) {
	go func() {
		// we cannot use threading.GoSafe here, because of import cycle.
		l.maybeCompressFile(file)
		l.maybeDeleteOutdatedFiles()
	}()
}

func (l *BackupLogger) rotate() error {
	if l.fp != nil {
		err := l.fp.Close()
		l.fp = nil
		if err != nil {
			return err
		}
	}

	_, err := os.Stat(l.filename)
	if err == nil && len(l.backup) > 0 {
		backupFilename := l.getBackupFilename()
		err = os.Rename(l.filename, backupFilename)
		if err != nil {
			return err
		}

		l.postRotate(backupFilename)
	}

	l.backup = l.rule.BackupFileName()
	if l.fp, err = os.Create(l.filename); err == nil {
		CloseOnExec(l.fp)
	}

	return err
}

func (l *BackupLogger) startWorker() {
	l.waitGroup.Add(1)

	go func() {
		defer l.waitGroup.Done()

		for {
			select {
			case event := <-l.channel:
				l.write(event)
			case <-l.done:
				return
			}
		}
	}()
}

func (l *BackupLogger) write(v []byte) {
	if l.rule.ShallRotate() {
		if err := l.rotate(); err != nil {
			log.Println(err)
		} else {
			l.rule.MarkRotated()
		}
	}
	if l.fp != nil {
		if _, err := l.fp.Write(v); err != nil {
			ErrorFormat("write l.fp.Write err %+v", err)
		}
	}
}

func compressLogFile(file string) {
	start := time.Now()
	InfoFormat("compressing log file: %s", file)
	if err := gzipFile(file); err != nil {
		ErrorFormat("compress error: %s", err)
	} else {
		InfoFormat("compressed log file: %s, took %s", file, time.Since(start))
	}
}

func getNowDate() string {
	return time.Now().Format(dateFormat)
}

func gzipFile(file string) error {
	in, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		if err := in.Close(); err != nil {
			ErrorFormat("gzipFile.in.Close err %+v", err)
		}
	}()

	out, err := os.Create(fmt.Sprintf("%s.gz", file))
	if err != nil {
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			ErrorFormat("gzipFile.out.Close err %+v", err)
		}
	}()

	w := gzip.NewWriter(out)
	if _, err = io.Copy(w, in); err != nil {
		return err
	} else if err = w.Close(); err != nil {
		return err
	}

	return os.Remove(file)
}
