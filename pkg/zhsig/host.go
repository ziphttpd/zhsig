package zhsig

import (
	"crypto/rsa"
	fpath "path/filepath"
)

const (
	sigExt = ".zhsig"
	// PrivatePemName は秘密鍵の出力ファイル名です。
	PrivatePemName = "private.pem"
	// PublicPemName は公開鍵の出力ファイル名です。
	PublicPemName = "public.pem"
	// PublicSigName は公開鍵の署名ファイル名です。
	PublicSigName = "public.pem" + sigExt
	// CatalogName はカタログファイル名です。
	CatalogName = "catalog.json"
	privatePath = "private"
	publicPath  = "public"
	storePath   = "store"
)

// Host はホストを表します。
type Host interface {
	// Host はホストを返します。
	Host() string
	// BasePath は基準パスを返します。
	BasePath() string

	// BaseURL は公開鍵などのURLを返します。
	BaseURL() string
	// FileURL は BaseURL のファイルのURLを返します。
	FileURL(file string) string
	// PublicKeyURL は公開鍵のURLを返します。
	PublicKeyURL() string
	// PublicKeySigURL は公開鍵自己署名のURLを返します。
	PublicKeySigURL() string
	// CatalogURL はカタログファイルのURLを返します。
	CatalogURL() string

	// PrivateKeyPath は秘密鍵のパスを返します。
	PrivateKeyPath() string
	// PrivateKeyFile は秘密鍵ファイルのパスを返します。
	PrivateKeyFile() string

	// MyStorePath は自分のファイルの置き場所を返します。
	MyStorePath() string
	// MyFile は自分のファイルのパスを返します。
	MyFile(filename string) string
	// MyPublicKeyFile は自分の公開鍵ファイルのパスを返します。
	MyPublicKeyFile() string
	// MyPublicKeySig は自分の公開鍵署名ファイルのパスを返します。
	MyPublicKeySig() string
	// MyFileSig は自分の署名ファイルのパスを返します。
	MyFileSig(name string) string
	// MyCatalogFile は自分のカタログのパスを返します。
	MyCatalogFile() string

	// StorePath は保存パスを返します。
	StorePath() string
	// File はファイルのパスを返します。
	File(name string, file string) string
	// SigFile はファイルの署名パスを返します。
	SigFile(name string) string
	// SigURL は署名ファイルのURLを返します。
	SigURL(name string) string
	// CatalogFile はカタログファイルのパスを返します。
	CatalogFile() string
	// PublicKeyFile は公開鍵ファイルのパスを返します。
	PublicKeyFile() string
	// Names は name の一覧を返します。
	Names() []string

	// PrivateKey は秘密鍵をキャッシュして返します。
	PrivateKey() (*rsa.PrivateKey, error)
	// PublicKey は公開鍵をキャッシュして返します。
	PublicKey() (*rsa.PublicKey, error)
}

// HostInst はサイトを表します。
type HostInst struct {
	base string
	host string
	// 公開鍵のキャッシュ
	privateKey *rsa.PrivateKey
	// 公開鍵のキャッシュ
	publicKey *rsa.PublicKey
}

// ScanHosts はダウンロードしたHostを列挙します。
func ScanHosts(baseDir string) []*HostInst {
	dirs := []*HostInst{}
	if fis, err := di.ReadDir(fpath.Join(baseDir, storePath)); err == nil {
		for _, fi := range fis {
			if fi.IsDir() {
				dirs = append(dirs, NewHost(baseDir, fi.Name()))
			}
		}
	}
	return dirs
}

// NewHost は Host を生成します。
func NewHost(baseDir, host string) *HostInst {
	s := &HostInst{base: baseDir, host: host}
	return s
}

// Host はホスト名を返します。
func (s *HostInst) Host() string {
	return s.host
}

// BasePath は公開鍵の基準パスを返します。
func (s *HostInst) BasePath() string {
	return s.base
}

// PrivateKeyPath は秘密鍵のパスを返します。
func (s *HostInst) PrivateKeyPath() string {
	return fpath.Join(s.base, privatePath, s.host)
}

