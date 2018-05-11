package dbwriter

import (
	"testing"
	"database/sql"
	"time"
	"github.com/sirupsen/logrus"
	"os"
	"github.com/pkg/errors"
	"fmt"
	"io/ioutil"
)

var transactionAttemp = 0

const testSqlDir = "sql_test/"

func TestMain(m *testing.M) {
	os.Mkdir(testSqlDir, 0777)
	file, _ := os.Create(testSqlDir + "deleted1.sql")
	file.Close()
	for i := 0; i < 200; i++ {
		file, _ := os.Create(fmt.Sprintf("%stest%d.sql",testSqlDir, i))
		file.Close()
	}
	v := m.Run()
	os.Remove(testSqlDir + "deleted1.sql")
	os.Remove(testSqlDir)
	os.Exit(v)
}

func TestAppend(t *testing.T) {
	shutdown := &testShutdownController{}
	logInstance := logrus.New()
	logInstance.Out = ioutil.Discard
	config := WriteConfig{&testSql{}, logInstance, testSqlDir, 1, 10, 100, shutdown}
	m := New(config)
	for i := 0; i < 105; i++ {
		m.Append(testModel{1, "1"})
		time.Sleep(time.Millisecond)
	}
	go func() {
		time.Sleep(90 * time.Millisecond)
		shutdown.turnOff = true
	}()

	time.Sleep(200 * time.Millisecond)
}

type testSql struct {
	pingAttempts int
	execAttempts int
}

func (m *testSql) Begin() (transactionInterface, error) { return &testTransaction{}, nil }

func (m *testSql) Exec(query string, args ...interface{}) (sql.Result, error) {
	m.execAttempts++
	if m.execAttempts > 1 {
		return nil, nil
	} else {
		return nil, errors.New("Transaction error")
	}
}

func (m *testSql) Ping() error {
	m.pingAttempts++
	if m.pingAttempts > 20 {
		return nil
	} else {
		return errors.New("db error")
	}
}

func (m *testSql) SetConnMaxLifetime(duration time.Duration) {}

type testModel struct {
	Id     int    `gorm:"column:id"`
	String string `gorm:"column:string;"`
}

func (testModel) SchemaName() string { return "test" }

func (testModel) TableName() string { return "test" }

type testTransaction struct{}

func (*testTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	transactionAttemp++
	if transactionAttemp > 2 {
		return nil, nil
	} else {
		return nil, errors.New("Transaction error")
	}
}

func (*testTransaction) Rollback() error {
	return nil
}

func (m *testTransaction) Commit() error {
	transactionAttemp++
	if transactionAttemp > 1 {
		return nil
	} else {
		return errors.New("Transaction error")
	}
}

type testShutdownController struct{ turnOff bool }

func (*testShutdownController) Done() {}

func (m *testShutdownController) IsSwitchingOff() bool { return m.turnOff }
