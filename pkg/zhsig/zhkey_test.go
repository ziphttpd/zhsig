package zhsig

import (
	"testing"
)

func TestZhkey(t *testing.T) {
	err := TempSpace(func(tempdir string) error {
		host := NewHost(tempdir, "ziphttpd.com")
		err := CreateZhkey(host)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Log(err)
	}
}
