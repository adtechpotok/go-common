package dbwriter

import (
	"testing"
	"database/sql"
	"time"
	"github.com/sirupsen/logrus"
	"os"
	"github.com/pkg/errors"
	"fmt"
)

var transactionAttemp = 0

const testsqlDir = "sql_test/"

func TestMain(m *testing.M) {
	os.Mkdir(testsqlDir, 0777)
	file, _ := os.Create(testsqlDir + "deleted1.sql")
	file.Close()
	for i := 0; i < 200; i++ {
		file, _ := os.Create(testsqlDir + "test" + fmt.Sprintf("%d", i) + ".sql")
		file.Close()
	}
	v := m.Run()
	os.Remove(testsqlDir + "deleted1.sql")
	os.Remove(testsqlDir)
	os.Exit(v)
}

func TestAppend(t *testing.T) {
	shutdown := &testShutdownController{}
	config := WriteConfig{&testSql{}, logrus.New(), testsqlDir, 1, 10, 100,shutdown}
	m := New(config)
	for i := 0; i < 100; i++ {
		m.Append(testModel{1, "1"})
		time.Sleep(time.Millisecond)
	}
	go func(){
		time.Sleep(90*time.Millisecond)
		shutdown.turnOff = true
	}()
	time.Sleep(1 * time.Second)
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

type testShutdownController struct { turnOff bool}

func (*testShutdownController) Done() {}

func (m *testShutdownController) IsSwitchingOff() bool { return m.turnOff}


