// Package zhsig は zhsig ファイルを元として署名を検証するユーティリティです。
package zhsig

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/url"
	fpath "path/filepath"

	"github.com/xorvercom/util/pkg/json"
)

// SigInfo は署名情報です。
type SigInfo interface {
	// Name はファイル名に依存しないユニーク名です
	Name() string
	// Host はホスト名です。
	Host() string
	// URL は署名対象ファイルのURLです。
	URL() string
	// File は署名対象ファイルのファイル名です。
	File() string
	// Sig はファイルの署名をBASE64エンコーディングした文字列です。
	Sig() string
	// VerifyFile はファイルを再検証します
	VerifyFile(filename string) error
	// DownloadFile はファイルを検証ダウンロードします。
	DownloadFile() error
	// Save は署名を保存します。
	Save(filename string) error
}

const (
	pathHost = "host"
	pathName = "name"
	pathURL  = "url"
	pathSig  = "sig"
)

type sig struct {
	// host はホスト名です。
	host string
	// name はファイル名に依存しないユニーク名です
	name string
	// url は署名対象ファイルのURLです。
	url string
	// sig はファイルの署名をBASE64エンコーディングした文字列です。
	sig string
}

// sigInfo は zhsig ファイルの内容です。
type sigInfo struct {
	sig  *sig
	host Host
}

// Name はファイル名に依存しないユニーク名です
func (z *sigInfo) Name() string {
	return z.sig.name
}

// Host はホスト名です。
func (z *sigInfo) Host() string {
	return z.sig.host
}

// URL は署名対象ファイルのURLです。
func (z *sigInfo) URL() string {
	return z.sig.url
}

// File は署名対象ファイルのファイル名です。
func (z *sigInfo) File() string {
	url, err := url.Parse(z.sig.url)
	if err != nil {
		return ""
	}
	return fpath.Base(url.Path)
}

// Sig はファイルの署名をBASE64エンコーディングした文字列です。
func (z *sigInfo) Sig() string {
	return z.sig.sig
}

// VerifyFile はファイルを再検証します
func (z *sigInfo) VerifyFile(filename string) error {
	var err error
	publicKey, err := z.host.PublicKey()
	if err != nil {
		return err
	}
	sigBytes, err := base64.StdEncoding.DecodeString(z.Sig())
	if err != nil {
		return err
	}

	// ファイル名は name 由来ではない元ファイル名を使用する
	fileHash, err := CalcHashFile(filename)
	if err != nil {
		return err
	}

	// 検証
	err = di.VerifyPSS(publicKey, crypto.SHA512, fileHash, sigBytes, nil)
	return err
}

// DownloadFile はファイルを検証ダウンロードします。
func (z *sigInfo) DownloadFile() error {
	file := z.host.File(z.Name(), z.File())
	err := GetFileHTTP(z.URL(), file)
	if err != nil {
		return err
	}
	err = z.VerifyFile(file)
	if err != nil {
		di.Remove(file)
	}
	return err
}

// Save は署名を保存します。
func (z *sigInfo) Save(path string) error {
	di.MkdirAll(fpath.Dir(path), 0755)
	elem := json.NewElemObject()
	elem.Put(pathHost, json.NewElemString(z.Host()))
	elem.Put(pathName, json.NewElemString(z.Name()))
	elem.Put(pathURL, json.NewElemString(z.URL()))
	elem.Put(pathSig, json.NewElemString(z.Sig()))
	return json.SaveToJSONFile(path, elem, true)
}

