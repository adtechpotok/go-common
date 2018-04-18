package viewwriter

import (
	"sync"
	"time"
	"fmt"
	"os"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"github.com/pkg/errors"
	orm "gitlab.rambler.ru/More/go-orm"
)

type WriteConfig struct {
	Db                *gorm.DB
	Log               *logrus.Logger
	FilePath          string
	ServerId          int
	TickTimeMs        time.Duration
	MaxConnectTimeSec time.Duration
}

const attempsLimit = 2

var serverId = 0
var writeBuffer = make([]orm.SchemaPotok, 0)
var mysqlBuffer = make([]orm.SchemaPotok, 0)
var attemps = 0
var mutex = &sync.RWMutex{}

func AppendData(data ...orm.SchemaPotok) {
	mutex.Lock()
	defer mutex.Unlock()
	mysqlBuffer = append(mysqlBuffer, data...)
}

func Writer(conf WriteConfig) {
	c := time.Tick(conf.TickTimeMs * time.Millisecond)
	serverId = conf.ServerId
	conf.Db.DB().SetConnMaxLifetime(time.Second * conf.MaxConnectTimeSec)
	for range c {
		mysqlWrite(&conf)
	}
}

func mysqlWrite(conf *WriteConfig) {

	//Смотрим остались ли данные с прошлой записи
	if len(writeBuffer) == 0 { // если данных нет, заполняем их из текущего буфера
		mutex.Lock()
		writeBuffer = mysqlBuffer
		mysqlBuffer = make([]orm.SchemaPotok, 0)
		mutex.Unlock()
	}

	if conf.Db.DB().Ping() == nil { //если конект есть
		err := readFiles(conf) //записываем все файлы
		if err != nil {
			return
		}
	}

	if len(writeBuffer) == 0 { // если данных нет прекращаем работу
		return
	}
	conf.Log.Infof("Got %d element for mysql write", len(writeBuffer))
	if attemps > attempsLimit { //используем попытки для понимания пишем мы в файл или в базу
		if err := fileWrite(conf); err != nil {
			conf.Log.Error(err, "Cannot write to file.")
			return
		}
		writeBuffer = make([]orm.SchemaPotok, 0)
	}

	if conf.Db.DB().Ping() == nil { //если конект есть
		i := 0
		for _, val := range getQueryData() {
			_, err := conf.Db.DB().Exec(val)
			if err != nil {
				conf.Log.Error("Error while writing in db", err)
				attemps++
				writeBuffer = writeBuffer[i:]
				return
			}
			i++
		}

		conf.Log.Infof("Inserted %d rows", i)
		writeBuffer = make([]orm.SchemaPotok, 0)
		attemps = 0
	} else {
		attemps++
	}

}

func fileWrite(conf *WriteConfig) error {

	fileName := fmt.Sprintf("%s%d.sql", conf.FilePath, time.Now().UnixNano())
	file, err := os.Create(fileName)
	defer file.Close()

	if err != nil {
		return err
	}

	file.WriteString(strings.Join(getQueryData(), "\n"))
	return nil
}

func getQueryData() []string {
	var res []string
	for _, val := range writeBuffer {
		k := makeInsertQuery(val)
		res = append(res, k)
	}
	return res
}

func readFiles(conf *WriteConfig) error {
	files, err := ioutil.ReadDir("./" + conf.FilePath)
	if err != nil {
		conf.Log.Error("Cannot read from file", err)
	}

	for _, f := range files {
		b, err := ioutil.ReadFile(conf.FilePath + f.Name())
		if err != nil {
			conf.Log.Error(err)
		}
		str := strings.Split(string(b), "\n")
		i := 0
		tx, _ := conf.Db.DB().Begin()
		for _, val := range str {
			_, err := tx.Exec(val)
			if err != nil {
				conf.Log.Error("Error while writing in db from file", err)
				attemps++
				tx.Rollback()
				return err
			}
			i++
		}

		conf.Log.Infof("FROM FILE inserted %d rows", i)

		err = os.Remove(conf.FilePath + f.Name())
		if err != nil {
			tx.Rollback()
			return errors.New("Cannot delete file")

		}
		tx.Commit()

	}
	if len(files) > 0 {
		attemps = 0
	}

	return nil
}
