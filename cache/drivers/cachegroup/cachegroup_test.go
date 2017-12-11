package cachegroup

import (
	"testing"

	"github.com/herb-go/herb/cache"

	"bytes"
	"encoding/json"
	"time"

	_ "github.com/herb-go/herb/cache/drivers/freecache"
)

func newTestCache(ttl int64) *cache.Cache {
	c := cache.New()
	err := c.Open("cachegroup", json.RawMessage(testConfig), ttl)
	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	return c
}

func TestNameConflict(t *testing.T) {
	var err error
	defaultTTL := int64(1)
	testKey := "testKey"
	testDataModel := "test"
	testDataBytes := []byte("testbytes")
	testDataInt := int64(12345)
	var resultDataModel string
	var resultDataBytes []byte
	var resultInt int64
	c := newTestCache(defaultTTL)
	err = c.Set(testKey, testDataModel, cache.DefualtTTL)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKey, &resultDataModel)
	if err != nil {
		t.Fatal(err)
	}
	if resultDataModel != testDataModel {
		t.Errorf("Cache get result error %s", resultDataModel)
	}
	err = c.SetCounter(testKey, testDataInt, cache.DefualtTTL)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKey, &resultDataModel)
	if err != nil {
		t.Fatal(err)
	}
	if resultDataModel != testDataModel {
		t.Errorf("Cache get result error %s", resultDataModel)
	}
	resultInt, err = c.GetCounter(testKey)
	if err != nil {
		t.Fatal(err)
	}
	if resultInt != testDataInt {
		t.Errorf("Cache getCounter result error %d", testDataInt)
	}
	err = c.SetBytesValue(testKey, testDataBytes, cache.DefualtTTL)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKey, &resultDataModel)
	if err == nil && testDataModel == resultDataModel {
		t.Fatal(err)
	}
	resultDataBytes, err = c.GetBytesValue(testKey)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(resultDataBytes, testDataBytes) != 0 {
		t.Errorf("Cache get result error %s", resultDataModel)
	}
	resultInt, err = c.GetCounter(testKey)
	if err != nil {
		t.Fatal(err)
	}
	if resultInt != testDataInt {
		t.Errorf("Cache getCounter result error %d", testDataInt)
	}

}

