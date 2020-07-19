package zhsig

import (
	"fmt"
	"os"
	fpath "path/filepath"
	"testing"
)

func TestCreateSig(t *testing.T) {
	err := TempSpace(func(tempdir string) error {
		hostname := "www.example.com"
		groupname := "group"
		docname := "doc.txt"
		datafile := fpath.Join(tempdir, "doc-1.0.0.txt")
		host := NewHost(tempdir, hostname)

		err := CreateSig(host, groupname, docname, datafile)
		if err == nil {
			// データファイルが無かったのに正常
			t.Fatal()
			return fmt.Errorf("err")
		}

		// ファイルを作成
		f, _ := os.Create(datafile)
		f.WriteString("test")
		f.Close()

		err = CreateSig(host, groupname, docname, datafile)
		if err != nil {
			t.Fatal(err)
			return err
		}
		// 自分の公開鍵などをダウンロードフォルダにコピー
		err = copyfile(host.PublicKeyFile(), host.MyPublicKeyFile())
		if err != nil {
			t.Fatal(err)
			return err
		}
		err = copyfile(host.SigFile(docname), host.MyFileSig(docname))
		if err != nil {
			t.Fatal(err)
			return err
		}
		err = copyfile(host.CatalogFile(), host.MyCatalogFile())
		if err != nil {
			t.Fatal(err)
			return err
		}

		// 認証のテスト
		sigfile := host.SigFile(docname)
		if false == fileExists(sigfile) {
			t.Fatal(sigfile)
			return fmt.Errorf("err")
		}
		sig, err := ReadSig(host, docname)
		if err != nil {
			t.Fatal(err)
			return err
		}
		err = sig.VerifyFile(datafile)
		if err != nil {
			t.Fatal(err)
			return err
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
