package zhsig

import (
	"os"
	fpath "path/filepath"

	"github.com/xorvercom/util/pkg/fileutil"
)

// Update はホストのファイルをダウンロードします。
func Update(host Host) []error {
	errs := []error{}
	store := host.StorePath()
	if err := DownloadCatalog(host); err == nil {
		cat, err := ReadCatalog(host.CatalogFile())
		if err != nil {
			// TODO: エラー処理
			errs = append(errs, err)
			return errs
		}
		// cats[host.Host()] = cat
		for groupname, group := range cat.Groups {
			// ドキュメントグループが未ダウンロードの場合にはスキップ
			if fi, err := os.Stat(fpath.Join(store, groupname)); err != nil {
				continue
			} else {
				if false == fi.IsDir() {
					continue
				}
			}
			for docid := range group.Docs {
				need := false
				if sig, err := ReadSig(host, docid); err != nil {
					// 未ダウンロードか何かで署名が読み出せなかった
					need = true
				} else {
					// 本体はダウンロードされている
					target := host.File(docid, sig.File())
					if _, err := os.Stat(target); err == nil {
						// 署名変更チェック
						sigstr := FetchSig(host, docid)
						if sigstr != "" && sig.Sig() != sigstr {
							// 署名が更新されていた
							need = true
						}
					} else {
						// 本体がない
						need = true
					}
				}
				if need {
					// ダウンロード実行
					if sig, err := DownloadSig(host, docid); err == nil {
						if err := sig.DownloadFile(); err != nil {
							// TODO: ダウンロード失敗
							fileutil.FileIfDelete(host.SigFile(docid))
							errs = append(errs, err)
						}
					} else {
						errs = append(errs, err)
					}
				}
			}
		}
	}
	return errs
}

// Download はホストのファイルをダウンロードします。
func Download(host Host, group string) []error {
	errs := []error{}
	if err := DownloadCatalog(host); err == nil {
		cat, err := ReadCatalog(host.CatalogFile())
		if err != nil {
			// TODO: エラー処理
			errs = append(errs, err)
			return errs
		}
		if group, ok := cat.Groups[group]; ok {
			for docid := range group.Docs {
				need := false
				if sig, err := ReadSig(host, docid); err != nil {
					// 未ダウンロードか何かで署名が読み出せなかった
					need = true
				} else {
					// 本体はダウンロードされている
					target := host.File(docid, sig.File())
					if _, err := os.Stat(target); err == nil {
						// 署名変更チェック
						sigstr := FetchSig(host, docid)
						if sigstr != "" && sig.Sig() != sigstr {
							// 署名が更新されていた
							need = true
						}
					} else {
						// 本体がない
						need = true
					}
				}
				if need {
					// ダウンロード実行
					if sig, err := DownloadSig(host, docid); err == nil {
						if err := sig.DownloadFile(); err != nil {
							// TODO: ダウンロード失敗
							fileutil.FileIfDelete(host.SigFile(docid))
							errs = append(errs, err)
						}
					} else {
						errs = append(errs, err)
					}
				}
			}
		}
	}
	return errs
}
