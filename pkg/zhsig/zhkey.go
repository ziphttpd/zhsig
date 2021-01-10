package zhsig

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
)

// CreateZhkey はヘルプ証明書を作成します。
func CreateZhkey(host Host) (reterr error) {
	di.MkdirAll(host.PrivateKeyPath(), 0755)
	di.MkdirAll(host.MyStorePath(), 0755)
	// ファイル
	pathPriPem := host.PrivateKeyFile()
	pathPubPem := host.MyPublicKeyFile()
	pathPubSig := host.MyPublicKeySig()
	urlFile := host.PublicKeyURL()
	di.Remove(pathPriPem)
	di.Remove(pathPubPem)
	di.Remove(pathPubSig)

	defer func() {
		err := recover()
		if err != nil {
			di.Remove(pathPriPem)
			di.Remove(pathPubPem)
			di.Remove(pathPubSig)
			reterr = err.(error)
		} else {
			reterr = nil
		}
	}()

	// 鍵ペア
	//priKey, err := rsa.GenerateKey(rand.Reader, 2048)
	priKey, err := di.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	pubKey := priKey.Public()

	// PKCS1
	//priBytes := x509.MarshalPKCS1PrivateKey(priKey)
	priBytes := di.MarshalPKCS1PrivateKey(priKey)
	var pubBytes []byte
	if rsaPublicKeyPointer, ok := pubKey.(*rsa.PublicKey); ok {
		//pubBytes = x509.MarshalPKCS1PublicKey(rsaPublicKeyPointer)
		pubBytes = di.MarshalPKCS1PublicKey(rsaPublicKeyPointer)
	}

	// private.pem の作成
	//priFile, err := os.Create(tempPriPem)
	priFile, err := di.Create(pathPriPem)
	if err != nil {
		panic(err)
	}
	// pem にエンコード
	//err = pem.Encode(priFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: priBytes})
	err = di.PemEncode(priFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: priBytes})
	if err != nil {
		panic(err)
	}
	priFile.Close()

	// public.pem の作成
	//pubFile, err := os.Create(tempPubPem)
	pubFile, err := di.Create(pathPubPem)
	if err != nil {
		panic(err)
	}
	//err = pem.Encode(pubFile, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes})
	err = di.PemEncode(pubFile, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes})
	if err != nil {
		panic(err)
	}
	pubFile.Close()

	// 共通鍵の署名を作成
	s, err := CreateSignatureBASE64(priKey, pathPubPem)
	if err != nil {
		panic(err)
	}
	siginfo := &sigInfo{sig: &sig{
		name: PublicPemName,
		host: host.Host(),
		url:  urlFile,
		sig:  s,
	}, host: host}

	return siginfo.Save(pathPubSig)
}
