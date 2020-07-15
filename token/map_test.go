package token

import (
	"bytes"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	m := NewMap()
	token, err := m.Create("owner", []byte("12345"), NeverExpired)
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
	if err != ErrIDNotFound {
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
	if err != ErrIDNotFound {
		panic(err)
	}
	err = m.Refresh(token.ID, &expired)
	if err != nil {
		panic(err)
	}
	loaded, err = m.Load(token.ID)
	if err != ErrIDNotFound {
		panic(err)
	}
	token, err = GeneratAndCreate(m, BytesGenerator(15), "owner", &expired)
	if token == nil || err != nil {
		t.Fatal(loaded, err)
	}
	err = m.Refresh(token.ID, &expired)
	if err != ErrIDNotFound {
		panic(err)
	}
	err = m.Revoke(token.ID)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err = m.Load(token.ID)
	if loaded != nil || err != ErrIDNotFound {
		t.Fatal(loaded, err)
	}
}

func TestEncoding(t *testing.T) {
	token := New()
	token.ID = "test"
	token.Secret = []byte{1, 2, 3, 4, 5}
	encoded, err := Base64Encoding.Encode(token.Secret)
	if err != nil {
		panic(err)
	}
	if encoded == string(token.Secret) {
		t.Fatal(encoded)
	}
	decoded, err := Base64Encoding.Decode(encoded)
	if err != nil {
		panic(err)
	}
	if bytes.Compare(decoded, token.Secret) != 0 {
		t.Fatal(decoded)
	}
	encoded, err = StringEncoding.Encode(token.Secret)
	if err != nil {
		panic(err)
	}
	if encoded != string(token.Secret) {
		t.Fatal(encoded)
	}
	decoded, err = StringEncoding.Decode(string(token.Secret))
	if err != nil {
		panic(err)
	}
	if bytes.Compare(decoded, token.Secret) != 0 {
		t.Fatal(decoded)
	}
}
