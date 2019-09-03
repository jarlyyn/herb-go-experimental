package uuid

import (
	"github.com/jarlyyn/herb-go-experimental/uniqueid"

	"testing"
)

func newUUIDGenerator() *uniqueid.Generator {
	g := uniqueid.NewGenerator()
	o := uniqueid.NewOptionConfigMap()
	o.Driver = "uuid"
	err := o.ApplyTo(g)
	if err != nil {
		panic(err)
	}
	return g
}

func newUUIDGeneratorV4() *uniqueid.Generator {
	g := uniqueid.NewGenerator()
	o := uniqueid.NewOptionConfigMap()
	o.Config.Set("Version", 4)
	o.Driver = "uuid"
	err := o.ApplyTo(g)
	if err != nil {
		panic(err)
	}
	return g
}

func TestUUID(t *testing.T) {
	generator := newUUIDGenerator()
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

		last = id
	}
}

func TestUUIDV4(t *testing.T) {
	generator := newUUIDGeneratorV4()
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
		last = id
	}
}

func BenchmarkUUID(b *testing.B) {
	generator := newUUIDGenerator()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			generator.GenerateID()
		}
	})
}

func BenchmarkUUIDV4(b *testing.B) {
	generator := newUUIDGeneratorV4()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			generator.GenerateID()
		}
	})
}
