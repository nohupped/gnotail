package generic

import (
	"github.com/rjeczalik/notify"
	"flag"
	"regexp"
	"os"
)

// FileToPattern is used to map the file to its corresponding slice of patterns.
var FileToPattern map[string][]string

// Mapping holds generic configuration (port, loglevel, addr to which the udp datagram is sent and 
// slice of FileMeta struct).
type Mapping struct {
	Port, Loglevel *int
	UDPConnAddr *string
	MetaMapping map[string]FileMeta
}


// FileMeta holds the meta data and a trigger channel to interrupt the 
// goroutine for each file when it is tailed. It is nested inside Mapping{} struct.
type FileMeta struct {
	Hostname string
	Filename string
	Notifychan chan notify.EventInfo
	Offset int64
	Patterns []*regexp.Regexp
	TriggerChan chan string
	FD *os.File
}// Maps filename to inotify.EventInfo chan for setting watchers for each files

// CreateMappings is used to parse the configuration file and create 
// the mappings required by the program to tail and filter each file.
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
	hostname := flag.String("hostname", "", "Hostname to be tagged in the filtered output. Defaults to system hostname.")
	flag.Parse()
	if *conf == "" || conf == nil {
		panic("Give a configuration.")
	}
	if *hostname == "" {
		var err error
		*hostname, err = os.Hostname()
		if err != nil {
			panic(err)
		}
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

		// Updates each file's meta information
		mapping.MetaMapping[k] = FileMeta{
			Hostname: *hostname,
			Patterns:   transformPatterns(v),
			Notifychan: make(chan notify.EventInfo, 2048),
			Offset:     stat.Size(), // This is not updated as tail follows.
			TriggerChan: make(chan string),
			Filename: k,
		}
	}

	return *mapping
}

// Converts the patterns mentioned in the config file to slice of regex objects.
func transformPatterns(p []string) []*regexp.Regexp{
	var patterns []*regexp.Regexp
	for _, i := range p {
		pattern := regexp.MustCompile(i)
		patterns = append(patterns, pattern)
	}
	return patterns
}
