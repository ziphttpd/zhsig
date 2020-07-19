package main

import (
	"flag"
	"os"
	fpath "path/filepath"

	"github.com/ziphttpd/zhsig/pkg/zhsig"
)

func main() {
	var (
		hostname string
		dir      string
	)
	flag.StringVar(&dir, "dir", "", "configuration directory")
	flag.StringVar(&hostname, "host", "", "hostname (ex. ziphttpd.com)")
	flag.Parse()
	if dir == "" {
		// 無ければ実行ファイルのディレクトリ
		exe, _ := os.Executable()
		dir = fpath.Dir(exe)
	}

	if hostname == "" {
		// 全ホストのカタログのダウンロード
		for _, host := range zhsig.ScanHosts(dir) {
			zhsig.Download(host)
		}
	} else {
		// 特定ホストからのダウンロード
		host := zhsig.NewHost(dir, hostname)
		zhsig.Download(host)
	}
}
