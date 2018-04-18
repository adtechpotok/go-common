package src

import "testing"

// Тест покрывает инициализацию объекта и метод GetCache
func TestStringDbModelInit(t *testing.T) {
	v := StringDbModel{}

	v.GetCache()
	if !v.isInit {
		t.Error("String Db Model should be inited before GetCache")
	}
}

//Тест покрывает инициализацию объекта и методы FindInCache, AddToCache
func TestStringDbModelFindInCache(t *testing.T) {
	testCases := []testStringDbModelPotok{
		{true, "1"},
		{false, "1"},
		{true, "2"},
		{false, "2"},
	}
	v := StringDbModel{}

	res := v.FindInCache("1")
	if !v.isInit {
		t.Error("String Db Model should be inited before FindInCache")
	}

	if res != nil {
		t.Error("Not empty result")
	}
	for _, tc := range testCases {
		v.AddToCache(testStringDbModelPotok{tc.active, tc.id})
		if !v.isInit {
			t.Error("String Db Model should be inited before AddToCache")
		}
		res = v.FindInCache(tc.id)

		if tc.IsActive() && res == nil {
			t.Error("Empty result after put")
		}

		if !tc.IsActive() && res != nil {
			t.Error("Not active result found in cache")
		}
	}

}

//Тест покрывает инициализацию объекта и методы FindInCache, AddToCache, ClearCache
func TestStringDbModelClearCache(t *testing.T) {
	v := StringDbModel{}

	v.AddToCache(testStringDbModelPotok{true, "1"})
	if !v.isInit {
		t.Error("String Db Model should be inited before AddToCache")
	}
	res := v.FindInCache("1")
	if res == nil {
		t.Error("Empty result after put")
	}

	res = v.FindInCache("2")
	if res != nil {
		t.Error("Found bad result")
	}

	v.ClearCache()
	if v.Len() != 0 {
		t.Error("Not empty result after clear")
	}
}

//Проверка добавления двух значений с одинаковым id
func TestStringDbModelAddToCacheDuplicate(t *testing.T) {
	v := StringDbModel{}

	testCases := []testStringDbModelPotok{
		{true, "1"},
		{true, "1"},
	}
	for _, tc := range testCases {
		v.AddToCache(testStringDbModelPotok{tc.active, tc.id})
	}
	if v.Len() != 1 {
		t.Error("Remove of duplicates crash")
	}
}

type testStringDbModelPotok struct {
	active bool
	id     string
}

func (m testStringDbModelPotok) IsActive() bool { return m.active }
func (m testStringDbModelPotok) GetId() string  { return m.id }

func (m testStringDbModelPotok) OutDated() bool    { return false }
