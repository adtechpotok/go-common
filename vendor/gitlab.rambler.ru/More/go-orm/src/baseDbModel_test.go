package src

import "testing"

func TestBaseDbModelInit(t *testing.T) {
	var v BaseDbModel
	v = BaseDbModel{}
	v.GetCache()
	if !v.isInit {
		t.Error("Base Model should be inited before GetCache")
	}
}
func TestBaseDbModelGet(t *testing.T) {
	var v BaseDbModel

	v = BaseDbModel{}
	res := v.FindInCache(1)
	if !v.isInit {
		t.Error("Base Model should be inited before FindInCache")
	}

	if res != nil {
		t.Error("Not empty result")
	}

	v = BaseDbModel{}
	v.AddToCache(testPotok{true, 1})
	if !v.isInit {
		t.Error("Base Model should be inited before GetCache")
	}
	res = v.FindInCache(1)
	if res == nil {
		t.Error("Empty result after put")
	}

}
func TestBaseDbModelGetInactive(t *testing.T) {
	var v BaseDbModel

	v = BaseDbModel{}
	v.AddToCache(testPotok{false, 1})

	res := v.FindInCache(1)
	if res != nil {
		t.Error("Found inactive result")
	}
}

func TestBaseDbModelClear(t *testing.T) {
	var v BaseDbModel

	v = BaseDbModel{}
	v.AddToCache(testPotok{true, 1})

	v.ClearCache()
	res := v.FindInCache(1)

	if res != nil {
		t.Error("Not empty result after clear")
	}
}

type testPotok struct {
	active bool
	id     int
}

func (m testPotok) IsActive() bool { return m.active }
func (m testPotok) GetId() int     { return m.id }
