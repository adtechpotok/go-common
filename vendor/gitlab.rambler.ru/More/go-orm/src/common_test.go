package src

import "testing"

func TestFullTableName(t *testing.T) {
	v := FullTableName(testStuct{})
	if v != "1.2" {
		t.Error("Expected 1.2, got ", v)
	}
}

type testStuct struct{}

func (m testStuct) SchemaName() string { return "1" }
func (m testStuct) TableName() string  { return "2" }
