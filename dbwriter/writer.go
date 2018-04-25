package dbwriter

import (
	"time"
	"fmt"
	"os"
	"github.com/adtechpotok/silog"
	"io/ioutil"
	"strings"
	"github.com/pkg/errors"
	orm "github.com/adtechpotok/go-orm"
	"database/sql"
)

const attempsLimit = 2

type BaseWriter interface {
	mysqlWrite(conf *WriteConfig)
	setConnMaxLifetime(seconds time.Duration)
}

type WriteConfig struct {
	FilePath          string
	ServerId          int
	TickTimeMs        time.Duration
	MaxConnectTimeSec time.Duration
}

type Writer struct {
	Db          sql.DB
	Logger      silog.StandardLogger
	writeBuffer []orm.SchemaPotok
	mysqlBuffer []orm.SchemaPotok
	mutex       *sync.RWMutex
	attemps     int
}

func (m *Writer) setConnMaxLifetime(seconds time.Duration) {
	m.Db.SetConnMaxLifetime(time.Second * seconds)
}

func New(config WriteConfig) *Writer {
	m := &Writer{config: config}
	m.mutex = &sync.RWMutex{}
	go m.work(&config)

	return m
}

func (m *Writer) Append(data ...orm.SchemaPotok) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.mysqlBuffer = append(m.mysqlBuffer, data...)
}

func (m *Writer) work(config *WriteConfig) {
	c := time.Tick(m.config.TickTimeMs * time.Millisecond)
	m.config.Db.DB().SetConnMaxLifetime(time.Second * m.config.MaxConnectTimeSec)
	for range c {
		m.mysqlWrite(config)
	}
}

func (m *Writer) mysqlWrite(conf *WriteConfig) {

	//Смотрим остались ли данные с прошлой записи
	if len(m.writeBuffer) == 0 { // если данных нет, заполняем их из текущего буфера
		m.mutex.Lock()
		m.writeBuffer = m.mysqlBuffer
		m.mysqlBuffer = make([]orm.SchemaPotok, 0)
		m.mutex.Unlock()
	}

	if m.Db.Ping() == nil { //если коннект есть
		err := m.readFiles() //записываем все файлы
}
		if err != nil {
			return
		}
	}

	if len(m.writeBuffer) == 0 { // если данных нет прекращаем работу
		return
	}
  
	m.Logger.Infof("Got %d element for mysql write", len(m.writeBuffer))
	if m.attemps > attempsLimit { //используем попытки для понимания пишем мы в файл или в базу
		if err := m.fileWrite(conf); err != nil {
			m.Logger.Error(err, "Cannot write to file.")
			return
		}
		m.writeBuffer = make([]orm.SchemaPotok, 0)
	}

	if m.Db.Ping() == nil { //если конект есть
		i := 0
		for _, val := range m.getQueryData() {
			_, err := m.Db.Exec(val)
			if err != nil {
				m.Logger.Error("Error while writing in db", err)
				m.attemps++
				m.writeBuffer = m.writeBuffer[i:]
				return
			}
			i++
		}

		m.Logger.Infof("Inserted %d rows", i)
		m.writeBuffer = make([]orm.SchemaPotok, 0)
		m.attemps = 0
	} else {
		m.attemps++
	}

}


func (m *Writer) getQueryData() []string {
	var res []string
	for _, val := range m.writeBuffer {
		k := makeInsertQuery(val)
		res = append(res, k)
	}
	return res
}

func (m *Writer) fileWrite(conf *WriteConfig) error {
	fileName := fmt.Sprintf("%s%d.sql", m.config.FilePath, time.Now().UnixNano())
	file, err := os.Create(fileName)
	defer file.Close()

	if err != nil {
		return err
	}

	file.WriteString(strings.Join(m.getQueryData(), "\n"))
	return nil
}

func (m *Writer) readFiles(conf *WriteConfig) error {
	files, err := ioutil.ReadDir("./" + conf.FilePath)
	if err != nil {
		m.Logger.Error("Cannot read from file", err)
	}
	fileCounter := 0
	for _, f := range files {
		if fileCounter > 100 {
			return nil
		}
		b, err := ioutil.ReadFile(m.config.FilePath + f.Name())
		if err != nil {
			m.Logger.Error(err)
		}
		str := strings.Split(string(b), "\n")
		i := 0
		tx, _ := m.Db.Begin()
		for _, val := range str {
			_, err := tx.Exec(val)
			if err != nil {
				m.Logger.Error("Error while writing in db from file", err)
				m.attemps++
				tx.Rollback()
				return err
			}
			i++
		}

		m.Logger.Infof("FROM FILE inserted %d rows", i)

		deletedFileName := m.config.FilePath + "deleted" + f.Name()
		err = os.Rename(m.config.FilePath+f.Name(), deletedFileName)
		if err != nil {
			tx.Rollback()
			return errors.New("Cannot delete file")
		}
		err = tx.Commit()
		if err == nil {
			err = os.Remove(deletedFileName)
			if err != nil {
				m.config.Log.Errorf("Cannot delete file %s", deletedFileName)

			}
		} else {
			os.Rename(deletedFileName, m.config.FilePath+f.Name())
			m.config.Log.Errorf("%v transaction err", err)
		}
		fileCounter++
	}
	if len(files) > 0 {
		m.attemps = 0
	}

	return nil
}
