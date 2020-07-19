package zhsig

import (
	"crypto/rsa"
	"fmt"
)

// PrivateKey は秘密鍵を取得します。
func PrivateKey(host Host) (*rsa.PrivateKey, error) {
	f := host.PrivateKeyFile()
	if false == fileExists(f) {
		err := CreateZhkey(host)
		if err != nil {
			return nil, err
		}
	}
	return DecodePrivateKeyFile(f)
}

// DecodePrivateKeyFile は pem ファイルから秘密鍵を取り出します。
func DecodePrivateKeyFile(pemfile string) (*rsa.PrivateKey, error) {
	//pemBytes, err := ioutil.ReadFile(pemfile)
	pemBytes, err := di.ReadFile(pemfile)
	if err != nil {
		return nil, err
	}
	return DecodePrivateKeyBytes(pemBytes)
}

// DecodePrivateKeyBytes は pem ファイルのバイト列から秘密鍵を取り出します。
func DecodePrivateKeyBytes(pemByte []byte) (*rsa.PrivateKey, error) {
	//pb, _ := pem.Decode(pemByte)
	pb, _ := di.PemDecode(pemByte)
	if pb == nil {
		return nil, fmt.Errorf("invalid private key pem")
	}
	if pb.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key type : %s", pb.Type)
	}
	//key, err := x509.ParsePKCS1PrivateKey(pb.Bytes)
	key, err := di.ParsePKCS1PrivateKey(pb.Bytes)
	return key, err
}
