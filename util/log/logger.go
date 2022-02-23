package log

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const (
	Trace = iota
	Debug
	Info
	Warn
	Error
	Fatal
)

var (
	colorMap map[int]string
	Dog *Logger
)

func init() {
	colorMap = map[int]string{
		Trace: "\033[90m",   // dark gray
		Debug: "\033[96m",   // cyan
		Info:  "\033[92m",   // green
		Warn:  "\033[93m",   // yellow
		Error: "\033[91m",   // red
		Fatal: "\033[1;91m", // bold red; 41;1;37m white&bold with red background
	}
	Dog = &Logger{
		Level: Trace,
		Name: "WatchDog",
		DateFmt: "2006-01-02 15:04:05",
	}
}

type Logger struct {
	Level    int
	Name     string
	// header
	DateFmt	 string
	NoCaller bool // don't print caller
}

func (logger *Logger) Trace(args ...interface{}) {
	logger.log(Trace, args...)
}

func (logger *Logger) Tracef(format string, args ...interface{}) {
	logger.log(Trace, fmt.Sprintf(format, args...))
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.log(Debug, args...)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.log(Debug, fmt.Sprintf(format, args...))
}

func (logger *Logger) Info(args ...interface{}) {
	logger.log(Info, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.log(Info, fmt.Sprintf(format, args...))
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.log(Warn, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.log(Warn, fmt.Sprintf(format, args...))
}

func (logger *Logger) Error(args ...interface{}) {
	logger.log(Error, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.log(Error, fmt.Sprintf(format, args...))
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.log(Fatal, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.log(Fatal, fmt.Sprintf(format, args...))
}

func (logger *Logger) log(level int, args ...interface{}) {
	if level < logger.Level {
		return
	}

	// make string builder
	sb := strings.Builder{}

	// format arguments
	for _, arg := range args {
		sb.WriteString(fmt.Sprintf("%v ", arg))
	}

	// trim last space
	rs := sb.String()
	if sb.Len() > 0 {
		rs = rs[:len(rs)-1]
	}

	// add color
	cs := dye(rs, level)

	// add header
	if !logger.NoCaller || len(logger.DateFmt) > 0 {
		header := ""

		// add date
		if len(logger.DateFmt) > 0 {
			now := time.Now()
			header += now.Format(logger.DateFmt) + " "
		}
		// add caller
		if (!logger.NoCaller) {
			pc, _, _, _ := runtime.Caller(2)
			caller := runtime.FuncForPC(pc).Name()
			header += caller
		}

		cs = "[" + strings.Trim(header, " ") + "] " + cs
	}

	fmt.Println(cs + "\r")
}

// add color to string
func dye(string string, level int) string {
	if level < Trace || level > Fatal {
		return string
	}
	return colorMap[level] + string + "\033[0m"
}