// CreateSig はdatafileの署名を作成します。
func CreateSig(host Host, groupname string, docname string, datafile string) error {
	// 秘密鍵
	pkey, err := host.PrivateKey()
	if err != nil {
		return err
	}

	// datafile の署名を作成
	sigstr, err := CreateSignatureBASE64(pkey, datafile)
	if err != nil {
		return err
	}
	siginfo := &sigInfo{sig: &sig{
		name: docname,
		host: host.Host(),
		url:  host.FileURL(datafile),
		sig:  sigstr,
	}, host: host}
	siginfo.Save(host.MyFileSig(docname))

	// 新規の登録ならばカタログに追加
	catalogfile := host.MyCatalogFile()
	var cat *Catalog
	if false == fileExists(catalogfile) {
		// カタログが作られていない
		cat = NewCatalog()
		err = json.SaveToJSONFile(catalogfile, cat.JSON(), true)
	}
	cat, err = ReadCatalog(catalogfile)
	if err != nil {
		return err
	}

	group, ok := cat.Groups[groupname]
	if false == ok {
		// グループが新規追加
		group = &GroupInfo{Title: docname, Description: docname, Docs: map[string]*DocInfo{}}
		cat.Groups[groupname] = group
	}
	if _, ok := group.Docs[docname]; false == ok {
		// ドキュメントが新規追加
		group.Docs[docname] = &DocInfo{Title: docname, Description: docname}
		err = json.SaveToJSONFile(catalogfile, cat.JSON(), true)
	}

	return err
}

// ReadSig は file の署名を取得します。
func ReadSig(host Host, name string) (SigInfo, error) {
	file := host.SigFile(name)
	if di.FileExists(file) {
		// ダウンロードしてあった署名を読む
		return sigFromFile(host, name)
	}
	return nil, fmt.Errorf("%s is not exist", file)
}

// FetchSig はサイト上の署名を取得します。
func FetchSig(host Host, name string) string {
	sigurl := host.SigURL(name)
	if bytes, err := GetBytesHTTP(sigurl); err == nil {
		if elem, err := json.LoadFromJSONByte(bytes); err == nil {
			if sig, err := sigFromElement(host, elem); err == nil {
				return sig.Sig()
			}
		}
	}
	return ""
}

// DownloadSig は URL から署名をダウンロードします。
func DownloadSig(host Host, name string) (SigInfo, error) {
	err := sigDownload(host, name)
	if err != nil {
		return nil, err
	}
	// ダウンロードした署名を読む
	return sigFromFile(host, name)
}

func sigFromElement(host Host, elem json.Element) (*sigInfo, error) {
	// json読み出し
	siginfo := &sig{}
	if es, ok := json.QueryElemString(elem, pathHost); ok {
		siginfo.host = es.Text()
	} else {
		return nil, fmt.Errorf("host not exist")
	}
	if es, ok := json.QueryElemString(elem, pathName); ok {
		siginfo.name = es.Text()
	} else {
		return nil, fmt.Errorf("name not exist : %s", siginfo.host)
	}
	if es, ok := json.QueryElemString(elem, pathURL); ok {
		siginfo.url = es.Text()
	} else {
		return nil, fmt.Errorf("%s", siginfo.host+":"+siginfo.name)
	}
	if es, ok := json.QueryElemString(elem, pathSig); ok {
		siginfo.sig = es.Text()
	} else {
		return nil, fmt.Errorf("%s", siginfo.host+":"+siginfo.name)
	}

	sig := &sigInfo{sig: siginfo, host: host}
	return sig, nil
}

// sigFromFile はファイルから署名を読み込みます。
func sigFromFile(host Host, name string) (*sigInfo, error) {
	file := host.SigFile(name)
	elem, err := json.LoadFromJSONFile(file)
	if err != nil {
		return nil, err
	}
	return sigFromElement(host, elem)
}

func sigDownload(host Host, name string) error {
	sigurl := host.SigURL(name)
	path := host.SigFile(name)
	// di.MkdirAll(host.MyStorePath(), 0755)
	return GetFileHTTP(sigurl, path)
}

// CreateSignatureBASE64 はファイルの署名を生成します。
func CreateSignatureBASE64(key *rsa.PrivateKey, dataFile string) (string, error) {
	// ハッシュ
	h, err := CalcHashFile(dataFile)
	if err != nil {
		return "", err
	}
	// 署名
	//s, err := rsa.SignPSS(rand.Reader, key, crypto.SHA512, h, nil)
	s, err := di.SignPSS(rand.Reader, key, crypto.SHA512, h, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(s), nil
}
