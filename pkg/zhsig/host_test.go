package zhsig

import (
	"os"
	fpath "path/filepath"
	"testing"
)

func TestHost1(t *testing.T) {
	err := TempSpace(func(tempdir string) error {
		host := NewHost(tempdir, "ziphttpd.com")
		t.Log("Host: " + host.Host())
		t.Log("BasePath: " + host.BasePath())
		// private
		t.Log("PrivateKeyPath: " + host.PrivateKeyPath())
		t.Log("PrivateKeyFile: " + host.PrivateKeyFile())
		// public
		t.Log("MyStorePath: " + host.MyStorePath())
		t.Log("MyPublicKeyFile: " + host.MyPublicKeyFile())
		t.Log("MyPublicKeySig: " + host.MyPublicKeySig())
		// URL
		t.Log("BaseURL: " + host.BaseURL())
		t.Log("PublicKeyURL: " + host.PublicKeyURL())
		t.Log("PublicKeySigURL: " + host.PublicKeySigURL())
		t.Log("CatalogURL: " + host.CatalogURL())
		// store
		t.Log("StorePath: " + host.StorePath())
		t.Log("PublicKeyFile: " + host.PublicKeyFile())
		t.Log("SigFile: " + host.SigFile("test.txt"))

		for _, name := range host.Names() {
			t.Log("name: " + name)
		}

		testTxt := host.File("test.txt", "test_2020.txt")
		t.Log("File: " + testTxt)
		os.MkdirAll(fpath.Dir(testTxt), 0755)
		txt, err := os.Create(testTxt)
		if err != nil {
			return err
		}
		txt.Close()

		testSig := host.SigFile("test.txt")
		t.Log(testSig)
		os.MkdirAll(fpath.Dir(testSig), 0755)
		sig, err := os.Create(testSig)
		if err != nil {
			return err
		}
		sig.Close()

		for _, name := range host.Names() {
			t.Log("name: " + name)
		}
		return nil
	})
	if err != nil {
		t.Log(err)
	}
}
