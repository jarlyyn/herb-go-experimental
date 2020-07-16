package token

import (
	"bytes"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	m := NewMap()
	var _ Manager = m

	token, err := CreateManagedToken(m, "owner", NeverExpired)
	if err != nil {
		panic(err)
	}
	loaded, err := m.Load(token.ID)
	if err != nil {
		panic(err)
	}
	if loaded.ID != token.ID || bytes.Compare(loaded.Secret, token.Secret) != 0 {
		t.Fatal(loaded)
	}
	err = m.Update("notexist", []byte("abcde"))
	if err != ErrTokenNotFound {
		panic(err)
	}
	err = m.Update(token.ID, []byte("abcde"))
	if err != nil {
		panic(err)
	}
	loaded, err = m.Load(token.ID)
	if err != nil {
		panic(err)
	}
	if bytes.Compare(loaded.Secret, []byte("abcde")) != 0 {
		t.Fatal(loaded)
	}
	o, err := IdentifyToken(m, loaded)
	if o != "owner" || err != nil {
		t.Fatal(o, err)
	}
	o, err = Identify(m, loaded.ID, []byte("12345"))
	if o != "" || err != nil {
		t.Fatal(o, err)
	}
	o, err = Identify(m, "notexists", loaded.Secret)
	if o != "" || err != nil {
		t.Fatal(o, err)
	}
	expired := time.Now().Add(-time.Hour)
	err = m.Refresh("notexist", &expired)
	if err != ErrTokenNotFound {
		panic(err)
	}
	err = m.Refresh(token.ID, &expired)
	if err != nil {
		panic(err)
	}
	loaded, err = m.Load(token.ID)
	if err != ErrTokenNotFound {
		panic(err)
	}
	token = New()
	err = m.Regenerate(token)
	if err != nil {
		panic(err)
	}
	token, err = CreateManagedToken(m, "owner", &expired)
	if token == nil || err != nil {
		t.Fatal(loaded, err)
	}
	err = m.Refresh(token.ID, &expired)
	if err != ErrTokenNotFound {
		panic(err)
	}
	err = m.Revoke(token.ID)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err = m.Load(token.ID)
	if loaded != nil || err != ErrTokenNotFound {
		t.Fatal(loaded, err)
	}
}
