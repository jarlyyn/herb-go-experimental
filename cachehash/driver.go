package cachehash

import (
	"time"

	"github.com/herb-go/herb/cache"
)

type Driver struct {
	cache.DriverUtil
	Store Store
}

//Set bytes data to cache by given key.
func (d *Driver) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {

}

//Update bytes data to cache by given key only if the cache exist.
func (d *Driver) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {

}

//Get bytes data from cache by given key.
func (d *Driver) GetBytesValue(key string) ([]byte, error) {

}

//Delete data in cache by given key.
func (d *Driver) Del(key string) error {

}

//Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
func (d *Driver) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {

}

//Set int val in cache by given key.Count cache and data cache are in two independent namespace.
func (d *Driver) SetCounter(key string, v int64, ttl time.Duration) error {

}

//Get int val from cache by given key.Count cache and data cache are in two independent namespace.
func (d *Driver) GetCounter(key string) (int64, error) {

}

//Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
func (d *Driver) DelCounter(key string) error {

}

//Set callback to handler error raised when gc.
func (d *Driver) SetGCErrHandler(f func(err error)) {

}
func (d *Driver) Expire(key string, ttl time.Duration) error {

}
func (d *Driver) ExpireCounter(key string, ttl time.Duration) error {

}
func (d *Driver) MGetBytesValue(keys ...string) (map[string][]byte, error) {

}
func (d *Driver) MSetBytesValue(map[string][]byte, time.Duration) error {

}

//Close cache.
func (d *Driver) Close() error {

}

//Delete all data in cache.
func (d *Driver) Flush() error {

}
