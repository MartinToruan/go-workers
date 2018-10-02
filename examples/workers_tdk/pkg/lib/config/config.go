package config

import (
	"log"
	"sync"

	"github.com/tokopedia/tdk/go/config"
)

type Config struct {
	Appname string
}

var cfg Config
var once = sync.Once{}

func GetConfig() Config {
	once.Do(func() {
		err := config.Read(&cfg, "/Users/nakama/go/src/workers/examples/workers_tdk/config/workers_tdk.{TKPENV}.yaml", "/etc/workers_tdk/workers_tdk.{TKPENV}.yaml")
		if err != nil {
			log.Fatal(err)
		}
	})
	return cfg
}
