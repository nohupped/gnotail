package poller

import (
	"github.com/rjeczalik/notify"
	log "github.com/nohupped/glog"
	"sync"
	"runtime"
	"gnotail/generic"
	"time"
	"gnotail/tailer"
	"net"
)

func SetNewNotifier(filepath string, c chan notify.EventInfo, meta generic.FileMeta, udpclient *net.UDPConn) error{
	log.Infoln("Setting new notifier to", filepath)
	if err := notify.Watch(filepath, c, notify.Write); err != nil {
		log.Errorln(err)
		return err
	}
	time.Sleep(time.Second*1)
	tailer.Tail(meta, udpclient)
	return nil
}

func ReadInotify(r <-chan notify.EventInfo, trigger <-chan string, wg *sync.WaitGroup, meta generic.FileMeta, udpclient *net.UDPConn) {
	defer wg.Done()
	for ; ; {
		select {
		case event := <- r:
			log.Debugln(event.Event(), "on", event.Path())
			//log.Debugln(meta.Offset)
			tailer.Tail(meta, udpclient)
		case done := <-trigger:
			log.Debugln("Received rotate trigger for", done)
			runtime.Goexit()
		}

	}
}
