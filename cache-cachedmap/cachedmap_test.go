package cachedmap

import (
	"encoding/json"
	"testing"

	"github.com/herb-go/herb/cache"
	_ "github.com/herb-go/herb/cache/drivers/freecache"
)

func newTestCache(ttl int64) *cache.Cache {
	config := json.RawMessage("{\"Size\": 10000000}")
	c := cache.New()
	err := c.Open("freecache", config, ttl)
	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	return c
}

type testmodel struct {
	Keyword string
	Content int
}

const valueKey = "valueKey"
const valueKeyAadditional = "valueKeyAadditional"
const valueKeyChanged = "valueKeyChanged"
const wrongDataKey = "wrongdata"

var WrongData = []byte("wrongdata")

const startValue = 1
const changedValue = 2
const mapCreatorKeyword = "mapCreatorKeyword"

const creatorKeyword = "creatorKeyword"

var rawData map[string]int

type testmodelmap map[string]*testmodel

func (m *testmodelmap) NewMapElement(key string) error {
	(*m)[key] = &testmodel{
		Keyword: mapCreatorKeyword,
		Content: 0,
	}
	return nil
}
func (m *testmodelmap) LoadMapElements(keys ...string) error {
	for _, v := range keys {
		(*m)[v].Content = rawData[v]
	}
	return nil
}
func creator(m *testmodelmap) func(key string) error {
	return func(key string) error {
		(*m)[key] = &testmodel{
			Keyword: creatorKeyword,
			Content: 0,
		}
		return nil
	}
}

func loader(m *testmodelmap) func(keys ...string) error {
	return func(keys ...string) error {
		for _, v := range keys {
			(*m)[v].Content = rawData[v]
		}
		return nil
	}
}
func TestMap(t *testing.T) {
	rawData = map[string]int{
		valueKey:            startValue,
		valueKeyAadditional: startValue,
		valueKeyChanged:     startValue,
	}
	c := newTestCache(-1)
	var err error
	var tm = testmodelmap{}
	err = LoadCachedMap(&tm, c, valueKey, valueKeyAadditional, valueKeyAadditional)
	if err != nil {
		t.Fatal(err)
	}

	if val := tm[valueKey].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKey].Keyword; val != mapCreatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Keyword; val != mapCreatorKeyword {
		t.Error(val)
	}
	delete(tm, valueKeyAadditional)
	rawData[valueKey] = changedValue
	rawData[valueKeyAadditional] = changedValue
	rawData[valueKeyChanged] = changedValue
	err = LoadCachedMap(&tm, c, valueKey, valueKeyAadditional, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}
	if val := tm[valueKey].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKey].Keyword; val != mapCreatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Keyword; val != mapCreatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyChanged].Content; val != changedValue {
		t.Error(val)
	}
	if val := tm[valueKeyChanged].Keyword; val != mapCreatorKeyword {
		t.Error(val)
	}
	c.SetBytesValue(wrongDataKey, WrongData, 0)
	err = LoadCachedMap(&tm, c, wrongDataKey)
	if err == nil {
		t.Fatal(err)
	}

}

func TestMapLoad(t *testing.T) {
	rawData = map[string]int{
		valueKey:            startValue,
		valueKeyAadditional: startValue,
		valueKeyChanged:     startValue,
	}
	c := newTestCache(-1)
	var err error
	var tm = testmodelmap{}
	err = Load(&tm, c, loader(&tm), creator(&tm), valueKey, valueKeyAadditional)
	if err != nil {
		t.Fatal(err)
	}
	if val := tm[valueKey].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKey].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	rawData[valueKey] = changedValue
	rawData[valueKeyAadditional] = changedValue
	rawData[valueKeyChanged] = changedValue
	err = Load(&tm, c, loader(&tm), creator(&tm), valueKeyAadditional, valueKeyChanged)
	if err != nil {
		t.Fatal(err)
	}
	if val := tm[valueKey].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKey].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Content; val != startValue {
		t.Error(val)
	}
	if val := tm[valueKeyAadditional].Keyword; val != creatorKeyword {
		t.Error(val)
	}
	if val := tm[valueKeyChanged].Content; val != changedValue {
		t.Error(val)
	}
	if val := tm[valueKeyChanged].Keyword; val != creatorKeyword {
		t.Error(val)
	}
}
