package cachetree_test

import (
	"testing"
	"time"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herbconfig/loader"
	_ "github.com/herb-go/herbconfig/loader/drivers/jsonconfig"
)

func newTestCache(ttl int64) *cache.Cache {
	c := cache.New()
	oc := cache.NewOptionConfig()
	err := loader.LoadConfig("json", []byte(testConfig), oc)
	if err != nil {
		panic(err)
	}
	oc.TTL = ttl
	err = c.Init(oc)

	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	return c
}

func TestTree(t *testing.T) {
	c := newTestCache(1800)
	if c.DefualtTTL() != 1800*time.Second {
		t.Fatal(c)
	}
	node1 := c.Node("node1")
	if node1.DefualtTTL() != 1800*time.Second {
		t.Fatal(node1.DefualtTTL())
	}
	node2 := c.Node("test/test2/")
	if node2.DefualtTTL() != 2400*time.Second {
		t.Fatal(node2.DefualtTTL())
	}
	node3 := c.Node("test/test")
	if node3.DefualtTTL() != 3600*time.Second {
		t.Fatal(node3.DefualtTTL())
	}
	node4 := c.Node("alias/test")
	if node4.DefualtTTL() != 4800*time.Second {
		t.Fatal(node4.DefualtTTL())
	}
	c1 := c.Collection("node1")
	if c1.DefualtTTL() != 1800*time.Second {
		t.Fatal(c1.DefualtTTL())
	}
	c2 := c.Collection("test/test2/")
	if c2.DefualtTTL() != 2400*time.Second {
		t.Fatal(c2.DefualtTTL())
	}
	c3 := c.Collection("test/test")
	if c3.DefualtTTL() != 3600*time.Second {
		t.Fatal(c3.DefualtTTL())
	}
	c4 := c.Collection("alias/test")
	if c4.DefualtTTL() != 4800*time.Second {
		t.Fatal(c4.DefualtTTL())
	}
}
