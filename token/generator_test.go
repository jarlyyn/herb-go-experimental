package token

import "testing"

func TestListGenerator(t *testing.T) {
	list := []byte("abcde")
	listmap := map[byte]bool{}
	for _, v := range list {
		listmap[v] = true
	}
	g := &ListGenerator{
		List: list,
		Min:  256,
	}
	s, err := g.Generate()
	if err != nil {
		panic(err)
	}
	if len(s) != 256 {
		t.Fatal(s)
	}
	for _, v := range s {
		if listmap[v] == false {
			t.Fatal(s)
		}
	}
	g = &ListGenerator{
		List: list,
		Min:  15,
		Max:  256,
	}
	s, err = g.Generate()
	if err != nil {
		panic(err)
	}
	if len(s) > 256 || len(s) < 15 {
		t.Fatal(s)
	}

}

func TestBytesGenerator(t *testing.T) {
	token := New()
	token.ID = "test"
	err := Regenerate(BytesGenerator(15), token)
	if token == nil || len(token.Secret) != 15 || err != nil {
		t.Fatal(token, err)
	}
}
