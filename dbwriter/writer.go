package viewwriter

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
	attemps     int
}

func (m *Writer) setConnMaxLifetime(seconds time.Duration) {
	m.Db.SetConnMaxLifetime(time.Second * seconds)
}

func (m *Writer) mysqlWrite(conf *WriteConfig) {

	//Смотрим остались ли данные с прошлой записи
	if len(m.writeBuffer) == 0 { // если данных нет, заполняем их из текущего буфера
		mutex.Lock()
		m.writeBuffer = mysqlBuffer
		mysqlBuffer = make([]orm.SchemaPotok, 0)
		mutex.Unlock()
	}

	if m.Db.Ping() == nil { //если конект есть
		err := m.readFiles(conf) //записываем все файлы
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

	fileName := fmt.Sprintf("%s%d.sql", conf.FilePath, time.Now().UnixNano())
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

	for _, f := range files {
		b, err := ioutil.ReadFile(conf.FilePath + f.Name())
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

		err = os.Remove(conf.FilePath + f.Name())
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
