package zhsig

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
)

// PublicKey は共通鍵を取得します。
func PublicKey(host Host) (pub *rsa.PublicKey, err error) {
	if di.FileExists(host.PublicKeyFile()) {
		// ファイルから
		return DecodePublicKeyFile(host.PublicKeyFile())
	}
	return publicKeyFromURL(host)
}

// publicKeyFromURL は署名データ中のURLから最新の公開鍵と署名ファイルを読みだして公開鍵を検証して準備します。
func publicKeyFromURL(host Host) (*rsa.PublicKey, error) {
	// 公開鍵ファイル
	pemBytes, err := GetBytesHTTP(host.PublicKeyURL())
	if err != nil {
		return nil, err
	}

	// 公開鍵の署名ファイル
	keySig, err := DownloadSig(host, PublicPemName)
	if err != nil {
		return nil, err
	}

	// 公開鍵の実際のハッシュ
	pemHash := CalcHashBytes(pemBytes)

	// 公開鍵
	publicKey, err := DecodePublicKeyBytes(pemBytes)
	if err != nil {
		return nil, err
	}

	// 公開鍵の署名データ
	sigBytes, err := base64.StdEncoding.DecodeString(keySig.Sig())
	if err != nil {
		return nil, err
	}

	// 検証
	err = di.VerifyPSS(publicKey, crypto.SHA512, pemHash, sigBytes, nil)
	if err != nil {
		return nil, err
	}

	// 保存
	pubFile, err := di.Create(host.PublicKeyFile())
	if err != nil {
		return nil, err
	}
	_, err = pubFile.Write(pemBytes)
	if err != nil {
		pubFile.Close()
		di.Remove(host.PublicKeyFile())
		return nil, err
	}
	pubFile.Close()

	return publicKey, err
}

// DecodePublicKeyFile は pem ファイルから公開鍵を取り出します。
func DecodePublicKeyFile(pemfile string) (*rsa.PublicKey, error) {
	//pemBytes, err := ioutil.ReadFile(pemfile)
	pemBytes, err := di.ReadFile(pemfile)
	if err != nil {
		return nil, err
	}
	return DecodePublicKeyBytes(pemBytes)
}

// DecodePublicKeyBytes は pem ファイルのバイト列から公開鍵を取り出します。
func DecodePublicKeyBytes(pemByte []byte) (*rsa.PublicKey, error) {
	//pb, _ := pem.Decode(pemByte)
	pb, _ := di.PemDecode(pemByte)
	if pb == nil {
		return nil, fmt.Errorf("invalid public key pem")
	}
	if pb.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("invalid public key type : %s", pb.Type)
	}
	//key, err := x509.ParsePKCS1PublicKey(pb.Bytes)
	key, err := di.ParsePKCS1PublicKey(pb.Bytes)
	return key, err
}
