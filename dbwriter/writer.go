package dbwriter

import (
	"database/sql"
	"fmt"
	orm "github.com/adtechpotok/go-orm"
	"github.com/adtechpotok/silog"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

// Config to writer
type WriteConfig struct {
	Db                *sql.DB         // Db instance
	Log               silog.Logger    // Log instance
	FilePath          string          // Path to write file
	ServerId          int             // Current server id
	TickTimeMs        time.Duration   // Tick for work
	MaxConnectTimeSec time.Duration   // Connect max time limit
	ShutdownControl   ShutdownControl // Shutdown switcher
	NoStart           bool
}

// Write sql to DB or file
type Writer struct {
	config      WriteConfig
	writeBuffer []orm.SchemaPotok
	mysqlBuffer []orm.SchemaPotok
	mutex       *sync.RWMutex
	attemps     int
}

// Return new writer instance
func New(config WriteConfig) *Writer {
	config.FilePath = strings.TrimRight(config.FilePath, "/") + "/"
	m := &Writer{config: config}
	m.mutex = &sync.RWMutex{}
	if !m.config.NoStart {
		go m.work()
	}
	return m
}

const attempsLimit = 2

// Append data to mysqlBuffer
func (m *Writer) Append(data ...orm.SchemaPotok) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.mysqlBuffer = append(m.mysqlBuffer, data...)
}

// Start writer work
func (m *Writer) work() {
	c := time.Tick(m.config.TickTimeMs * time.Millisecond)
	m.config.Db.SetConnMaxLifetime(time.Second * m.config.MaxConnectTimeSec)
	for range c {
		done := m.mysqlWrite()
		if done {
			return
		}
	}
}

// Write mysql to file or db
func (m *Writer) mysqlWrite() bool {
	// If there is data from last tick
	if len(m.writeBuffer) == 0 { // if data is empty, take it from current buffer
		m.mutex.Lock()
		m.writeBuffer = m.mysqlBuffer
		m.mysqlBuffer = make([]orm.SchemaPotok, 0)
		m.mutex.Unlock()
	}

	if m.config.Db.Ping() == nil { // if connection is active
		err := m.readFiles() // writing all data from files
		if err != nil {
			return false
		}
	}

	if len(m.writeBuffer) == 0 { // if there is no data - return
		if m.config.ShutdownControl.IsSwitchingOff() {
			m.config.ShutdownControl.Done()
			return true
		}
		return false
	}
	m.config.Log.Debugf("Got %d element for mysql write", len(m.writeBuffer))
	if m.attemps > attempsLimit { // if we have made more attemps than limited
		if err := m.fileWrite(); err != nil {
			m.config.Log.Error(err, "Cannot write to file.")
			return false
		}
		m.writeBuffer = make([]orm.SchemaPotok, 0)
	}

	if m.config.Db.Ping() == nil { // if connect is active
		i := 0
		for _, val := range m.getQueryData() {
			_, err := m.config.Db.Exec(val)
			if err != nil {
				m.config.Log.Error("Error while writing in db", err)
				m.attemps++
				m.writeBuffer = m.writeBuffer[i:]
				return false
			}
			i++
		}

		m.config.Log.Debugf("Inserted %d rows", i)
		m.writeBuffer = make([]orm.SchemaPotok, 0)
		m.attemps = 0
	} else {
		m.attemps++
	}
	return false
}

// Write sql to file
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

// Get sql from writeBuffer
func (m *Writer) getQueryData() []string {
	var res []string
	for _, val := range m.writeBuffer {
		k := makeInsertQuery(val, m.config.ServerId)
		res = append(res, k)
	}
	return res
}

// Read files with sql
func (m *Writer) readFiles() error {
	var deletedFileCounter int
	files, err := ioutil.ReadDir("./" + m.config.FilePath)
	if err != nil {
		m.config.Log.Error("Cannot read from file", err)
	}
	fileCounter := 0
	for _, f := range files {
		if fileCounter > 100 {
			return nil
		}
		if strings.Contains(f.Name(), "deleted") { //if transaction failed file would be marked deleted, we should skip it
			deletedFileCounter++
			continue
		}
		b, err := ioutil.ReadFile(m.config.FilePath + f.Name())
		if err != nil {
			m.config.Log.Error(err)
		}
		str := strings.Split(string(b), "\n")
		i := 0
		tx, _ := m.config.Db.Begin()
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

		m.config.Log.Debugf("FROM FILE inserted %d rows", i)

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
	if len(files)-deletedFileCounter > 0 {
		m.attemps = 0
	}

	return nil
}

type ShutdownControl interface {
	Done()
	IsSwitchingOff() bool
}
