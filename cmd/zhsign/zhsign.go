package main

import (
	"flag"
	"fmt"
	"os"
	fpath "path/filepath"

	"github.com/ziphttpd/zhsig/pkg/zhsig"
)

func main() {
	var (
		hostname string
		err      error
		dir      string
		group    string
		name     string
		file     string
	)
	flag.StringVar(&hostname, "host", "", "hostname (ex. ziphttpd.com)")
	flag.StringVar(&dir, "dir", "", "configuration directory")
	flag.StringVar(&group, "group", "", "document group name")
	flag.StringVar(&name, "name", "", "document name")
	flag.StringVar(&file, "file", "", "document file")
	flag.Parse()
	if hostname == "" {
		flag.PrintDefaults()
		return
	}
	if dir == "" {
		// 無ければ実行ファイルのディレクトリ
		exe, _ := os.Executable()
		dir = fpath.Dir(exe)
	}
	// ホストの関連情報生成
	host := zhsig.NewHost(dir, hostname)
	// site.json の生成
	// https://github.com/ziphttpd/zhsig/issues/1 の暫定対応
	err = zhsig.CreateSiteFile(host)
	if err != nil {
		fmt.Println(err)
	}

	// 指定されたファイルを署名
	if name != "" || file != "" {
		if name != "" && file != "" {
			if f, err := os.Stat(file); os.IsNotExist(err) || f.IsDir() {
				// TODO: ファイルなし
			} else {
				err = zhsig.CreateSig(host, group, name, file)
			}
		} else {
			flag.PrintDefaults()
		}
	}
	if err != nil {
		fmt.Println(err)
	}
}