// PrivateKeyFile は秘密鍵ファイルのパスを返します。
func (s *HostInst) PrivateKeyFile() string {
	return fpath.Join(s.PrivateKeyPath(), PrivatePemName)
}

// MyStorePath は自分のファイルの置き場所を返します。
func (s *HostInst) MyStorePath() string {
	return fpath.Join(s.base, publicPath, s.host)
}

// MyFile は自分のファイルのパスを返します。
func (s *HostInst) MyFile(filename string) string {
	_, f := fpath.Split(filename)
	return fpath.Join(s.MyStorePath(), f)
}

// MyPublicKeyFile は自分の公開鍵ファイルのパスを返します。
func (s *HostInst) MyPublicKeyFile() string {
	return s.MyFile(PublicPemName)
}

// MyPublicKeySig は自分の公開鍵署名ファイルのパスを返します。
func (s *HostInst) MyPublicKeySig() string {
	return s.MyFile(PublicSigName)
}

// MyFileSig は自分の署名ファイルのパスを返します。
func (s *HostInst) MyFileSig(name string) string {
	return s.MyFile(name + sigExt)
}

// MyCatalogFile は自分のカタログのパスを返します。
func (s *HostInst) MyCatalogFile() string {
	return s.MyFile(CatalogName)
}

// BaseURL は公開鍵の基準パスを返します。
func (s *HostInst) BaseURL() string {
	return "https://" + s.host + "/sig/"
}

// FileURL は署名のURLを返します。
func (s *HostInst) FileURL(filename string) string {
	_, f := fpath.Split(filename)
	return s.BaseURL() + f
}

// PublicKeyURL は公開鍵のURLを返します。
func (s *HostInst) PublicKeyURL() string {
	return s.FileURL(PublicPemName)
}

// PublicKeySigURL は公開鍵自己署名のURLを返します。
func (s *HostInst) PublicKeySigURL() string {
	return s.FileURL(PublicSigName)
}

// CatalogURL は公カタログのURLを返します。
func (s *HostInst) CatalogURL() string {
	return s.FileURL(CatalogName)
}

// StorePath は保存ディレクトリを返します。
func (s *HostInst) StorePath() string {
	return fpath.Join(s.base, storePath, s.host)
}

// File はファイルのパスを返します。
// ファイルはSigInfo.File()を使ってください。
func (s *HostInst) File(name string, filename string) string {
	_, f := fpath.Split(filename)
	return fpath.Join(s.StorePath(), name, f)
}

// SigFile は署名ファイルのパスを返します。
func (s *HostInst) SigFile(name string) string {
	return fpath.Join(s.StorePath(), name+sigExt)
}

// SigURL は署名のURLを返します。
func (s *HostInst) SigURL(name string) string {
	return s.FileURL(name + sigExt)
}

// CatalogFile はカタログファイルのパスを返します。
func (s *HostInst) CatalogFile() string {
	return fpath.Join(s.StorePath(), CatalogName)
}

// PublicKeyFile は公開鍵ファイルのパスを返します。
func (s *HostInst) PublicKeyFile() string {
	return fpath.Join(s.StorePath(), PublicPemName)
}

// Names は name の一覧を返します。
func (s *HostInst) Names() []string {
	ret := []string{}
	fis, err := di.ReadDir(s.StorePath())
	if err != nil {
		return ret
	}
	for _, fi := range fis {
		if fi.IsDir() {
			name := fi.Name()
			if fileExists(name + sigExt) {
				ret = append(ret, name)
			}
		}
	}
	return ret
}

// PrivateKey は秘密鍵をキャッシュして返します。
func (s *HostInst) PrivateKey() (*rsa.PrivateKey, error) {
	if s.privateKey != nil {
		return s.privateKey, nil
	}
	k, err := PrivateKey(s)
	if err == nil {
		s.privateKey = k
	}
	return k, err
}

// PublicKey は公開鍵をキャッシュして返します。
func (s *HostInst) PublicKey() (*rsa.PublicKey, error) {
	if s.publicKey != nil {
		return s.publicKey, nil
	}
	k, err := PublicKey(s)
	if err == nil {
		s.publicKey = k
	}
	return k, err
}
