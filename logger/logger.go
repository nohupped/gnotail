package logger

import (
	log "github.com/nohupped/glog"
	"io"
)
// InitLog initializes the logger and sets the default loglevel.
func InitLog(loglevel int, output io.Writer) *log.Logger{
	l := log.New(output, "", log.Lshortfile)
	log.SetFlags(log.Lshortfile)
	l.SetLogLevel(loglevel)
	log.SetStandardLogLevel(loglevel)
	return l
}
