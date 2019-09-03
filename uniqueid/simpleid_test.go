package uniqueid

import "testing"

func newSimpleIDGenerator() *Generator {
	g := NewGenerator()
	o := NewOptionConfigMap()
	o.Driver = "simpleid"
	o.Config.Set("Suff", "-test")
	err := o.ApplyTo(g)
	if err != nil {
		panic(err)
	}
	return g
}
func TestSimpleID(t *testing.T) {
	generator := newSimpleIDGenerator()
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

func BenchmarkSimpleID(b *testing.B) {
	generator := newSimpleIDGenerator()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			generator.GenerateID()
		}
	})
}