func TestCloseAndFlush(t *testing.T) {
	defaultTTL := int64(1)
	testKey := "testKey"
	testDataModel := "test"
	var resultDataModel string
	c := newTestCache(defaultTTL)
	err := c.Set(testKey, testDataModel, cache.DefualtTTL)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKey, &resultDataModel)
	if err != nil {
		t.Fatal(err)
	}
	if resultDataModel != testDataModel {
		t.Errorf("Cache get result error %s", resultDataModel)
	}
	err = c.Flush()
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKey, &resultDataModel)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}

	err = c.Close()
	if err != nil {
		t.Fatal(err)
	}
}
func TestSearchUpdate(t *testing.T) {
	var err error
	defaultTTL := int64(1)
	c := newTestCache(defaultTTL)
	defer c.Close()
	testKey := "testkey"
	testKeyUpdate := "testkeyupdate"
	testKeyBytes := "testkeybytes"
	testKeyBytesUpdate := "testkeybytesupdate"
	testDataModel := "test"
	var resultDataModel string
	testDataBytes := []byte("testbytes")
	err = c.Set(testKey, testDataModel, cache.TTLForever)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKey, &resultDataModel)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Set(testKeyUpdate, testDataModel, cache.TTLForever)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKeyUpdate, &resultDataModel)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Del(testKeyUpdate)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Update(testKeyUpdate, testDataModel, cache.TTLForever)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKeyUpdate, &resultDataModel)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}

	err = c.SetBytesValue(testKeyBytes, testDataBytes, cache.TTLForever)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyBytes)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetBytesValue(testKeyBytesUpdate, testDataBytes, cache.TTLForever)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyBytesUpdate)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Del(testKeyBytesUpdate)
	if err != nil {
		t.Fatal(err)
	}
	err = c.UpdateBytesValue(testKeyBytesUpdate, testDataBytes, cache.TTLForever)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyBytesUpdate)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
}
func TestSearchByPrefix(t *testing.T) {
	defaultTTL := int64(1)
	var testPrefix = "prefix"
	var testPrefix2 = "prefix"
	c := newTestCache(defaultTTL)
	defer c.Close()
	_, err := c.SearchByPrefix(testPrefix)
	if err != cache.ErrFeatureNotSupported {
		t.Fatal(err)
	}
	_, err = c.SearchCounterByPrefix(testPrefix2)
	if err != cache.ErrFeatureNotSupported {
		t.Fatal(err)
	}
}
func TestCounter(t *testing.T) {
	defaultTTL := int64(1)
	testKey := "testKey"
	testInitVal := int64(1)
	testIncremeant := int64(2)
	testTargetResultInt := int64(3)
	var resultDataInt int64
	c := newTestCache(defaultTTL)
	err := c.SetCounter(testKey, testInitVal, cache.DefualtTTL)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	resultDataInt, err = c.GetCounter(testKey)
	if err != nil {
		t.Fatal(err)
	}
	if resultDataInt != testInitVal {
		t.Errorf("GetCounter error %d ", resultDataInt)
	}
	resultDataInt, err = c.IncrCounter(testKey, testIncremeant, cache.DefualtTTL)
	if err != nil {
		t.Fatal(err)
	}
	if resultDataInt != testTargetResultInt {
		t.Errorf("IncrCounter error %d ", resultDataInt)
	}
	resultDataInt, err = c.GetCounter(testKey)
	if err != nil {
		t.Fatal(err)
	}
	if resultDataInt != testTargetResultInt {
		t.Errorf("GetCounter error %d ", resultDataInt)
	}
}
func TestDefaulTTL(t *testing.T) {
	defaultTTL := int64(1)
	testKey := "testKey"
	testKey2 := "testKey2"
	testKey3 := "testKey3"
	testDataModel := "test"
	var resultDataModel string
	testDataBytes := []byte("testbytes")
	var resultDataBytes []byte
	testDataInt := int64(1)
	var resultDataInt int64
	c := newTestCache(defaultTTL)
	defer c.Close()
	err := c.Set(testKey, testDataModel, cache.DefualtTTL)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetBytesValue(testKey2, testDataBytes, cache.DefualtTTL)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetCounter(testKey3, testDataInt, cache.DefualtTTL)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKey, &resultDataModel)
	if err != nil {
		t.Fatal(err)
	}
	if resultDataModel != testDataModel {
		t.Errorf("Cache get result error %s", resultDataModel)
	}
	resultDataBytes, err = c.GetBytesValue(testKey2)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(resultDataBytes, testDataBytes) != 0 {
		t.Errorf("Cache getBytesValue result error %s", string(resultDataBytes))
	}
	resultDataInt, err = c.GetCounter(testKey3)
	if err != nil {
		t.Fatal(err)
	}
	if resultDataInt != testDataInt {
		t.Errorf("Cache get result error %d", testDataInt)
	}
	time.Sleep(2000 * time.Millisecond)
	err = c.Get(testKey, &resultDataModel)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKey2)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	resultDataInt, err = c.GetCounter(testKey3)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
}
func TestTTL(t *testing.T) {
	var err error
	defaultTTL := int64(-1)
	c := newTestCache(defaultTTL)
	defer c.Close()

	testKeyTTLForver := "forever"
	testKeyTTLForverBytes := "foreverbytes"
	testKeyTTLForverCounter := "forevercounter"
	testKeyTTL1Second := "1second"
	testKeyTTL1SecondBytes := "1secondbytes"
	testKeyTTL1SecondCounter := "1secondcounter"
	testKeyTTL3Second := "3second"
	testKeyTTL3SecondBytes := "3secondbytes"
	testKeyTTL3SecondCounter := "3secondcounter"
	testKeyTTLRefresh := "refresh"
	testKeyTTLRefreshBytes := "refreshbytes"
	testKeyTTLRefreshCounter := "refreshcounter"
	testKeyTTLExpire := "expire"
	testKeyTTLExpireBytes := "expirebytes"
	testKeyTTLExpireCounter := "expirecounter"

	testDataModel := "12345"
	testDataBytes := []byte("12345byte")
	testDataInt := int64(99999)
	var resultModelData string
	err = c.Set(testKeyTTLForver, testDataModel, -1)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetBytesValue(testKeyTTLForverBytes, testDataBytes, -1)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetCounter(testKeyTTLForverCounter, testDataInt, -1)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Set(testKeyTTL1Second, testDataModel, 1*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetBytesValue(testKeyTTL1SecondBytes, testDataBytes, 1*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetCounter(testKeyTTL1SecondCounter, testDataInt, 1*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Set(testKeyTTL3Second, testDataModel, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetBytesValue(testKeyTTL3SecondBytes, testDataBytes, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetCounter(testKeyTTL3SecondCounter, testDataInt, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Set(testKeyTTLRefresh, testDataModel, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetBytesValue(testKeyTTLRefreshBytes, testDataBytes, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetCounter(testKeyTTLRefreshCounter, testDataInt, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Set(testKeyTTLExpire, testDataModel, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetBytesValue(testKeyTTLExpireBytes, testDataBytes, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetCounter(testKeyTTLExpireCounter, testDataInt, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Get(testKeyTTLForver, &resultModelData)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.GetBytesValue(testKeyTTLForverBytes)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTLForverCounter)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKeyTTL1Second, &resultModelData)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyTTL1SecondBytes)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTL1SecondCounter)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKeyTTL3Second, &resultModelData)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyTTL3SecondBytes)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTL3SecondCounter)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKeyTTLRefresh, &resultModelData)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyTTLRefreshBytes)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTLRefreshCounter)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKeyTTLExpire, &resultModelData)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyTTLExpireBytes)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTLExpireCounter)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2000 * time.Millisecond)
	err = c.Get(testKeyTTLForver, &resultModelData)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyTTLForverBytes)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTLForverCounter)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKeyTTL1Second, &resultModelData)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyTTL1SecondBytes)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTL1SecondCounter)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	err = c.Get(testKeyTTL3Second, &resultModelData)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyTTL3SecondBytes)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTL3SecondCounter)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKeyTTLRefresh, &resultModelData)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyTTLRefreshBytes)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTLRefreshCounter)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Set(testKeyTTLRefresh, testDataModel, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetBytesValue(testKeyTTLRefreshBytes, testDataBytes, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetCounter(testKeyTTLRefreshCounter, testDataInt, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Expire(testKeyTTLExpire, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Expire(testKeyTTLExpireBytes, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	err = c.ExpireCounter(testKeyTTLExpireCounter, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2000 * time.Millisecond)
	err = c.Get(testKeyTTLForver, &resultModelData)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyTTLForverBytes)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTLForverCounter)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get(testKeyTTL1Second, &resultModelData)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyTTL1SecondBytes)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTL1SecondCounter)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	err = c.Get(testKeyTTL3Second, &resultModelData)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyTTL3SecondBytes)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTL3SecondCounter)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	err = c.Get(testKeyTTLRefresh, &resultModelData)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyTTLRefreshBytes)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTLRefreshCounter)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Get(testKeyTTLExpire, &resultModelData)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue(testKeyTTLExpireBytes)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetCounter(testKeyTTLExpireCounter)
	if err != nil {
		t.Fatal(err)
	}
}
