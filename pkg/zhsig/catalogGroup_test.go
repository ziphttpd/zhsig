package zhsig

import (
	"testing"

	"github.com/xorvercom/util/pkg/json"
)

func TestGroupInfo(t *testing.T) {
	group := &GroupInfo{Docs: map[string]*DocInfo{}}
	var group2 *GroupInfo
	groupname := "test"
	str := "abcd"

	group2 = decodeGroup(group.JSON().Text(), groupname)
	if group2.Title != groupname {
		t.Fatal(group2.Title, groupname)
	}

	group.Title = str
	group2 = decodeGroup(group.JSON().Text(), groupname)
	if group2.Title != group.Title {
		t.Fatal(group2.Title, group.Title)
	}
	if group2.Description != group.Title {
		t.Fatal(group2.Description, group.Title)
	}

	group.Description = groupname
	group2 = decodeGroup(group.JSON().Text(), groupname)
	if group2.Title != group.Title {
		t.Fatal(group2.Title, group.Title)
	}
	if group2.Description != group.Description {
		t.Fatal(group2.Description, group.Description)
	}

	dn := group.DocNames()
	if len(dn) != 0 {
		t.Fatal(len(dn), 0)
	}

	// ドキュメントのテスト
	docname := "docname"
	doc := &DocInfo{Title: "doctitle", Description: "docdesc"}
	group.Docs[docname] = doc
	//	log.Print("JSON:", group.JSON().Text())
	group2 = decodeGroup(group.JSON().Text(), groupname)
	if doc2, ok := group2.Docs[docname]; ok {
		if doc2.Title != doc.Title {
			t.Fatal(doc2.Title, doc.Title)
		}
		if doc2.Description != doc.Description {
			t.Fatal(doc2.Description, doc.Description)
		}
	} else {
		t.Fatal("err")
	}

	dn = group.DocNames()
	if len(dn) != 1 {
		t.Fatal(len(dn), 1)
	}
}
func decodeGroup(jsonstr string, docname string) *GroupInfo {
	if ele, err := json.LoadFromJSONByte([]byte(jsonstr)); err == nil {
		if eleo, ok := ele.AsObject(); ok {
			return DecodeGroupInfo(eleo, docname)
		}
	}
	return &GroupInfo{}
}
