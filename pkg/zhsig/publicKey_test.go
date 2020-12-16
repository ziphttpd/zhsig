package zhsig

import (
	"encoding/pem"
	"fmt"
	"testing"
)

func TestPublicKeyFetchFromURL(t *testing.T) {
	err := TempSpace(func(tempdir string) error {
		host := NewHost(tempdir, "www.xorver.com")
		_, err := PublicKeyFetchFromURL(host)
		if err != nil {
			t.Fatal("PublicKey", err)
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPublicKeyFromURL(t *testing.T) {
	err := TempSpace(func(tempdir string) error {
		host := NewHost(tempdir, "www.xorver.com")
		_, err := PublicKey(host)
		if err != nil {
			t.Fatal("PublicKey", err)
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecodePublicKey(t *testing.T) {
	savedi := *di
	// 念のため
	defer func() { *di = savedi }()

	di.ReadFile = func(name string) ([]byte, error) {
		return nil, fmt.Errorf("Read error")
	}
	_, err := DecodePublicKeyFile("not exists file")
	if err == nil {
		t.Fatal()
	}
	*di = savedi

	di.PemDecode = func(data []byte) (p *pem.Block, rest []byte) {
		return nil, nil
	}
	_, err = DecodePublicKeyBytes(nil)
	if err == nil {
		t.Fatal()
	}
	t.Log(err)
	*di = savedi

	di.PemDecode = func(data []byte) (p *pem.Block, rest []byte) {
		return &pem.Block{Type: "PUBLIC KEY"}, nil
	}
	_, err = DecodePublicKeyBytes(nil)
	if err == nil {
		t.Fatal()
	}
	t.Log(err)
	*di = savedi
}
