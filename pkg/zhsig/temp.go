package zhsig

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// TempSpace では一時領域を作成し、callback() を実行した後に削除します。
func TempSpace(callback func(tempdir string) error) (err error) {
	r := make([]byte, 8)
	rand.Read(r)
	tempdir := filepath.Join(os.TempDir(), "tempspace_"+hex.EncodeToString(r))

	// 一時領域を作成
	err = os.MkdirAll(tempdir, 0600)
	if err != nil {
		return err
	}
	// 一時領域を削除
	defer func() {
		os.RemoveAll(tempdir)
		os.Remove(tempdir)
	}()
	// パニックをキャッチ
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("Run() return %+v", r)
		}
	}()

	// TempSpace を実行
	err = callback(tempdir)

	return err
}
