package src

import "testing"

// Тест покрывает инициализацию объекта и метод GetCache
func TestBaseSliceDbModelInit(t *testing.T) {
	v := BaseSlicedDbModel{}

	v.GetCache()
	if !v.isInit {
		t.Error("Base Slice Model should be inited before GetCache")
	}
}

//Тест покрывает инициализацию объекта и методы FindInCache, AddToCache
func TestBaseSlicedDbModelFindInCache(t *testing.T) {
	testCases := []testSlicePotok{
		{1, true, 1},
		{2, false, 1},
		{3, true, 2},
		{4, false, 1},
		{5, true, 2},
		{6, false, 1},
		{7, true, 2},
		{8, false, 1},
		{9, false, 2},
		{10, false, 1},
	}
	v := BaseSlicedDbModel{}

	res := v.FindInCache(1)
	if !v.isInit {
		t.Error("Base Slice Model should be inited before FindInCache")
	}

	if res != nil {
		t.Error("Not empty result")
	}
	for _, tc := range testCases {
		v.AddToCache(testSlicePotok{tc.relationKey, tc.active, tc.id})
		if !v.isInit {
			t.Error("Base Slice Model should be inited before AddToCache")
		}
		res = v.FindInCache(tc.id)
		if res == nil {
			t.Error("Empty result after put")
		}
	}

}

//Тест покрывает инициализацию объекта и методы FindInCache, AddToCache, ClearCache
func TestBaseSlicedDbModelClearCache(t *testing.T) {
	v := BaseSlicedDbModel{}

	v.AddToCache(testSlicePotok{1, true, 1})
	if !v.isInit {
		t.Error("Base Slice Model should be inited before AddToCache")
	}
	res := v.FindInCache(1)
	if res == nil {
		t.Error("Empty result after put")
	}
	res = v.FindInCache(2)
	if res != nil {
		t.Error("Found bad result")
	}

	v.ClearCache()
	if len(v.cache) != 0 {
		t.Error("Not empty result after clear")
	}
}

//Проверка работы функции Len()
func TestBaseSlicedDbModelLen(t *testing.T) {
	v := BaseSlicedDbModel{}

	v.AddToCache(testSlicePotok{1, true, 1})
	v.AddToCache(testSlicePotok{2, true, 2})
	v.AddToCache(testSlicePotok{3, true, 3})
	if !v.isInit {
		t.Error("Base Slice Model should be inited before AddToCache")
	}
	if 3 != v.Len() {
		t.Error("Len function dont't show actual length")
	}
	v.AddToCache(testSlicePotok{4, true, 4})
	v.AddToCache(testSlicePotok{3, true, 3})
	if 4 != v.Len() {
		t.Error("Len function dont't show actual length")
	}
}

//Проверка добавления двух значений с одинаковым relationKey
func TestBaseSlicedDbModelAddToCacheDuplicate(t *testing.T) {
	testCases := []testSlicePotok{
		{1, true, 1},
		{1, true, 1},
	}
	v := BaseSlicedDbModel{}

	for _, tc := range testCases {
		v.AddToCache(testSlicePotok{tc.relationKey, tc.active, tc.id})
	}
	if len(v.cache) != 1 {
		t.Error("Remove of duplicates crash")
	}
}

type testSlicePotok struct {
	relationKey int
	active      bool
	id          int
}

func (m testSlicePotok) GetRelationKey() int { return m.relationKey }
func (m testSlicePotok) IsActive() bool      { return m.active }
func (m testSlicePotok) GetId() int          { return m.id }
