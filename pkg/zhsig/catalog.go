package zhsig

import (
	"fmt"

	"github.com/xorvercom/util/pkg/json"
)

const (
	catalogPeer   = "peer"
	catalogGroups = "groups"
)

// Catalog はサイトのカタログです。
type Catalog struct {
	Peer   *PeerInfo
	Groups map[string]*GroupInfo
}

// GroupNames はグループの識別子の一覧を返します。
func (c *Catalog) GroupNames() []string {
	ret := []string{}
	for key := range c.Groups {
		ret = append(ret, key)
	}
	return ret
}

// JSON はjsonエレメントを返します。
func (c *Catalog) JSON() json.ElemObject {
	ret := json.NewElemObject()
	if c.Peer != nil {
		ret.Put(catalogPeer, c.Peer.JSON())
	}
	groups := json.NewElemObject()
	for key, val := range c.Groups {
		groups.Put(key, val.JSON())
	}
	ret.Put(catalogGroups, groups)
	return ret
}

// DownloadCatalog はカタログファイルをダウンロードします。
func DownloadCatalog(host Host) error {
	var cat json.ElemObject
	// ダウンロード
	err := GetFileHTTP(host.CatalogURL(), host.CatalogFile())
	if err != nil {
		return err
	}

	// 読み込み
	if e, err := json.LoadFromJSONFile(host.CatalogFile()); err == nil {
		if o, ok := e.AsObject(); ok {
			cat = o
		} else {
			return fmt.Errorf("%s is invalid Catalog", host.CatalogFile())
		}
	} else {
		return err
	}

	// 証明書情報を追加する
	if peer, err := GetPeerInfo(host.Host()); err == nil {
		cat.Put(catalogPeer, peer.JSON())
	} else {
		return err
	}

	// 更新
	return json.SaveToJSONFile(host.CatalogFile(), cat, true)
}

// NewCatalog はカタログの情報のデフォルト値を生成します。
func NewCatalog() *Catalog {
	return &Catalog{Groups: map[string]*GroupInfo{}}
}

// ReadCatalog はカタログファイルを読み出します。
func ReadCatalog(file string) (*Catalog, error) {
	e, err := json.LoadFromJSONFile(file)
	if err != nil {
		return nil, err
	}
	if eobj, ok := e.AsObject(); ok {
		return DecodeCatalog(eobj), nil
	}
	return nil, fmt.Errorf("%s is invalid Catalog", file)
}

// DecodeCatalog はカタログの情報をデコードします。
func DecodeCatalog(eobj json.ElemObject) *Catalog {
	cat := NewCatalog()
	// ピア情報
	if po, ok := eobj.Child(catalogPeer).AsObject(); ok {
		cat.Peer = DecodePeerInfo(po)
	}
	// グループ情報
	if ego, ok := eobj.Child(catalogGroups).AsObject(); ok {
		cat.Groups = DecodeGroupInfoMap(ego)
	}
	return cat
}
