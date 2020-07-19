package zhsig

import (
	// crypto
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	fpath "path/filepath"

	// encoding
	//"encoding/json"
	"encoding/pem"

	// io
	"io"
	"io/ioutil"

	// その他
	"bytes"
	"net/http"
	"os"
)

// テスト時に coverage を向上させるためライブラリ呼び出しを置き換え可能にする実験的対策
type distruct struct {
	Compare                func([]byte, []byte) int
	ReadFile               func(filename string) ([]byte, error)
	ReadAll                func(io.Reader) ([]byte, error)
	ReadDir                func(dirname string) ([]os.FileInfo, error)
	ParsePKCS1PrivateKey   func(der []byte) (*rsa.PrivateKey, error)
	ParsePKCS1PublicKey    func(der []byte) (*rsa.PublicKey, error)
	MarshalPKCS1PrivateKey func(key *rsa.PrivateKey) []byte
	MarshalPKCS1PublicKey  func(key *rsa.PublicKey) []byte
	PemDecode              func(data []byte) (p *pem.Block, rest []byte)
	PemEncode              func(out io.Writer, b *pem.Block) error
	Open                   func(name string) (*os.File, error)
	Create                 func(name string) (*os.File, error)
	Remove                 func(name string) error
	MkdirAll               func(path string, perm os.FileMode) error
	GenerateKey            func(random io.Reader, bits int) (*rsa.PrivateKey, error)
	SignPSS                func(rand io.Reader, priv *rsa.PrivateKey, hash crypto.Hash, hashed []byte, opts *rsa.PSSOptions) ([]byte, error)
	VerifyPSS              func(pub *rsa.PublicKey, hash crypto.Hash, hashed []byte, sig []byte, opts *rsa.PSSOptions) error
	HTTPGet                func(url string) (resp *http.Response, err error)
	FileExists             func(filename string) bool
}

var di *distruct

func init() {
	di = &distruct{
		Compare:                bytes.Compare,
		ReadFile:               ioutil.ReadFile,
		ReadAll:                ioutil.ReadAll,
		ReadDir:                ioutil.ReadDir,
		ParsePKCS1PrivateKey:   x509.ParsePKCS1PrivateKey,
		ParsePKCS1PublicKey:    x509.ParsePKCS1PublicKey,
		MarshalPKCS1PrivateKey: x509.MarshalPKCS1PrivateKey,
		MarshalPKCS1PublicKey:  x509.MarshalPKCS1PublicKey,
		PemDecode:              pem.Decode,
		PemEncode:              pem.Encode,
		Open:                   os.Open,
		Create:                 os.Create,
		Remove:                 os.Remove,
		MkdirAll:               os.MkdirAll,
		GenerateKey:            rsa.GenerateKey,
		SignPSS:                rsa.SignPSS,
		VerifyPSS:              rsa.VerifyPSS,
		HTTPGet:                http.Get,
		FileExists:             fileExists,
	}
}

func fileExists(filename string) bool {
	f, err := os.Stat(filename)
	return err == nil && false == f.IsDir()
}

func dirExists(filename string) bool {
	f, err := os.Stat(filename)
	return err == nil && f.IsDir()
}

func copyfile(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	dir, _ := fpath.Split(dst)
	os.MkdirAll(dir, os.ModeDir)
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()
	_, err = io.Copy(d, s)
	return err
}

func keys(m map[string]interface{}) []string {
	ret := []string{}
	for key := range m {
		ret = append(ret, key)
	}
	return ret
}
