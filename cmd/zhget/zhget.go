package main

import (
	"flag"
	"os"
	fpath "path/filepath"

	"github.com/ziphttpd/zhsig/pkg/zhsig"
)

// const (
// 	listURL      = "https://ziphttpd.com/api/v1/list"
// 	sitelistname = "sitelist.json"
// )

func main() {
	var (
		hostname  string
		dir       string
		groupname string
	)
	flag.StringVar(&dir, "dir", "", "configuration directory")
	flag.StringVar(&hostname, "host", "", "hostname (ex. ziphttpd.com)")
	flag.StringVar(&groupname, "group", "", "document group (ex. ziphttpd)")
	flag.Parse()
	if dir == "" {
		// 無ければ実行ファイルのディレクトリ
		exe, _ := os.Executable()
		dir = fpath.Dir(exe)
	}

	// ziphttpd.com からサイト一覧を得る
	// sitelistfile := fpath.Join(dir, zhsig.StorePath, sitelistname)
	// if err := zhsig.GetFileHTTP(listURL, sitelistfile); err == nil {
	// 	// noop
	// }
	// sitelist := []string{}
	// if elem, err := json.LoadFromJSONFile(sitelistfile); err == nil {
	// 	if eo, ok := elem.AsObject(); ok {
	// 		if el, ok := eo.Child("Hosts").AsArray(); ok {
	// 			if es, ok := el.AsString(); ok {
	// 				sitelist = append(sitelist, es.Text())
	// 			}
	// 		}
	// 	}
	// }
	// sort.Strings(sitelist)

	if hostname == "" {
		// 以前にダウンロードした全ホストからのアップデート
		for _, host := range zhsig.ScanHosts(dir) {
			zhsig.Update(host)
		}
	} else {
		// 指定ホストからのダウンロード
		host := zhsig.NewHost(dir, hostname)
		if groupname == "" {
			// 既にダウンロードしているグループのアップデート
			zhsig.Update(host)
		} else {
			// 指定グループをダウンロード (一度ダウンロードしなければアップデートされない)
			zhsig.Download(host, groupname)
		}
	}
}
