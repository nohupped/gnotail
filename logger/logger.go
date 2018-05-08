package logger

import (
	log "github.com/nohupped/glog"
	"io"
)

func InitLog(loglevel int, output io.Writer) *log.Logger{
	l := log.New(output, "", log.Lshortfile)
	log.SetFlags(log.Lshortfile)
	l.SetLogLevel(loglevel)
	log.SetStandardLogLevel(loglevel)
	return l
}
