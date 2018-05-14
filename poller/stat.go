package poller

import (
	"os"
	log"github.com/nohupped/glog"
	"sync"
	"time"
)

// FileStat is supposed to run in a goroutine and continuously check the os.Stat of the filename every 1 second.
// If the size of the file is reduced (truncated, or moved and a new file is created), this will write the name of the file
// to the w channel (supposed to be the trigger channel.) This will do a stat everytime, and doesn't hold an open fd.
func FileStat(filename string, w chan string, wg *sync.WaitGroup) {
	var offset int64
	offset = -1
	for ; ; {
		time.Sleep(time.Second * 1)
		fdStat, err := os.Stat(filename)
		if err != nil {
			log.Errorln(err)
			continue
		}
		//log.Debugln(filename, "size:", fdStat.Size(), "offset:", offset)
		if fdStat.Size() >= offset {
			offset = fdStat.Size()
			continue
		} else {
			offset = -1
			log.Infoln(filename, "truncated, sending notification from Filestat")
			w <- filename
			continue
		}
	}
	wg.Done()
}

