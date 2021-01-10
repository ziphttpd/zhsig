package zhsig

import (
	"fmt"
	"io"
	fpath "path/filepath"
)

// GetBytesHTTP は URL からバイト列を読み出します。
func GetBytesHTTP(url string) ([]byte, error) {
	res, err := di.HTTPGet(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("%s return %s", url, res.Status)
	}
	return di.ReadAll(res.Body)
}

// GetFileHTTP は URL からファイルにダウンロードします。
func GetFileHTTP(url string, file string) error {
	res, err := di.HTTPGet(url)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("%s return %s", url, res.Status)
	}
	di.MkdirAll(fpath.Dir(file), 0755)
	f, err := di.Create(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, res.Body)
	return err
}
