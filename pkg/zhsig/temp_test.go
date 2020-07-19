package zhsig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	var tdir string
	err := TempSpace(func(tempdir string) error {
		tdir = tempdir
		data := "ABCD"
		filename := filepath.Join(tempdir, "test.txt")
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(data)
		if err != nil {
			return err
		}
		f.Close()
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		if string(b) != data {
			return fmt.Errorf("err")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if fileExists(tdir) {
		t.Fatal(tdir)
	}
}
func TestRunPanic(t *testing.T) {
	err := TempSpace(func(tempdir string) error {
		panic("err")
	})
	if err == nil {
		t.Fatal()
	}
}
