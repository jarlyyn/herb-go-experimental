package cachedmap

import (
	"reflect"

	"github.com/herb-go/herb/cache"
)

type CachedMap interface {
	NewMapElement(string) error
	LoadMapElements(keys ...string) error
}

func unmarshalMapElement(mapvalue reflect.Value, creator func(string) error, key string, data []byte) (err error) {
	err = creator(key)
	if err != nil {
		return err
	}
	var v = mapvalue.MapIndex(reflect.ValueOf(key)).Interface()
	err = cache.UnmarshalMsgpack(data, &v)
	return err
}
func Load(cm interface{}, c cache.Cacheable, loader func(keys ...string) error, creator func(string) error, keys ...string) error {
	var keysmap = make(map[string]bool, len(keys))
	var mapvalue = reflect.Indirect(reflect.ValueOf(cm))
	var filteredKeys = make([]string, len(keys))
	var filteredKeysLength = 0

	for k := range keys {
		if keysmap[keys[k]] == true {
			continue
		}
		keysmap[keys[k]] = true
		if !mapvalue.MapIndex(reflect.ValueOf(keys[k])).IsValid() {

			filteredKeys[filteredKeysLength] = keys[k]
			filteredKeysLength++
		}
	}
	filteredKeys = filteredKeys[:filteredKeysLength]
	results, err := c.MGetBytesValue(filteredKeys...)
	if err != nil {
		return err
	}
	var uncachedKeys = make([]string, len(results))
	var uncachedKeysLength = 0
	for k := range results {
		if results[k] == nil {
			err = creator(k)
			if err != nil {
				return err
			}
			uncachedKeys[uncachedKeysLength] = k
			uncachedKeysLength++
		} else {
			err = unmarshalMapElement(mapvalue, creator, k, results[k])
			if err != nil {
				return err
			}
		}
	}
	uncachedKeys = uncachedKeys[:uncachedKeysLength]
	err = loader(uncachedKeys...)
	if err != nil {
		return err
	}
	var data = make(map[string][]byte, len(uncachedKeys))
	for k := range uncachedKeys {
		v := mapvalue.MapIndex(reflect.ValueOf(uncachedKeys[k])).Interface()
		data[uncachedKeys[k]], err = cache.MarshalMsgpack(v)
		if err != nil {
			return err
		}
	}
	return c.MSetBytesValue(data, 0)
}
func LoadCachedMap(cm CachedMap, c cache.Cacheable, keys ...string) error {
	return Load(cm, c, cm.LoadMapElements, cm.NewMapElement, keys...)
}
