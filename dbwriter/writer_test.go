package dbwriter

import (
	"testing"
	"time"
	"github.com/sirupsen/logrus"
	"os"
	"fmt"
	"io/ioutil"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/pkg/errors"
)

const testSqlDir = "sql_test/"

func TestMain(m *testing.M) {
	os.Mkdir(testSqlDir, 0777)
	file, _ := os.Create(testSqlDir + "deleted1.sql")
	file.Close()
	for i := 0; i < 200; i++ {
		file, _ := os.Create(fmt.Sprintf("%stest%d.sql", testSqlDir, i))
		file.Close()
	}
	v := m.Run()
	os.Remove(testSqlDir + "deleted1.sql")
	os.Remove(testSqlDir)
	os.Exit(v)
}

func TestAppend(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	shutdown := &testShutdownController{}
	logInstance := logrus.New()
	logInstance.Out = ioutil.Discard
	/** Первая запись в файл всегда ошибочка */
	mock.ExpectBegin()
	mock.ExpectExec("").WillReturnError(errors.New("test"))
	mock.ExpectRollback()

	/** Для каждого нашего запроса, содаем успешный  */
	for i := 0; i < 200; i++ {
		sucessTransaction(mock)
	}

	config := WriteConfig{Db: db, Log: logInstance, FilePath: testSqlDir, ServerId: 1, TickTimeMs: 10, MaxConnectTimeSec: 100, ShutdownControl: shutdown}
	m := New(config)
	for i := 0; i < 105; i++ {
		/* записываем успешно все данные */
		mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
		m.Append(testModel{1, "1"})
		time.Sleep(time.Millisecond)
	}

	/* проверяем кейс безуспешной записи, которая удается при восстанавление данных из файла*/
	m.Append(testModel{1, "1"})
	mock.ExpectExec("").WillReturnError(errors.New("1"))
	mock.ExpectExec("").WillReturnError(errors.New("1"))

	sucessTransaction(mock)

	go func() {
		time.Sleep(90 * time.Millisecond)
		shutdown.turnOff = true
	}()

	time.Sleep(200 * time.Millisecond)
}

type testModel struct {
	Id     int    `gorm:"column:id"`
	String string `gorm:"column:string;"`
}

func (testModel) SchemaName() string { return "test" }

func (testModel) TableName() string { return "test" }

type testShutdownController struct{ turnOff bool }

func (*testShutdownController) Done()                  {}
func (m *testShutdownController) IsSwitchingOff() bool { return m.turnOff }

func sucessTransaction(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}
