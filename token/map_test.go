package token

import (
	"bytes"
	"testing"
)

func TestMap(t *testing.T) {
	m := NewMap()
	token := New()
	token.ID = "test"
	err := Regenerate(BytesGenerator(15), token)
	if token == nil || len(token.Secret) != 15 || err != nil {
		t.Fatal(token, err)
	}
	token, err = m.Create("owner", token.Secret, NeverExpired)
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
	err = m.Revoke("test")
	if err != nil {
		panic(err)
	}
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
