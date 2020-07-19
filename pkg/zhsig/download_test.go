package zhsig

import (
	"fmt"
	"testing"
)

func TestDownload(t *testing.T) {
	err := TempSpace(func(tempdir string) error {
		host := NewHost(tempdir, "www.xorver.com")

		errs := Download(host)
		for _, err := range errs {
			t.Error(err)
		}
		if len(errs) != 0 {
			return fmt.Errorf("error")
		}

		if dirExists(host.MyStorePath()) {
			t.Fatal(host.MyStorePath())
		}
		if false == dirExists(host.StorePath()) {
			t.Fatal(host.StorePath())
		}

		cat, err := ReadCatalog(host.CatalogFile())
		if err != nil {
			return err
		}
		if len(cat.Groups) == 0 {
			return fmt.Errorf("error")
		}
		for _, gname := range cat.GroupNames() {
			g := cat.Groups[gname]
			if len(g.Docs) == 0 {
				err = fmt.Errorf("no Docs")
				t.Fatal(err)
				return err
			}
			for _, docname := range g.DocNames() {
				d := g.Docs[docname]
				t.Log(d.Title)
				sig, err := ReadSig(host, docname)
				if err != nil {
					t.Fatal(err)
					return err
				}
				file := host.File(docname, sig.File())
				if fileExists(file) == false {
					err = fmt.Errorf("%s not found", file)
					t.Fatal(err)
					return err
				}
			}
		}

		// 再実行
		errs = Download(host)
		for _, err := range errs {
			t.Error(err)
		}
		if len(errs) != 0 {
			return fmt.Errorf("error")
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
