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
// NewTailHandler returns an fd for the filename. This is stored in the FileMeta struct.
// This is kept open until the FileStat founds a truncated/removed file and triggers the rotate.
// This open fd is closed from the main. 
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

// Parsed struct is used to json marshal the message and to add relevant metadata along with it. 
type Parsed struct {
	Hostname string `json:"hostname"`
	Filename string `json:"filename"`
	Rule     string `json:"rule"`
	Message  string `json:"message"`
}

// Tail will read from the fd and marks the last seek offset and write a udp datagram until an EOF error or an 
// empty slice is occured and the function exits. It doesn't run continuously. This function is called again
// later when an inotify event is occured.
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
						Hostname: meta.Hostname,
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
