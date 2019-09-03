package snowflake

import (
	"github.com/jarlyyn/herb-go-experimental/uniqueid"

	"testing"
)

func newSnowFlakeGenerator() *uniqueid.Generator {
	g := uniqueid.NewGenerator()
	o := uniqueid.NewOptionConfigMap()
	o.Driver = "snowflake"
	err := o.ApplyTo(g)
	if err != nil {
		panic(err)
	}
	return g
}

func TestSnowFlake(t *testing.T) {
	generator := newSnowFlakeGenerator()
	var last = ""
	var usedmap = map[string]bool{}
	for i := 0; i < 1000; i++ {
		id, err := generator.GenerateID()
		if err != nil {
			t.Fatal(err)
		}
		if usedmap[id] {
			t.Fatal(id)
		}
		usedmap[id] = true
		if last == id {
			t.Fatal(id)
		}
		if last >= id {
			t.Fatal(id)
		}
		last = id
	}
}

func BenchmarkSnowFlake(b *testing.B) {
	generator := newSnowFlakeGenerator()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			generator.GenerateID()
		}
	})
}
