package poller

import (
	"os"
	log"github.com/nohupped/glog"
	"sync"
	"time"
)

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

