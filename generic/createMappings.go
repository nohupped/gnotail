package generic

import (
	"github.com/rjeczalik/notify"
	"flag"
	"regexp"
	"os"
)
var FileToPattern map[string][]string

type Mapping struct {
	Port, Loglevel *int
	UDPConnAddr *string
	MetaMapping map[string]FileMeta
}


type FileMeta struct {
	Filename string
	Notifychan chan notify.EventInfo
	Offset int64
	Patterns []*regexp.Regexp
	TriggerChan chan string
	FD *os.File
}// Maps filename to inotify.EventInfo chan for setting watchers for each files

func CreateMappings() Mapping{
	mapping := new(Mapping)

	conf := flag.String("conf", "", `path to configuration of json format. Example format would be 
{
   "/var/log/syslog": [
       "cron",
       "su.*"
   ],
   "/var/log/auth.log": [
       "sudo",
       "pam_unix"
   ]
}`)

	mapping.Port = flag.Int("port", 9999, "Portnumber to which the output be sent")
	mapping.Loglevel = flag.Int("loglevel", 2, "Log level when printing to STDOUT. ErrorLevel: 0, " +
		"WarnLevel: 1, InfoLevel: 2, DebugLevel: 3. Defaults to InfoLevel.")
	mapping.UDPConnAddr = flag.String("udp_conn_addr", "127.0.0.1", "address to start udp server")
	flag.Parse()
	if *conf == "" || conf == nil {
		panic("Give a configuration.")
	}

	if *mapping.Loglevel > 3 {
		panic("Set loglevel less than 4. \n ErrorLevel: 0, WarnLevel: 1, InfoLevel: 2, DebugLevel: 3")
	}

	fd, err := Opener(*conf)
	defer fd.Close()
	if err != nil {
		panic(err)
	}
	mapping.MetaMapping = make(map[string]FileMeta)
	ConfParser(fd, &FileToPattern)
	for k, v := range FileToPattern {
		stat, err := os.Stat(k)
		if err != nil {
			panic(err)
		}

		mapping.MetaMapping[k] = FileMeta{
			Patterns:   transformPatterns(v),
			Notifychan: make(chan notify.EventInfo, 2048),
			Offset:     stat.Size(), // This is not updated as tail follows.
			//Offset: 0,
			TriggerChan: make(chan string),
			Filename: k,
		}
	}

	return *mapping
}

func transformPatterns(p []string) []*regexp.Regexp{
	var patterns []*regexp.Regexp
	for _, i := range p {
		pattern := regexp.MustCompile(i)
		patterns = append(patterns, pattern)
	}
	return patterns
}
