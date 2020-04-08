package cachetree

import (
	"sort"
	"strings"
	"time"

	"github.com/herb-go/logger"

	"github.com/herb-go/herb/cache"
)

type children []*child

// Len is the number of elements in the collection.
func (c *children) Len() int {
	return len(*c)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (c *children) Less(i, j int) bool {
	return len((*c)[i].key) > len((*c)[j].key)
}

// Swap swaps the elements with indexes i and j.
func (c *children) Swap(i, j int) {
	v := (*c)[i]
	(*c)[i] = (*c)[j]
	(*c)[j] = v
}

type child struct {
	cache *cache.Cache
	key   string
}

func newChild(key string, cache *cache.Cache) *child {
	return &child{
		key:   key,
		cache: cache,
	}
}

type Tree struct {
	Debug bool
	*cache.Cache
	children *children
}

func NewTree() *Tree {
	c := children{}
	return &Tree{
		children: &c,
	}
}
func (t *Tree) find(key string) (string, *cache.Cache) {
	for k := range *(t.children) {
		if strings.HasPrefix(key, (*t.children)[k].key) {
			return key[len((*t.children)[k].key):], (*t.children)[k].cache
		}
	}
	return key, nil
}

func (t *Tree) SetUtil(u *cache.Util) {
	uc := u.Clone()
	uc.CollectionFactory = t.collectionFactory
	uc.NodeFactory = t.nodeFactory
	t.Driver.SetUtil(uc)
}

func (t *Tree) collectionFactory(c cache.Cacheable, key string, ttl time.Duration) *cache.Collection {
	s := c.FinalKey(key)
	if t.Debug {
		logger.Debug("cachetree collection created:" + s)
	}
	if k, d := t.find(s); d != nil {
		return cache.NewCollection(d, k, ttl)
	}
	return cache.DefaultCollectionFactory(c, key, ttl)
}
func (t *Tree) nodeFactory(c cache.Cacheable, key string) *cache.Node {
	s := c.FinalKey(key)
	if t.Debug {
		logger.Debug("cachetree node created:" + s)
	}
	if k, d := t.find(s); d != nil {
		return cache.NewNode(d, k)
	}
	return cache.DefaultNodeFactory(c, key)
}

type Config struct {
	Debug    bool
	Root     *cache.OptionConfig
	Children map[string]*cache.OptionConfig
}

func (c *Config) Create() (cache.Driver, error) {
	var err error
	d := NewTree()
	d.Debug = c.Debug
	root, err := cache.NewSubCache(c.Root)
	if err != nil {
		return nil, err
	}
	d.Cache = root
	for k := range c.Children {
		c, err := cache.NewSubCache(c.Children[k])
		if err != nil {
			return nil, err
		}
		*d.children = append(*d.children, newChild(cache.Key(k), c))
	}
	sort.Sort(d.children)
	return d, nil
}

func init() {
	cache.Register("tree", func(loader func(interface{}) error) (cache.Driver, error) {
		c := &Config{}
		err := loader(c)
		if err != nil {
			return nil, err
		}
		return c.Create()
	})
}
