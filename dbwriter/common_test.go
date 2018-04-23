package dbwriter

import (
	"testing"
	"time"
)

func TestParseTagSetting(t *testing.T) {
	phoneSlice := make([]int, 2)
	phoneSlice[0] = 1
	phoneSlice[1] = 2
	m := testStruct{
		1,
		phoneSlice,
		1,
		time.Now(),
		true,
		false,
		"23",
		0.0,
		"23",
		0,
	}

	stringSql := makeInsertQuery(m,1)

	if stringSql == "" {
		t.Error("All bad")
	}
}

type testStruct struct {
	Id           int       `gorm:"column:id"`
	Json         []int     `gorm:"column:json; fromJson"`
	ServerId     int       `gorm:"column:server_id"`
	Created      time.Time `gorm:"column:created"`
	BoolTrue     bool      `gorm:"column:bool_true"`
	BoolFalse    bool      `gorm:"column:bool_false"`
	String       string    `gorm:"column:string;"`
	Float        float64   `gorm:"column:float"`
	StringQuotes string    `gorm:"column:stringQuotes; clearQuotes"`
	noParse      int
}

func (m testStruct) SchemaName() string {
	return "Voximplant"
}
func (m testStruct) TableName() string {
	return "Test"
}

func (m testStruct) AfterSql() string {
	return ""
}

func TestClearQuotes(t *testing.T) {
	a := `t\nest\`
	b := clearQuotes(a)

	if b != `t\\nest\\` {
		t.Error("Quote was not cleared")
	}
}
