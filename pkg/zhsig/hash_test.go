package zhsig

import (
	"bytes"
	"fmt"
	"os"
	fpath "path/filepath"
	"strings"
	"testing"
)

type tstreader struct{}

func (t *tstreader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("Read error")
}

func TestCalcHash(t *testing.T) {
	var err error

	h1 := CalcHashBytes([]byte("test"))
	h2, err := CalcHashStream(strings.NewReader("test"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(h1, h2) != 0 {
		t.Fatal()
	}

	// Read error
	_, err = CalcHashStream(&tstreader{})
	if err == nil {
		t.Fatal(err)
	}
	t.Log(err)
}
func TestCalcHashFile(t *testing.T) {
	err := TempSpace(func(tempdir string) error {
		testfile := fpath.Join(tempdir, "test.txt")
		_, err := CalcHashFile(testfile)
		if err == nil {
			return fmt.Errorf("exists error")
		}
		f, err := os.Create(testfile)
		if err != nil {
			return err
		}
		f.WriteString("test")
		f.Close()
		_, err = CalcHashFile(testfile)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
