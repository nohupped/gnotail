package tailer

import (
	"bufio"
	log "github.com/nohupped/glog"
	"os"
	"io"
	"encoding/json"
	"gnotail/generic"
	"net"
)

var err error

func NewTailHandler(filename string) *os.File {
	var fd *os.File
	for i := 0; i <= 10; i++ {
		fd, err = os.Open(filename)
		if err != nil {
			log.Errorln(err)
			continue
		} else {
			break
		}
	}
	return fd
}

type Parsed struct {
	Filename string `json:"filename"`
	Rule     string `json:"rule"`
	Message  string `json:"message"`
}

func Tail(meta generic.FileMeta, udpclient *net.UDPConn) int64 {
	reader := bufio.NewReader(meta.FD)
	offset := 0
	for {
		line, _, _ := reader.ReadLine()
		if err == io.EOF {
			log.Errorln(err)
			break
		} else if len(line) == 0 {
			//log.Debugln("Empty line when reading", meta.Filename)
			break
		} else {
			for _, i := range meta.Patterns {
				if i.Match(line) {
					parsed := Parsed{
						Filename: meta.Filename,
						Rule:     i.String(),
						Message:  string(line),
					}
					msg, err := json.Marshal(parsed)
					if err != nil {
						log.Errorln("Couldn't json marshal", string(line))
					}
					//fmt.Println(string(msg))
					udpclient.Write(msg)
				} // else {
					//log.Errorln("Mismatch for rule", i.String(), string(line))
				//}
			}
			offset += len(line)
		}

	}
	return int64(offset)
}
