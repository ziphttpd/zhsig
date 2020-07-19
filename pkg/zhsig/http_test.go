package zhsig

import (
	"os"
	fpath "path/filepath"
	"testing"
)

func TestGetBytesHTTP(t *testing.T) {
	err := TempSpace(func(tempdir string) error {
		host := NewHost(tempdir, "www.xorver.com")
		_, err := GetBytesHTTP("xxx" + host.PublicKeyURL())
		if err == nil {
			t.Error()
		}
		_, err = GetBytesHTTP(host.PublicKeyURL() + ".notexist")
		if err == nil {
			t.Error()
		}
		_, err = GetBytesHTTP(host.PublicKeyURL())
		return err
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGetFileHTTP(t *testing.T) {
	err := TempSpace(func(tempdir string) error {
		host := NewHost(tempdir, "www.xorver.com")
		path := fpath.Join(tempdir, PublicPemName)
		err := GetFileHTTP("xxx"+host.PublicKeyURL(), path)
		if err == nil {
			t.Error()
		}
		err = GetFileHTTP(host.PublicKeyURL()+".notexist", path)
		if err == nil {
			t.Error()
		}
		err = GetFileHTTP(host.PublicKeyURL(), "/"+path)
		if err == nil {
			t.Error()
		}
		err = GetFileHTTP(host.PublicKeyURL(), path)
		if err != nil {
			return err
		}
		_, err = os.Stat(path)
		return err
	})
	if err != nil {
		t.Error(err)
	}
}
