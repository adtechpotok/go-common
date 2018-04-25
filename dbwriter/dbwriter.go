package viewwriter

import (
	orm "github.com/adtechpotok/go-orm"
	"sync"
	"time"
)

var serverId = 0
var mysqlBuffer = make([]orm.SchemaPotok, 0)
var mutex = &sync.RWMutex{}

func AppendData(data ...orm.SchemaPotok) {
	mutex.Lock()
	defer mutex.Unlock()
	mysqlBuffer = append(mysqlBuffer, data...)
}

func StartWriter(conf WriteConfig, writer BaseWriter) {
	c := time.Tick(conf.TickTimeMs * time.Millisecond)
	serverId = conf.ServerId
	writer.setConnMaxLifetime(time.Second * conf.MaxConnectTimeSec)
	for range c {
		//Тут мы можем очень удобно всё распределять по врайтерам
		writer.mysqlWrite(&conf)
	}
}
