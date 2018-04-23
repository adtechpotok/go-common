package dbwriter

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
	orm "github.com/adtechpotok/go-orm"
)

type WriteConfig struct {
	Db                *gorm.DB
	Log               *logrus.Logger
	FilePath          string
	ServerId          int
	TickTimeMs        time.Duration
	MaxConnectTimeSec time.Duration
}

type Writer struct {
	config      WriteConfig
	writeBuffer []orm.SchemaPotok
	mysqlBuffer []orm.SchemaPotok
	mutex       *sync.RWMutex
	attemps     int
}

func New(config WriteConfig) *Writer {
	m := &Writer{config: config}
	m.mutex = &sync.RWMutex{}
	go m.work()

	return m
}

const attempsLimit = 2

func (m *Writer) Append(data ...orm.SchemaPotok) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.mysqlBuffer = append(m.mysqlBuffer, data...)
}

func (m *Writer) work() {
	c := time.Tick(m.config.TickTimeMs * time.Millisecond)
	m.config.Db.DB().SetConnMaxLifetime(time.Second * m.config.MaxConnectTimeSec)
	for range c {
		m.mysqlWrite()
	}
}

func (m *Writer) mysqlWrite() {

	//Смотрим остались ли данные с прошлой записи
	if len(m.writeBuffer) == 0 { // если данных нет, заполняем их из текущего буфера
		m.mutex.Lock()
		m.writeBuffer = m.mysqlBuffer
		m.mysqlBuffer = make([]orm.SchemaPotok, 0)
		m.mutex.Unlock()
	}

	if m.config.Db.DB().Ping() == nil { //если коннект есть
		err := m.readFiles() //записываем все файлы
		if err != nil {
			return
		}
	}

	if len(m.writeBuffer) == 0 { // если данных нет прекращаем работу
		return
	}
	m.config.Log.Infof("Got %d element for mysql write", len(m.writeBuffer))
	if m.attemps > attempsLimit { //используем попытки для понимания пишем мы в файл или в базу
		if err := m.fileWrite(); err != nil {
			m.config.Log.Error(err, "Cannot write to file.")
			return
		}
		m.writeBuffer = make([]orm.SchemaPotok, 0)
	}

	if m.config.Db.DB().Ping() == nil { //если конект есть
		i := 0
		for _, val := range m.getQueryData() {
			_, err := m.config.Db.DB().Exec(val)
			if err != nil {
				m.config.Log.Error("Error while writing in db", err)
				m.attemps++
				m.writeBuffer = m.writeBuffer[i:]
				return
			}
			i++
		}

		m.config.Log.Infof("Inserted %d rows", i)
		m.writeBuffer = make([]orm.SchemaPotok, 0)
		m.attemps = 0
	} else {
		m.attemps++
	}

}

func (m *Writer) fileWrite() error {

	fileName := fmt.Sprintf("%s%d.sql", m.config.FilePath, time.Now().UnixNano())
	file, err := os.Create(fileName)
	defer file.Close()

	if err != nil {
		return err
	}

	file.WriteString(strings.Join(m.getQueryData(), "\n"))
	return nil
}

func (m *Writer) getQueryData() []string {
	var res []string
	for _, val := range m.writeBuffer {
		k := makeInsertQuery(val, m.config.ServerId)
		res = append(res, k)
	}
	return res
}

func (m *Writer) readFiles() error {
	files, err := ioutil.ReadDir("./" + m.config.FilePath)
	if err != nil {
		m.config.Log.Error("Cannot read from file", err)
	}

	for _, f := range files {
		b, err := ioutil.ReadFile(m.config.FilePath + f.Name())
		if err != nil {
			m.config.Log.Error(err)
		}
		str := strings.Split(string(b), "\n")
		i := 0
		tx, _ := m.config.Db.DB().Begin()
		for _, val := range str {
			_, err := tx.Exec(val)
			if err != nil {
				m.config.Log.Error("Error while writing in db from file", err)
				m.attemps++
				tx.Rollback()
				return err
			}
			i++
		}

		m.config.Log.Infof("FROM FILE inserted %d rows", i)

		err = os.Remove(m.config.FilePath + f.Name())
		if err != nil {
			tx.Rollback()
			return errors.New("Cannot delete file")

		}
		tx.Commit()

	}
	if len(files) > 0 {
		m.attemps = 0
	}

	return nil
}
