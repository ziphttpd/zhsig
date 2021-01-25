package zhsig

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"os"
)

// CreateSiteFile はサイト名署名ファイルを作成します。
func CreateSiteFile(host Host) error {
	var err error
	sitefile := host.MySiteFile()
	if _, err := os.Stat(sitefile); os.IsNotExist(err) {
		if pkey, err := PrivateKey(host); err == nil {
			nameHash := CalcHashBytes([]byte(host.Host()))
			if base64, err := SignatureBASE64(pkey, nameHash); err == nil {
				if file, err := os.Create(sitefile); err == nil {
					fmt.Fprintf(file, "{\"site\":\"%s\"}", base64)
					file.Close()
				}
			}
		}
	}
	return err
}

// VerifySiteFileBytes はサイト名署名を検証する
func VerifySiteFileBytes(publicKey *rsa.PublicKey, host Host, bytes []byte) bool {
	// 検証
	nameHash := CalcHashBytes([]byte(host.Host()))
	err := di.VerifyPSS(publicKey, crypto.SHA512, nameHash, bytes, nil)
	return err == nil
}
