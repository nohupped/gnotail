package generic

import (
	"io"
	"encoding/json"
	"io/ioutil"
	"os"
	log "github.com/nohupped/glog"
)

func ConfParser(c io.Reader, data *map[string][]string) {
	//conf := map[string][]string{}
	d, err := ioutil.ReadAll(c)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(d, &data)
	if err != nil {
		panic(err)
	}
	for k := range *data {
		fi, err := os.Stat(k)
		if err != nil {
			log.Errorln("Couldn't stat", k, "and won't be parsed.")
			delete(*data, k)
		}else if fi.IsDir() {
			log.Errorln(k, "is not a directory and won't be parsed")
			delete(*data, k)
		}
	}
}