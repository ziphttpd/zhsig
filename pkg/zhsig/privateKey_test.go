package zhsig

import (
	"encoding/pem"
	"fmt"
	"testing"
)

func TestDecodePrivateKey(t *testing.T) {
	savedi := *di
	// 念のため
	defer func() { *di = savedi }()

	di.ReadFile = func(name string) ([]byte, error) {
		return nil, fmt.Errorf("Read error")
	}
	_, err := DecodePrivateKeyFile("not exists file")
	if err == nil {
		t.Fatal()
	}
	*di = savedi

	di.PemDecode = func(data []byte) (p *pem.Block, rest []byte) {
		return nil, nil
	}
	_, err = DecodePrivateKeyBytes(nil)
	if err == nil {
		t.Fatal()
	}
	t.Log(err)
	*di = savedi

	di.PemDecode = func(data []byte) (p *pem.Block, rest []byte) {
		return &pem.Block{Type: "PRIVATE KEY"}, nil
	}
	_, err = DecodePrivateKeyBytes(nil)
	if err == nil {
		t.Fatal()
	}
	t.Log(err)
	*di = savedi
}
