package zhsig

// Download はホストのファイルをダウンロードします。
func Download(host Host) []error {
	errs := []error{}
	if err := DownloadCatalog(host); err == nil {
		cat, err := ReadCatalog(host.CatalogFile())
		if err != nil {
			// TODO: エラー処理
			errs = append(errs, err)
			return errs
		}
		// cats[host.Host()] = cat
		for _, group := range cat.Groups {
			for docid := range group.Docs {
				need := false
				if sig, err := ReadSig(host, docid); err != nil {
					// 未ダウンロードか何かで署名が読み出せなかった
					need = true
				} else {
					sigstr := FetchSig(host, docid)
					if sigstr != "" && sig.Sig() != sigstr {
						// 署名が更新されていた
						need = true
					}
				}
				if need {
					// ダウンロード実行
					if sig, err := DownloadSig(host, docid); err == nil {
						if err := sig.DownloadFile(); err != nil {
							// TODO: ダウンロード失敗
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
