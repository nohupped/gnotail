package main

import (
	"gnotail/generic"
	"gnotail/logger"
	"gnotail/poller"
	"gnotail/tailer"
	log "github.com/nohupped/glog"
	"os"
	"sync"
	"github.com/rjeczalik/notify"
	"gnotail/udp"
)

var mapping generic.Mapping
func init() {
	mapping = generic.CreateMappings()
	logger.InitLog(*mapping.Loglevel, os.Stdout)
	log.Debugf("%+v\n", mapping)
}

func main() {
	// Write to this channel when you want to notify about a log rotation
	rotateNotificationChan := make(chan string, 2048)

	wg := new(sync.WaitGroup)

	UDPClient := udp.NewUDPClient(*mapping.UDPConnAddr, *mapping.Port)


	for k, v := range mapping.MetaMapping {
		poller.SetNewNotifier(k, v.Notifychan, v, UDPClient)
		wg.Add(1)
		v.FD = tailer.NewTailHandler(k)
		v.FD.Seek(v.Offset, 0)
		go poller.FileStat(k, rotateNotificationChan, wg)
		wg.Add(1)
		go poller.ReadInotify(v.Notifychan, v.TriggerChan, wg, v, UDPClient)


	}

	for ; ; {
		select {
		case filename := <- rotateNotificationChan:
			log.Debugln("Trigger received for", filename)
			mapping.MetaMapping[filename].TriggerChan <- filename
			log.Debugln("Closing FD for", filename)
			mapping.MetaMapping[filename].FD.Close()
			notify.Stop(mapping.MetaMapping[filename].Notifychan)
			for ; ;  {
				err := poller.SetNewNotifier(filename, mapping.MetaMapping[filename].Notifychan, mapping.MetaMapping[filename], UDPClient)
				if err != nil {
					log.Debugln(err)
					continue
				}else {break}

			}

			tmp := generic.FileMeta{
				FD: tailer.NewTailHandler(filename),
				Offset: 0,
				Notifychan: mapping.MetaMapping[filename].Notifychan,
				TriggerChan: mapping.MetaMapping[filename].TriggerChan,
				Patterns: mapping.MetaMapping[filename].Patterns,
				Filename: mapping.MetaMapping[filename].Filename,
			}

			mapping.MetaMapping[filename] = tmp
			log.Infoln("Going to retail", filename)
			wg.Add(1)
			go poller.ReadInotify(mapping.MetaMapping[filename].Notifychan, mapping.MetaMapping[filename].TriggerChan, wg,
				mapping.MetaMapping[filename], UDPClient)
		}
	}

	wg.Wait()
}




