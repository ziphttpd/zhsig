package zhsig

import (
	"testing"

	"github.com/xorvercom/util/pkg/json"
)

func TestDocInfo(t *testing.T) {
	doc := &DocInfo{}
	var doc2 *DocInfo
	docname := "test"
	str := "abcd"

	doc2 = decodeDoc(doc.JSON().Text(), docname)
	if doc2.Title != docname {
		t.Fatal(doc2.Title, docname)
	}

	doc.Title = str
	doc2 = decodeDoc(doc.JSON().Text(), docname)
	if doc2.Title != doc.Title {
		t.Fatal(doc2.Title, doc.Title)
	}
	if doc2.Description != doc.Title {
		t.Fatal(doc2.Description, doc.Title)
	}

	doc.Description = docname
	doc2 = decodeDoc(doc.JSON().Text(), docname)
	if doc2.Title != doc.Title {
		t.Fatal(doc2.Title, doc.Title)
	}
	if doc2.Description != doc.Description {
		t.Fatal(doc2.Description, doc.Description)
	}
}
func decodeDoc(jsonstr string, docname string) *DocInfo {
	if ele, err := json.LoadFromJSONByte([]byte(jsonstr)); err == nil {
		if eleo, ok := ele.AsObject(); ok {
			return DecodeDocInfo(eleo, docname)
		}
	}
	return &DocInfo{}
}
