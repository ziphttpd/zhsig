package zhsig

import (
	"github.com/xorvercom/util/pkg/json"
)

const (
	docTitle       = "title"
	docDescription = "description"
)

// DocInfo はドキュメントの表示情報です。
type DocInfo struct {
	// Title はドキュメントの表示名です
	Title string
	// Description はドキュメントの内容の説明です
	Description string
}

// JSON はjsonエレメントを返します。
func (n *DocInfo) JSON() json.ElemObject {
	ret := json.NewElemObject()
	ret.Put(docTitle, json.NewElemString(n.Title))
	ret.Put(docDescription, json.NewElemString(n.Description))
	return ret
}

// DecodeDocInfo はドキュメントの情報をデコードします。
func DecodeDocInfo(en json.ElemObject, docname string) *DocInfo {
	n := &DocInfo{}
	// 表題
	if t, ok := en.Child(docTitle).AsString(); ok {
		n.Title = t.Text()
	}
	if n.Title == "" {
		n.Title = docname
	}
	// 注釈
	if t, ok := en.Child(docDescription).AsString(); ok {
		n.Description = t.Text()
	}
	if n.Description == "" {
		n.Description = n.Title
	}
	return n
}

// DecodeDocInfoMap はドキュメントの情報一覧をデコードします。
func DecodeDocInfoMap(names json.ElemObject) map[string]*DocInfo {
	ret := map[string]*DocInfo{}
	for _, docname := range names.Keys() {
		if en, ok := names.Child(docname).AsObject(); ok {
			// TODO; docname = strings.ToLower(docname)
			ret[docname] = DecodeDocInfo(en, docname)
		}
	}
	return ret
}
