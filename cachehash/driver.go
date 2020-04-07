package cachehash

import (
	"encoding/binary"
	"time"

	"github.com/herb-go/herb/cache"
)

type Driver struct {
	cache.DriverUtil
	Store        Store
	GcErrHanlder func(err error)
}

type context struct {
	hash     string
	unlocker func()
	data     *Hash
}

func (d *Driver) lockAndGetData(key string) (ctx *context, err error) {
	c := &context{}
	c.hash, err = d.Store.Hash(key)
	if err != nil {
		return nil, err
	}
	c.unlocker, err = d.Store.Lock(c.hash)
	if err != nil {
		return nil, err
	}
	c.data, err = d.Store.Load(c.hash)
	if err != nil {
		c.unlocker()
		return nil, err
	}
	return c, nil
}
func (d *Driver) save(ctx *context, status *Status) error {
	if ctx.data.isEmpty() {
		return d.Store.Delete(ctx.hash)
	}
	return d.Store.Save(ctx.hash, status, ctx.data)
}

//Set bytes data to cache by given key.
func (d *Driver) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.set(NewData(key, now.Add(ttl).Unix(), bytes), now.Unix())
	return d.save(ctx, status)
}

//Update bytes data to cache by given key only if the cache exist.
func (d *Driver) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.update(NewData(key, now.Add(ttl).Unix(), bytes), now.Unix())
	return d.save(ctx, status)

}

//Get bytes data from cache by given key.
func (d *Driver) GetBytesValue(key string) ([]byte, error) {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return nil, err
	}
	data := ctx.data.get(key, now.Unix())
	if data == nil {
		return nil, cache.ErrNotFound
	}
	return data.Data, nil
}

//Delete data in cache by given key.
func (d *Driver) Del(key string) error {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.delete(key, now.Unix())
	return d.save(ctx, status)
}

//Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
func (d *Driver) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	var v int64
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return 0, err
	}
	data := ctx.data.get(key, now.Unix())
	if data == nil {
		v = 0
	} else {
		v = int64(binary.BigEndian.Uint64(data.Data))
	}
	v = v + increment
	var bytes = make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(v))
	status := ctx.data.set(NewData(key, now.Add(ttl).Unix(), bytes), now.Unix())
	err = d.save(ctx, status)
	if err != nil {
		return 0, err
	}
	return v, nil
}

//Set int val in cache by given key.Count cache and data cache are in two independent namespace.
func (d *Driver) SetCounter(key string, v int64, ttl time.Duration) error {
	var bytes = make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(v))
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.set(NewData(key, now.Add(ttl).Unix(), bytes), now.Unix())
	return d.save(ctx, status)
}

//Get int val from cache by given key.Count cache and data cache are in two independent namespace.
func (d *Driver) GetCounter(key string) (int64, error) {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return 0, err
	}
	data := ctx.data.get(key, now.Unix())
	if data == nil {
		return 0, nil
	}
	v := binary.BigEndian.Uint64(data.Data)
	return int64(v), nil

}

//Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
func (d *Driver) DelCounter(key string) error {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.delete(key, now.Unix())
	return d.save(ctx, status)
}

//Set callback to handler error raised when gc.
func (d *Driver) SetGCErrHandler(f func(err error)) {
	d.GcErrHanlder = f
}
func (d *Driver) Expire(key string, ttl time.Duration) error {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.expired(key, now.Add(ttl).Unix(), now.Unix())
	return d.save(ctx, status)
}
func (d *Driver) ExpireCounter(key string, ttl time.Duration) error {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.expired(key, now.Add(ttl).Unix(), now.Unix())
	return d.save(ctx, status)
}
func (d *Driver) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	var result = map[string][]byte{}
	for k := range keys {
		bs, err := d.GetBytesValue(keys[k])
		if err != nil {
			return nil, err
		}
		result[keys[k]] = bs
	}
	return result, nil
}
func (d *Driver) MSetBytesValue(data map[string][]byte, ttl time.Duration) error {
	for key := range data {
		err := d.SetBytesValue(key, data[key], ttl)
		if err != nil {
			return err
		}
	}
	return nil
}

//Close cache.
func (d *Driver) Close() error {
	return d.Store.Close()
}

//Delete all data in cache.
func (d *Driver) Flush() error {
	return d.Store.Flush()
}
