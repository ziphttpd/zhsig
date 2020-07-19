package zhsig

import (
	"testing"

	"github.com/xorvercom/util/pkg/json"
)

func TestCatalog(t *testing.T) {
	// catalog := &Catalog{Peer: &PeerInfo{}, Groups: map[string]*GroupInfo{}}
	// var catalog2 *Catalog
	// str := "abcd"

	// catalog2 = decodeCatalog(catalog.JSON().Text())
	// if catalog2.Peer.CommonName != "" {
	// 	t.Fatal(catalog2.Peer.CommonName)
	// }

	// catalog.Title = str
	// catalog2 = decodeCatalog(catalog.JSON().Text())
	// if catalog2.Title != catalog.Title {
	// 	t.Fatal(catalog2.Title, catalog.Title)
	// }
	// if catalog2.Description != catalog.Title {
	// 	t.Fatal(catalog2.Description, catalog.Title)
	// }

	// catalog.Description = groupname
	// catalog2 = decodeCatalog(catalog.JSON().Text())
	// if catalog2.Title != catalog.Title {
	// 	t.Fatal(catalog2.Title, catalog.Title)
	// }
	// if catalog2.Description != catalog.Description {
	// 	t.Fatal(catalog2.Description, catalog.Description)
	// }

	// dn := catalog.DocNames()
	// if len(dn) != 0 {
	// 	t.Fatal(len(dn), 0)
	// }

	// // ドキュメントのテスト
	// docname := "docname"
	// doc := &DocInfo{Title: "doctitle", Description: "docdesc"}
	// catalog.Docs[docname] = doc
	// //	log.Print("JSON:", group.JSON().Text())
	// catalog2 = decodeCatalog(catalog.JSON().Text(), groupname)
	// if doc2, ok := catalog2.Docs[docname]; ok {
	// 	if doc2.Title != doc.Title {
	// 		t.Fatal(doc2.Title, doc.Title)
	// 	}
	// 	if doc2.Description != doc.Description {
	// 		t.Fatal(doc2.Description, doc.Description)
	// 	}
	// } else {
	// 	t.Fatal("err")
	// }

	// dn = catalog.DocNames()
	// if len(dn) != 1 {
	// 	t.Fatal(len(dn), 1)
	// }
}
func decodeCatalog(jsonstr string) *Catalog {
	if ele, err := json.LoadFromJSONByte([]byte(jsonstr)); err == nil {
		if eleo, ok := ele.AsObject(); ok {
			return DecodeCatalog(eleo)
		}
	}
	return &Catalog{}
}
